package resource

import (
	"context"
	"fmt"
)

type Params struct {
	Path        string
	Method      string
	QueryParams map[string]string
	PathParams  map[string]string
	Header      map[string]string
	Body        string
}

type Result struct {
	Body       string
	StatusCode int
}

func NewResult(msg string, sts int) *Result {
	return &Result{
		Body:       msg,
		StatusCode: sts,
	}
}

func NewResultWithErrorAndStatus(err error, sts int) (*Result, error) {
	return NewResult(err.Error(), sts), err
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
	}

	return nil, fmt.Errorf("invalid params: %#v", p)
}
