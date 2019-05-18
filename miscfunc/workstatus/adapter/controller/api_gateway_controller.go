package controller

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter"
)

type APIGatewayController struct {
	uc        workstatus.Usecase
	proxyReq  events.APIGatewayProxyRequest
	proxyResp events.APIGatewayProxyResponse
}

func NewAPIGatewayController(
	uc workstatus.Usecase,
	proxyReq events.APIGatewayProxyRequest) adapter.Controller {

	return &APIGatewayController{
		uc:        uc,
		proxyReq:  proxyReq,
		proxyResp: events.APIGatewayProxyResponse{},
	}
}

func (ac *APIGatewayController) Create(ctx context.Context) (*entity.Workstatus, error) {
	var workstatus entity.Workstatus

	err := json.Unmarshal([]byte(ac.proxyReq.Body), &workstatus)
	if err != nil {
		return nil, err
	}

	return ac.uc.Create(ctx, &workstatus)
}

func (ac *APIGatewayController) Desc(ctx context.Context) (*entity.DescWorkstatus, error) {
	return ac.uc.Desc(ctx)
}
