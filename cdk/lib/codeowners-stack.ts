import * as cdk from "@aws-cdk/core";

import {
  CodeownersApi,
  CodeownersLambda,
} from "./constructs/codeownersResources";

export class CodeownersStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const { codeownersLambda } = new CodeownersLambda(this, "CodeownersLambda");

    new CodeownersApi(this, "slackbot-codeowner", {
      lambda: codeownersLambda,
    });
  }
}
