package resource

import (
	"context"
)

type Crawler struct {
	Service string
}

func (crawler *Crawler) Get(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (crawler *Crawler) Post(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (crawler *Crawler) Put(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (crawler *Crawler) Delete(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}
