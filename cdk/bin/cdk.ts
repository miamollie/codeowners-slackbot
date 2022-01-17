#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { CodeownersStack } from "../lib/codeowners-stack";

const app = new cdk.App();
new CodeownersStack(app, "CdkStack");
