PROJECT := misc-functions-http-go1
CODE := hello-world
API_PATH := /hello
FUNC_NAME := MyGoAPIFunction

MAIN_GO := ./$(CODE)/main.go
CLIENT_GO := ./$(CODE)/client.go
EXE := ./$(CODE)/$(CODE)
LOCAL_EVENT := ./.event.json
API_PORT := 3999
TOOL_DIR := bin
PKG_DLV := github.com/go-delve/delve/cmd/dlv
DLV := $(TOOL_DIR)/dlv
DEBUG_PORT := 5986
TEMPLATE := ./template.yaml
OUTPUT_TEMPLATE := ./serverless-output.yaml
LOCAL_TEMPLATE := ./template.local.yaml
S3_BUCKET := jpshadowapps-lambdas
S3_PREFIX := $(PROJECT)
REGION := ap-northeast-1

#TARGET_PKGS = $(eval TARGET_PKGS := $(shell go list ./... | \grep -v 'vendor'))$(TARGET_PKGS)
TARGET_PKGS = $(shell go list ./... | \grep -v 'vendor')

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	@go get github.com/aws/aws-lambda-go/...@latest

.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 go build -o $(EXE) $(MAIN_GO)

.PHONY: clean
clean:
	@rm -f $(EXE)

.PHONY: rebuild
rebuild: clean build
	
###############################
#
# added targets
#
###############################


###
# for local development
###
mod-dl:
	@GO111MODULE=on go mod download

$(LOCAL_TEMPLATE):
	@[ -f $(LOCAL_TEMPLATE) ] && rm $(LOCAL_TEMPLATE);  cp template.yaml $(LOCAL_TEMPLATE)

api: build $(LOCAL_TEMPLATE)
	@sam local start-api -t $(LOCAL_TEMPLATE) -p $(API_PORT)

curl-get:
	@curl -XGET http://127.0.0.1:$(API_PORT)$(API_PATH)

gen-event:
	@sam local generate-event apigateway aws-proxy > $(LOCAL_EVENT)

###
# for local debug
###
build-dlv:
	@GOOS=linux GOARCH=amd64 go build -o $(DLV) $(PKG_DLV)

debug-build:
	@GOOS=linux GOARCH=amd64 go build -o $(EXE) -gcflags='-N -l' $(MAIN_GO)

# See: https://github.com/awslabs/aws-sam-cli/issues/1067#issuecomment-489406846
# 最新のsam＆delve＆lambciの組み合わせだと動かない
debug-api: debug-build $(LOCAL_TEMPLATE)
	@sam local start-api -t $(LOCAL_TEMPLATE) -p $(API_PORT) -d $(DEBUG_PORT) --debugger-path $(TOOL_DIR) --debug-args="-delveAPI=2" --debug

###
# for test
###
.PHONY: test
test:
	@go test -race ./$(CODE)/...

.PHONY: ci
ci: test

.PHONY: lint
lint:
	@golint -set_exit_status $(TARGET_PKGS)

###
# for package and deploy
###
validate:
	@sam validate

package: build
	@sam package --template-file $(TEMPLATE) --output-template-file $(OUTPUT_TEMPLATE) --s3-bucket $(S3_BUCKET) --s3-prefix $(S3_PREFIX)

# Note: sam cli実行IAMユーザーの権限は、cli用でよく使うPowerUserAccessだけではIAMやCloudFormation周りが不足している
#
# require "AWS_KMS_KEY_ID" env
# @sam deploy --template-file $(OUTPUT_TEMPLATE) --stack-name $(NAME)-stack --capabilities CAPABILITY_IAM --parameter-overrides KeyIdParameter=$$AWS_KMS_KEY_ID
deploy: package
	@sam deploy --template-file $(OUTPUT_TEMPLATE) --stack-name $(PROJECT)-stack --capabilities CAPABILITY_IAM

delete-stack:
	@aws cloudformation delete-stack --stack-name $(PROJECT)-stack

show-api-url:
	@aws cloudformation describe-stacks --stack-name $(PROJECT)-stack | \grep 'execute-api' | sed -e 's/"OutputValue": //g' | tr -d ',' | tr -d '"' | tr -d ' '

###
# for manage go modules
###
mod-tidy:
	@GO111MODULE=on go mod tidy -v

mod-clean:
	@go clean -modcache

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor -v

vendor-build:
	@GOOS=linux GOARCH=amd64 go build -o $(EXE) -mod vendor $(MAIN_GO)

vendor-debug-build:
	@GOOS=linux GOARCH=amd64 go build -gcflags='-N -l' -o $(EXE) -mod vendor $(MAIN_GO)


