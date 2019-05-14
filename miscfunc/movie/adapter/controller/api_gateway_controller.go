package controller

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie/adapter"
)

type APIGatewayController struct {
	uc        movie.Usecase
	proxyReq  events.APIGatewayProxyRequest
	proxyResp events.APIGatewayProxyResponse
}

func NewAPIGatewayController(
	uc movie.Usecase,
	proxyReq events.APIGatewayProxyRequest) adapter.Controller {

	return &APIGatewayController{
		uc:        uc,
		proxyReq:  proxyReq,
		proxyResp: events.APIGatewayProxyResponse{},
	}
}

func (ac *APIGatewayController) Create(ctx context.Context) (*entity.Movie, error) {
	var movie entity.Movie

	err := json.Unmarshal([]byte(ac.proxyReq.Body), &movie)
	if err != nil {
		return nil, err
	}

	return ac.uc.Create(ctx, &movie)
}
