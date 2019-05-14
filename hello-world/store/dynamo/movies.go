package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	TABLE_NAME = "movies"
)

type Movie struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (movie *Movie) PutItem(ctx context.Context) error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	svc := dynamodb.New(cfg)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]dynamodb.AttributeValue{
			"ID": dynamodb.AttributeValue{
				S: aws.String(movie.ID),
			},
			"Name": dynamodb.AttributeValue{
				S: aws.String(movie.Name),
			},
		},
	}
	req := svc.PutItemRequest(input)

	_, err = req.Send(ctx)
	if err != nil {
		return err
	}

	return nil
}
