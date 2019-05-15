package resource

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter/controller"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter/gateway"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/usecase"
)

type Workstatus struct {
}

func (ws *Workstatus) Get(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (ws *Workstatus) Post(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctrl, err := ws.newController(ctx, proxyReq)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	_, err = ctrl.Create(ctx)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	return NewResultWithHeader("", 201, map[string]string{"Content-Type": "application/json"}), nil
}

func (ws *Workstatus) Put(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (ws *Workstatus) Delete(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func newWorkstatusResource(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (*Workstatus, error) {
	return &Workstatus{}, nil
}

func (ws *Workstatus) newController(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (adapter.Controller, error) {
	gw, err := gateway.NewDynamoGateway()
	if err != nil {
		return nil, err
	}
	return controller.NewAPIGatewayController(usecase.NewDefaultUsecase(gw), proxyReq), nil
}
