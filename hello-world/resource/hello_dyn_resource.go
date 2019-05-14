package resource

import (
	"context"
	"encoding/json"

	"github.com/goldeneggg/misc-functions-http-go1/hello-world/store/dynamo"
)

type HelloDyn struct {
	Service string
}

func (helloDyn *HelloDyn) Get(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (helloDyn *HelloDyn) Post(ctx context.Context, p *Params) (*Result, error) {
	var movie dynamo.Movie

	err := json.Unmarshal([]byte(p.Body), &movie)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	err = movie.PutItem(ctx)
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
