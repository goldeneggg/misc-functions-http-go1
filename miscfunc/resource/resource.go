package resource

import (
	"context"
	"fmt"
)

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
