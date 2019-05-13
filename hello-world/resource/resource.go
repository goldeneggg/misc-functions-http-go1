package resource

import (
	"context"
	"errors"
	"fmt"
)

var (
	errNon200Response = errors.New("Non 200 Response found")
	errNoIP           = errors.New("No IP in HTTP response")
	errNotImplemented = errors.New("Not Implemented")
)

type Params struct {
	Path        string
	Method      string
	QueryParams map[string]string
	PathParams  map[string]string
	Header      map[string]string
	Body        string
	Stage       string
}

type Result struct {
	Body       string
	StatusCode int
	Header     map[string]string
}

func NewResult(msg string, sts int) *Result {
	return &Result{
		Body:       msg,
		StatusCode: sts,
	}
}

func NewResultWithHeader(msg string, sts int, header map[string]string) *Result {
	result := NewResult(msg, sts)
	result.Header = header

	return result
}

func NewResultWithErrorAndStatus(err error, sts int) (*Result, error) {
	return NewResult(err.Error(), sts), err
}

type Resource interface {
	Get(ctx context.Context, p *Params) (*Result, error)
	Post(ctx context.Context, p *Params) (*Result, error)
	Put(ctx context.Context, p *Params) (*Result, error)
	Delete(ctx context.Context, p *Params) (*Result, error)
}

func newResource(ctx context.Context, p *Params) (Resource, error) {
	switch p.Path {
	case "/hello":
		return &Hello{}, nil
	case "/crawler":
		return &Crawler{}, nil
	}

	return nil, fmt.Errorf("invalid params: %#v", p)
}

func Access(ctx context.Context, p *Params) (*Result, error) {
	r, err := newResource(ctx, p)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	switch p.Method {
	case "GET":
		return r.Get(ctx, p)
	case "POST":
		return r.Post(ctx, p)
	case "PUT":
		return r.Put(ctx, p)
	case "DELETE":
		return r.Delete(ctx, p)
	default:
		return NewResultWithErrorAndStatus(fmt.Errorf("invalid request method: %s", p.Method), 500)
	}
}
