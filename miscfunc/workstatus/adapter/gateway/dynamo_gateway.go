package gateway

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter"
)

const (
	TABLE_NAME = "workstatuss"
)

type DynamoGateway struct {
	db *dynamodb.DynamoDB
}

func NewDynamoGateway() (adapter.Gateway, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	return &DynamoGateway{dynamodb.New(cfg)}, nil
}

func (dg *DynamoGateway) Create(ctx context.Context, workstatus *entity.Workstatus) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]dynamodb.AttributeValue{
			"ID": dynamodb.AttributeValue{
				S: aws.String(workstatus.ID),
			},
			"YM": dynamodb.AttributeValue{
				N: aws.String(strconv.FormatInt(workstatus.YM, 10)),
			},
			"Desc": dynamodb.AttributeValue{
				S: aws.String(workstatus.Desc),
			},
		},
	}
	req := dg.db.PutItemRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return err
	}

	return nil
}
