package resource

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	HelloDefaultHTTPGetAddress = "https://checkip.amazonaws.com"
)

type Hello struct {
}

func (hello *Hello) Get(ctx context.Context, p *Params) (*Result, error) {
	r, err := http.Get(HelloDefaultHTTPGetAddress)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	if r.StatusCode != 200 {
		return NewResultWithErrorAndStatus(errNon200Response, r.StatusCode)
	}

	ip, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 400)
	}

	if len(ip) == 0 {
		return NewResultWithErrorAndStatus(errNoIP, 400)
	}

	return NewResult(fmt.Sprintf("Hello, %v", string(ip)), 200), nil
}

func (hello *Hello) Post(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (hello *Hello) Put(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (hello *Hello) Delete(ctx context.Context, p *Params) (*Result, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}
