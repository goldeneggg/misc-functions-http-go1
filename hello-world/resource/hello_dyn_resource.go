package resource

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type HelloDyn struct {
	Service string
}

type Movie struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (helloDyn *HelloDyn) Get(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (helloDyn *HelloDyn) Post(ctx context.Context, p *Params) (*Result, error) {
	var movie Movie

	err := json.Unmarshal([]byte(p.Body), &movie)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	svc := dynamodb.New(cfg)
	req := svc.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(p.QueryParams["tn"]),
		Item: map[string]dynamodb.AttributeValue{
			"ID": dynamodb.AttributeValue{
				S: aws.String(movie.ID),
			},
			"Name": dynamodb.AttributeValue{
				S: aws.String(movie.Name),
			},
		},
	})

	_, err = req.Send()
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	return NewResultWithHeader("", 201, map[string]string{"Content-Type": "application/json"}), nil
}

func (helloDyn *HelloDyn) Put(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (helloDyn *HelloDyn) Delete(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}