# DEPRECATE as follows

###
# for manage apigateway
###

# See: https://docs.aws.amazon.com/ja_jp/apigateway/latest/developerguide/create-api-using-awscli.html?shortFooter=true
# See: https://dev.classmethod.jp/cloud/aws/getting-started-with-api-gateway-lambda-integration/

# Note: check and save rest-api-id into .envrc after create
.PHONY: create-apigw
create-apigw:
	@aws apigateway create-rest-api --name $$LOCALENV_APIGW_NAME

.PHONY: check-apigw
check-apigw:
	@aws apigateway get-rest-api --rest-api-id $$LOCALENV_APIGW_ID

.PHONY: check-resources
check-resources:
	@aws apigateway get-resources --rest-api-id $$LOCALENV_APIGW_ID

.PHONY: add-resource-under-root
add-resource-under-root:
	@read -p 'Input path?: ' path; aws apigateway create-resource --rest-api-id $$LOCALENV_APIGW_ID --parent-id $$LOCALENV_APIGW_ROOT_RESOURCE_ID --path-part $$path

.PHONY: add-get-method
add-get-method:
	@read -p 'Input resource_id?: ' rid; aws apigateway put-method --rest-api-id $$LOCALENV_APIGW_ID --resource-id $$rid --http-method GET --authorization-type "NONE"

.PHONY: add-get-method-response-ok
add-get-method-response-ok:
	@read -p 'Input resource_id?: ' rid; aws apigateway put-method-response --rest-api-id $$LOCALENV_APIGW_ID --resource-id $$rid --http-method GET --status-code 200

.PHONY: put-http-integration-with-url
put-http-integration-with-url:
	@read -p 'Input resource_id?: ' rid; read -p 'Input through URL?: ' url; aws apigateway put-integration --rest-api-id $$LOCALENV_APIGW_ID --resource-id $$rid --http-method GET --type HTTP --integration-http-method GET --uri $$url

.PHONY: put-http-integration-response
put-http-integration-response:
	@read -p 'Input resource_id?: ' rid; aws apigateway put-integration-response --rest-api-id $$LOCALENV_APIGW_ID --resource-id $$rid --http-method GET --status-code 200 --selection-pattern ""

# ステージ = test へのデプロイ
.PHONY: deploy-apigw-test
deploy-apigw-test:
	@aws apigateway create-deployment --rest-api-id $$LOCALENV_APIGW_ID --stage-name test --stage-description 'Test stage'

# ステージ = test にデプロイされたAPIのURL確認
.PHONY: show-test-stage-url
show-test-stage-url:
	@echo https://$$LOCALENV_APIGW_ID.execute-api.$(REGION).amazonaws.com/test

.PHONY: curl-test-hoge
curl-test-hoge:
	@curl $(shell $(MAKE) show-test-stage-url)/hoge

#
# .PHONY: invoke-server
# invoke-server:
# 	@_LAMBDA_SERVER_PORT=$(INVOKE_PORT) go run $(MAIN_GO)
#
# .PHONY: invoke-client
# invoke-client:
# 	@_LAMBDA_SERVER_PORT=$(INVOKE_PORT) go run $(CLIENT_GO) "hoge"
#
# VARS := HOGE_ENV=huga
# IAM_ROLE := arn:aws:iam::271505164757:role/MyAPICentral
#
# .PHONY: create-func
# create-func: zip
# 	@aws lambda create-function --function-name $(FUNC_NAME) --zip-file fileb://./$(BINNAME).zip --runtime go1.x --handler $(BINNAME) --role $(IAM_ROLE) --region $(REGION)
#
# .PHONY: update-func
# update-func: zip
# 	@aws lambda update-function-code --function-name $(FUNC_NAME) --zip-file fileb://./$(BINNAME).zip --region $(REGION)
#
# .PHONY: update-func-conf
# update-func-conf:
# 	@aws lambda update-function-configuration --function-name $(FUNC_NAME) --environment Variables={$(VARS)} --region $(REGION)
#
# .PHONY: invoke-func
# invoke-func:
# 	@read -p 'Input version?: ' ver; aws lambda invoke --function-name $(FUNC_NAME) --payload file://input.json --qualifier $$ver result.json
#
# .PHONY: show-current-version
# bumpup:
# 	@read -p 'Input description of new version?: ' desc; aws lambda publish-version --function-name $(FUNC_NAME) --description $$desc
#
# .PHONY: bumpup
# bumpup:
# 	@read -p 'Input description of new version?: ' desc; aws lambda publish-version --function-name $(FUNC_NAME) --description $$desc