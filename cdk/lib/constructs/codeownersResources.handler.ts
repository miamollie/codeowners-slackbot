import * as lambda from "aws-lambda";
import { request, gql } from "graphql-request";
import { URLSearchParams } from "url";
import { createHmac, timingSafeEqual } from "crypto";

const API_ENDPOINT = process.env.GITHUB_GQL_API || "";
const requestHeaders = {
  authorization: `Bearer ${process.env.GITHUB_GQL_AUTH_TOKEN}`,
};

const SLACK_SIGNING_SECRET = process.env.SLACK_SIGNING_SECRET || "";

interface RespData {
  repository: {
    file1: { text: string };
    file2: { text: string };
    file3: { text: string };
  } | null;
  errors:
    | {
        type: string;
      }[]
    | null;
}

const QUERY = gql`
  query owners($name: String!, $owner: String!) {
    repository(name: $name, owner: $owner) {
      file1: object(expression: "master:CODEOWNERS") {
        ... on Blob {
          text
        }
      }
      file2: object(expression: "master:.github/CODEOWNERS") {
        ... on Blob {
          text
        }
      }
      file3: object(expression: "master:docs/CODEOWNERS") {
        ... on Blob {
          text
        }
      }
    }
  }
`;

interface VarType {
  name: string;
  owner: string;
}

export const handler = async (
  event: lambda.APIGatewayProxyEvent
): Promise<lambda.APIGatewayProxyResult> => {

  if (!verifySlackRequest(event)) {
    console.log("Failed to verify request from slack");
    return {
      statusCode: 500,
      body: "Unable to verify request from slack",
    };
  }

  const variables = getVarsFromParams(event) || getVarsFromBody(event);
  if (!variables) {
    console.log("Failed to get variables from request");
    return {
      statusCode: 200,
      body: "Sorry, couldn't read a repo and owner from that request",
    };
  }

  const resp = await request<RespData>(
    API_ENDPOINT,
    QUERY,
    variables,
    requestHeaders
  ).catch((e) => console.log("Error " + e));

  if (!resp) {
    console.log("Empty response from github api");
    return {
      statusCode: 200,
      body: "Sorry, couldn't get a response from github api",
    };
  }

  return {
    statusCode: 200,
    body: getMessageFromResp(resp),
  };
};

function getVarsFromParams(
  event: lambda.APIGatewayProxyEvent
): VarType | undefined {
  if (!event?.queryStringParameters?.text) {
    return undefined;
  }
  const args = event.queryStringParameters?.text.split("/");
  if (args.length !== 2) {
    return undefined;
  }

  if (!args[0] || !args[1]) {
    return undefined;
  }

  return { owner: args[0], name: args[1] };
}

function getVarsFromBody(
  event: lambda.APIGatewayProxyEvent
): VarType | undefined {
  if (!event?.body) {
    return undefined;
  }

  const params = new URLSearchParams(event.body);

  if (!params.has("text")) {
    return undefined;
  }

  const args = params.get("text")?.split("/");
  if (!args || args.length !== 2) {
    return undefined;
  }

  if (!args[0] || !args[1]) {
    return undefined;
  }

  return { owner: args[0], name: args[1] };
}

function getMessageFromResp(resp: RespData) {
  if (resp.errors) {
    console.log(resp.errors);
    return (
      resp.errors[0].type ||
      "Sorry, this repo doesn't exist or you do not have permission to view it"
    );
  }

  const owners =
    resp?.repository?.file1?.text ||
    resp?.repository?.file2?.text ||
    resp?.repository?.file3?.text;

  if (!owners) {
    return "Sorry, couldn't get a codeowners file for this repo – either you don't have access to view it, or it hasn't been added";
  }

  return `Codeowners: ${owners}`;
}

const SLACK_TIMESTAMP_HEADER = "X-Slack-Request-Timestamp";
const SLACK_SIGNATURE_HEADER = "X-Slack-Signature";
// Concat version number, request timestamp, body of request with colon as delimiter
// Hash basestring with HMAC SHA256 using signing secret as key
// compare to X-Slack-Signature header
function verifySlackRequest(event: lambda.APIGatewayProxyEvent) {
  if (!SLACK_SIGNING_SECRET) {
    console.log("signing secret not set");
    return false;
  }

  const { body } = event;

  const timestamp = event.headers[SLACK_TIMESTAMP_HEADER];
  const receivedSignature = event.headers[SLACK_SIGNATURE_HEADER];

  if (!body || !timestamp || !receivedSignature) {
    return false;
  }

  const [version, slack_hash] = receivedSignature.split("=");
  const baseString = `${version}:${timestamp}:${body}`;

  const generatedSignature = createHmac("sha256", SLACK_SIGNING_SECRET)
    .update(baseString, "utf8")
    .digest("hex");

  return timingSafeEqual(
    Buffer.from(generatedSignature, "utf8"),
    Buffer.from(slack_hash, "utf8")
  );
}
