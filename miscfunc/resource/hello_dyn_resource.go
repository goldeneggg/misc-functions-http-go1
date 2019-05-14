package resource

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie/adapter"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie/adapter/controller"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie/adapter/gateway"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie/usecase"
)

type HelloDyn struct {
	Service string
}

func (helloDyn *HelloDyn) Get(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (helloDyn *HelloDyn) Post(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctrl, err := helloDyn.newController(ctx, proxyReq)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	_, err = ctrl.Create(ctx)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	return NewResultWithHeader("", 201, map[string]string{"Content-Type": "application/json"}), nil
}

func (helloDyn *HelloDyn) Put(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (helloDyn *HelloDyn) Delete(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (helloDyn *HelloDyn) newController(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (adapter.Controller, error) {
	gw, err := gateway.NewDynamoGateway()
	if err != nil {
		return nil, err
	}
	return controller.NewAPIGatewayController(usecase.NewDefaultUsecase(gw), proxyReq), nil
}
