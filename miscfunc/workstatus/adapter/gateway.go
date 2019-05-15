package adapter

import (
	"context"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
)

type Gateway interface {
	Create(ctx context.Context, workstatus *entity.Workstatus) error
}
