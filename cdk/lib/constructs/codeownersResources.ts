import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-nodejs";
import { LambdaRestApi, } from "@aws-cdk/aws-apigateway";

export class CodeownersApi extends cdk.Construct {
  public readonly gateway: LambdaRestApi;

  constructor(
    scope: cdk.Construct,
    id: string,
    props: { lambda: lambda.NodejsFunction }
  ) {
    super(scope, id);

    this.gateway = new LambdaRestApi(this, "slackbot-codeowner", {
      handler: props.lambda,
      description: "api gateway for github codeowners slackbot",
    });
  }
}

export class CodeownersLambda extends cdk.Construct {
  public readonly codeownersLambda: lambda.NodejsFunction;

  constructor(scope: cdk.Construct, id: string) {
    super(scope, id);

    this.codeownersLambda = new lambda.NodejsFunction(this, "handler", {
      environment: {
        GITHUB_GQL_API: "https://api.github.com/graphql",
        GITHUB_GQL_AUTH_TOKEN: process.env.GITHUB_GQL_AUTH_TOKEN || "",
      },
    });
  }
}
