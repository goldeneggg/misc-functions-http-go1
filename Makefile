PROJECT := misc-functions-http-go1
CODE := miscfunc
API_PATH := /hello
FUNC_NAME := MyGoAPIFunction

PKG_AWS_LAMBDA_GO := github.com/aws/aws-lambda-go
PKG_AWS_SDK_GO_V2 := github.com/aws/aws-sdk-go-v2

MAIN_GO := ./$(CODE)/main.go
CLIENT_GO := ./$(CODE)/client.go
EXE := ./$(CODE)/$(CODE)
LOCAL_ENVS := ./localenvs.json
LOCAL_EVENT := ./.event.json
API_PORT := 3999
DOCKER_LOCAL_DYNAMO_NAME := localdynamo
DOCKER_LOCAL_DYNAMO_NETWORK := sam-dynamo-network
TESTDATA_CREATE_TABLE := testdata/skel-workstatus-create-table.json
TESTDATA_DATA := testdata/data.json
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
setup-deps:
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
# for local API Gateway development
###
cp-tmpl-local:
	@[ -f $(LOCAL_TEMPLATE) ] && rm $(LOCAL_TEMPLATE);  cp template.yaml $(LOCAL_TEMPLATE)

api: build cp-tmpl-local
	@sam local start-api -n $(LOCAL_ENVS) -t $(LOCAL_TEMPLATE) -p $(API_PORT)

# Note: must execute up-dynamo target before start-api
api-with-localdynamo: build cp-tmpl-local
	@sam local start-api -n $(LOCAL_ENVS) -t $(LOCAL_TEMPLATE) -p $(API_PORT) --docker-network $(DOCKER_LOCAL_DYNAMO_NETWORK)

api-samdev: build cp-tmpl-local
	@samdev local start-api -n $(LOCAL_ENVS) -t $(LOCAL_TEMPLATE) -p $(API_PORT)

# Note: must execute run-local-dynamo target before start-api
api-samdev-with-localdynamo: build cp-tmpl-local
	@samdev local start-api -n $(LOCAL_ENVS) -t $(LOCAL_TEMPLATE) -p $(API_PORT) --docker-network $(DOCKER_LOCAL_DYNAMO_NETWORK)

curl-get:
	@curl -XGET http://127.0.0.1:$(API_PORT)$(API_PATH)

curl-get-workstatus-desc:
	@curl -XGET "http://127.0.0.1:$(API_PORT)/workstatus?type=desc"

curl-post-workstatus:
	@curl -XPOST -d @$(TESTDATA_DATA) "http://127.0.0.1:$(API_PORT)/workstatus" 

setup-local-event:
	@sam local generate-event apigateway aws-proxy > $(LOCAL_EVENT)

###
# for local DynamoDB
###

# See:
# - https://docs.aws.amazon.com/ja_jp/amazondynamodb/latest/developerguide/DynamoDBLocal.html
# - https://github.com/aws-samples/aws-sam-java-rest

_pull-dynamo:
	@docker pull amazon/dynamodb-local

_create-docker-network-for-dynamo:
	@docker network create $(DOCKER_LOCAL_DYNAMO_NETWORK)

setup-dynamo-docker: _pull-dynamo _create-docker-network-for-dynamo

setup-dynamo-workstatus-table-skelton:
	@aws dynamodb create-table --generate-cli-skeleton > $(TESTDATA_CREATE_TABLE)

up-dynamo: _run-dynamo _create-workstatus-table

_run-dynamo:
	@docker run -d --rm -p 8000:8000 --name $(DOCKER_LOCAL_DYNAMO_NAME) --network $(DOCKER_LOCAL_DYNAMO_NETWORK) amazon/dynamodb-local

_create-workstatus-table:
	@aws dynamodb create-table --cli-input-json file://$(TESTDATA_CREATE_TABLE) --endpoint-url http://localhost:8000

down-dynamo:
	@docker stop $(DOCKER_LOCAL_DYNAMO_NAME)

delete-workstatus-table:
	@aws dynamodb delete-table --table-name localtable --endpoint-url http://localhost:8000

