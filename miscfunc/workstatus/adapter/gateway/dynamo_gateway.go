package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter"
)

const (
	LOCAL_DYNAMO_CONTAINER_NAME = "localdynamo"
	LOCAL_DYNAMO_PORT           = 8000
)

// DynamoGateway is adapter to DynamoDB workstatus table
type DynamoGateway struct {
	db *dynamodb.Client
}

func NewDynamoGateway() (adapter.Gateway, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	// See: https://github.com/aws/aws-sdk-go-v2/blob/master/example/aws/endpoints/customEndpoint/customEndpoint.go
	if os.Getenv("USE_LOCAL_DYNAMO") == "true" {
		localDynamoURL := fmt.Sprintf(
			"http://%s:%d",
			LOCAL_DYNAMO_CONTAINER_NAME,
			LOCAL_DYNAMO_PORT)
		cfg.EndpointResolver = aws.ResolveWithEndpointURL(localDynamoURL)
	}
	return &DynamoGateway{dynamodb.New(cfg)}, nil
}

func (dg *DynamoGateway) Create(ctx context.Context, workstatus *entity.Workstatus) error {
	content, err := json.Marshal(workstatus.Content)
	if err != nil {
		return err
	}

	item := map[string]dynamodb.AttributeValue{
		"ID": dynamodb.AttributeValue{
			S: aws.String(workstatus.ID),
		},
		"Content": dynamodb.AttributeValue{
			S: aws.String(string(content)),
		},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_TABLE_NAME")),
		Item:      item,
	}

	req := dg.db.PutItemRequest(input)

	_, err = req.Send(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dg *DynamoGateway) Desc(ctx context.Context) (*entity.DescWorkstatus, error) {
	return dg.describeTable(ctx)
}

func (dg *DynamoGateway) describeTable(ctx context.Context) (*entity.DescWorkstatus, error) {
	// Build the request with its input parameters
	req := dg.db.DescribeTableRequest(&dynamodb.DescribeTableInput{
		TableName: aws.String(os.Getenv("DYNAMO_TABLE_NAME")),
	})

	// Send the request, and get the response or error back
	out, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	var attrs []string
	for _, attr := range out.Table.AttributeDefinitions {
		attrs = append(attrs, attr.String())
	}
	status, _ := out.Table.TableStatus.MarshalValue()

	desc := &entity.DescWorkstatus{
		TableName: aws.StringValue(out.Table.TableName),
		Attrs:     attrs,
		Status:    status,
	}

	return desc, nil
}

func (dg *DynamoGateway) listTables(ctx context.Context) (*entity.DescWorkstatus, error) {
	// Build the request with its input parameters
	req := dg.db.ListTablesRequest(&dynamodb.ListTablesInput{})

	// Send the request, and get the response or error back
	_, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	// TODO other columns setting
	desc := &entity.DescWorkstatus{}

	return desc, nil
}
