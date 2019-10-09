# misc-functions-http-go1

## Requirements
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* AWS CLI already configured with Administrator permission
  * If using OS X and homebrew, run `brew install awscli`
* aws-sam-cli already configured with Administrator permission
  * If using OS X and homebrew, run `brew install aws-sam-cli`

## Setup process

```sh
# setup
$ git clone THIS_REPOS
$ cd THIS_REPOS
$ make mod-dl
```

## Local development

**Invoking function locally through local API Gateway**
See: miscfunc/resource/hello_resource.go

```sh
# run API on local
$ make api

# confirm API (See: main.go and miscfunc/resource/hello_resource.go)
$ make curl-get-hello
Hello, 111.108.8.42

# stop API
(press ctrl + c)
```

### With Local DynamoDB
See: miscfunc/resource/workstatus_resource.go

Setup local DynamoDB environment (only once)

```sh
$ make setup-dynamo-docker
```

Run local DynanoDB container  __Notice: This feature uses port 8000.__

```sh
$ make up-dynamo

# confirm running container
$ docker ps
CONTAINER ID        IMAGE                   COMMAND                  CREATED             STATUS              PORTS                    NAMES
XXXXXXXXXXXX        amazon/dynamodb-local   "java -jar DynamoDBL…"   43 seconds ago      Up 40 seconds       0.0.0.0:8000->8000/tcp   localdynamo
```

Invoke function with local DynamoDB

```sh
# run API on local with DynamoDB
$ make api-with-localdynamo

# get table describe from local DynamoDB table
$ make curl-get-workstatus-desc
{"table_name":"localtable","attrs":["{\n  AttributeName: \"ID\",\n  AttributeType: S\n}","{\n  AttributeName: \"Content\",\n  AttributeType: S\n}"],"status":"ACTIVE"}%

# post testdata to local DynamoDB table
$ make curl-post-workstatus
no data

# confirm posted testdata from local DynamoDB table
$ make scan-table
{
    "Items": [
        {
            "Content": {
                "S": "{\"ym\":201906,\"buffer_days_per_week\":1,\"desc\":\"テスト説明文\"}"
            },
            "ID": {
                "S": "1"
            }
        }
    ],
    "Count": 1,
    "ScannedCount": 1,
    "ConsumedCapacity": null
}

# stop API
(press ctrl + c)

# stop DynamoDB container
$ make down-dynamo
```

### Write and Run tests

```sh
$ make test

# Run golint
$ make lint

# Run go vet
$ make vet
```

## Packaging and deployment

### Edit and Validate template.yaml
Write or Modify template.yaml

```sh
$ vim template.yaml
```

Validate template.yaml

```sh
$ make validate
/Users/xxx/misc-functions-http-go1/template.yaml is a valid SAM Template
```

### Deploy to NOT production



### Deploy to production




## Managing go modules for dependencies
See: https://horizoon.jp/post/2019/04/18/contributing_with_gomodules/ (only Japanese)

### Check available versions

```sh
# for aws-lambda-go
$ make chk-versions-aws-lambda-go
:
v1.13.0
v1.13.1
v1.13.2

# for aws-sdk-go-v2
$ make chk-versions-aws-sdk-go-v2
:
v0.13.0
v0.14.0
v2.0.0-preview.1+incompatible
v2.0.0-preview.2+incompatible
v2.0.0-preview.3+incompatible
v2.0.0-preview.4+incompatible
```

### Update version

Assign target version using module-query

```sh
# for aws-lambda-go
$ make update-aws-lambda-go
Input Module Query(e.g. "<v1.20")?: <v1.14

query=<v1.14
go: downloading github.com/aws/aws-lambda-go v1.13.2
go: extracting github.com/aws/aws-lambda-go v1.13.2


# for aws-sdk-go-v2
$ make update-aws-sdk-go-v2
Input Module Query(e.g. "<v1.20")?: <v0.15

query=<v0.15
go: downloading github.com/aws/aws-sdk-go-v2 v0.14.0
go: extracting github.com/aws/aws-sdk-go-v2 v0.14.0
```



