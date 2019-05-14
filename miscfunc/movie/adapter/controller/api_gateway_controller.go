package controller

import (
	"context"

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
	proxyReq events.APIGatewayProxyRequest,
	proxyResp events.APIGatewayProxyResponse) adapter.Controller {

	return &APIGatewayController{
		uc:        uc,
		proxyReq:  proxyReq,
		proxyResp: proxyResp,
	}
}

func (ac *APIGatewayController) Create(ctx context.Context) (*entity.Movie, error) {
	// TODO
	return &entity.Movie{}, nil
}
