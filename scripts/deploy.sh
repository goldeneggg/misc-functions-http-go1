#!/bin/bash

set -eu

TEMPLATE=./template.yaml
OUTPUT_TEMPLATE=./serverless-output.yaml
S3_BUCKET=jpshadowapps-lambdas
PROJECT=misc-functions-http-go1

sam package --template-file ${TEMPLATE} \
  --output-template-file ${OUTPUT_TEMPLATE} \
  --s3-bucket ${S3_BUCKET} \
  --s3-prefix ${PROJECT}

sam deploy --template-file ${OUTPUT_TEMPLATE} \
  --stack-name ${PROJECT}-stack \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides KeyIdParameter=${AWS_KMS_KEY_ID}
