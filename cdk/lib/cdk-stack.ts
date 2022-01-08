import * as cdk from "@aws-cdk/core";
import { LambdaRestApi } from "@aws-cdk/aws-apigateway";
import * as lambda from "@aws-cdk/aws-lambda-go";

export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const l = new lambda.GoFunction(this, "handler", {
      entry: "../",
      environment: {
        GITHUB_GQL_API: "https://api.github.com/graphql",
        GITHUB_GQL_AUTH_TOKEN: process.env.GITHUB_GQL_AUTH_TOKEN || "",
      },
    });

    new LambdaRestApi(this, "slackbot-codeowner", {
      handler: l,
      description: "api gateway for golang slack handler",
    });
  }
}
