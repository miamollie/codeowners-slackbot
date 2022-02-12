#!/usr/bin/env bash
export $(cat .env | xargs)
cd ./cdk
npm run cdk deploy --all 
exit $?


