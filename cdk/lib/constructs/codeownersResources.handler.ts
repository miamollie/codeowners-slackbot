import * as lambda from "aws-lambda";
import { request, gql } from "graphql-request";
import { URLSearchParams } from "url";

const endpoint = process.env.GITHUB_GQL_API || "";
const requestHeaders = {
  authorization: `Bearer ${process.env.GITHUB_GQL_AUTH_TOKEN}`,
};

// interface RespData {
//   getCodeowners: { repository: {file1: }};
// }

const query = gql`
  query owners($name: String!, $owner: String!) {
    repository(name: $name, owner: $owner) {
      file1: object(expression: "master:CODEOWNERS") {
        ... on Blob {
          text
        }
      }
      file2: object(expression: "master:.git/CODEOWNERS") {
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
  console.log(event);

  const variables = getVarsFromParams(event) || getVarsFromBody(event);
  if (!variables) {
    return {
      statusCode: 200,
      body: "Sorry, couldn't read a repo and owner from that request",
    };
  }

  console.table(variables);

  // TODO add type to response https://github.com/prisma-labs/graphql-request/blob/master/examples/passing-more-options-to-fetch.ts
  const resp = await request(endpoint, query, variables, requestHeaders).catch(
    (e) => console.log("Error " + e)
  );

  console.table(resp);

  //TODO slack API for JS
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

  const params = new URLSearchParams(event.body)

  if (!params.has("text")) {
    return undefined
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

function getMessageFromResp(resp: any) {
  if (resp.errors) {
    console.log(resp.errors);
    return resp.errors[0].type || "Sorry, this repo doesn't exist or you do not have permission to view it";
  }

  const owners =
    resp?.repository?.file1?.text ||
    resp?.repository?.file2?.text ||
    resp?.repository?.file3?.text;

  if (!owners) {
    return "Sorry, couldn't get a codeowners file for this repo – either you don't have access to view it, or it hasn't been added";
  }

  return `Codeowners for this repo are: ${owners}`;
}
