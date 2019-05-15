package gateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter"
)

const (
	TABLE_NAME = "workstatus"
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
			"Content": dynamodb.AttributeValue{
				S: aws.String(workstatus.Content),
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