scan-table:
	@aws dynamodb scan --table-name localtable --endpoint-url http://localhost:8000

list-table:
	@aws dynamodb list-tables --endpoint-url http://localhost:8000

desc-table:
	@aws dynamodb describe-table --table-name workstatus --endpoint-url http://localhost:8000


###
# for local debug
###
build-dlv:
	@GOOS=linux GOARCH=amd64 go build -o $(DLV) $(PKG_DLV)

debug-build:
	@GOOS=linux GOARCH=amd64 go build -o $(EXE) -gcflags='-N -l' $(MAIN_GO)

# See: https://github.com/awslabs/aws-sam-cli/issues/1067#issuecomment-489406846
# 最新のsam＆delve＆lambciの組み合わせだと動かない
debug-api: build-dlv debug-build cp-tmpl-local
	@sam local start-api -t $(LOCAL_TEMPLATE) -p $(API_PORT) -d $(DEBUG_PORT) --debugger-path $(TOOL_DIR) --debug-args="-delveAPI=2" --debug

debug-api-samdev: build-dlv debug-build cp-tmpl-local
	@samdev local start-api -t $(LOCAL_TEMPLATE) -p $(API_PORT) -d $(DEBUG_PORT) --debugger-path $(TOOL_DIR) --debug-args="-delveAPI=2" --debug

###
# for test
###
.PHONY: test
test:
	@go test -race ./$(CODE)/...

.PHONY: lint
lint:
	@golint -set_exit_status $(TARGET_PKGS)

.PHONY: vet
vet:
	@go vet $(TARGET_PKGS)

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
	@sam deploy --template-file $(OUTPUT_TEMPLATE) --stack-name $(PROJECT)-stack --capabilities CAPABILITY_IAM --parameter-overrides KeyIdParameter=$$AWS_KMS_KEY_ID

describe-stack-latest-event:
	@aws cloudformation describe-stack-events --stack-name $(PROJECT)-stack --max-items 1

delete-stack:
	@aws cloudformation delete-stack --stack-name $(PROJECT)-stack

show-api-url:
	@aws cloudformation describe-stacks --stack-name $(PROJECT)-stack | \grep 'execute-api' | sed -e 's/"OutputValue": //g' | tr -d ',' | tr -d '"' | tr -d ' '

###
# for CI
###
.PHONY: ci
ci: test vet validate

lint-travis:
	@travis lint --org --debug .travis.yml

###
# for manage go modules
###
mod-dl:
	@GO111MODULE=on go mod download

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

chk_versions = go list -u -m -versions $1 | tr ' ' '\n'

chk-versions-aws-lambda-go:
	@$(call chk_versions,$(PKG_AWS_LAMBDA_GO))

chk-versions-aws-sdk-go-v2:
	@$(call chk_versions,$(PKG_AWS_SDK_GO_V2))

###
# for add new feature
###

# See: https://hackernoon.com/golang-clean-archithecture-efd6d7c43047
new-ca:
	@read -p 'Input new domain name?: ' name; \
		mkdir -p $(CODE)/entity && touch $(CODE)/entity/$$name.go && \
		mkdir -p $(CODE)/$$name/{delivery,repository,usecase} && \
		touch $(CODE)/$$name/repository.go $(CODE)/$$name/usecase.go $(CODE)/$$name/repository/.gitkeep $(CODE)/$$name/usecase/.gitkeep $(CODE)/$$name/delivery/.gitkeep

# original version
new-domain:
	@read -p 'Input new domain name?: ' name; \
		mkdir -p $(CODE)/entity && touch $(CODE)/entity/$$name.go && \
		mkdir -p $(CODE)/$$name/{usecase,adapter/{controller,gateway},infra} && \
		touch $(CODE)/$$name/usecase.go $(CODE)/$$name/adapter/controller.go $(CODE)/$$name/adapter/gateway.go && \
		touch $(CODE)/$$name/usecase/.gitkeep $(CODE)/$$name/adapter/controller/.gitkeep $(CODE)/$$name/adapter/gateway/.gitkeep $(CODE)/$$name/infra/.gitkeep

# DEPRECATEs as follows

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
