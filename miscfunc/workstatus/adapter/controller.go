package adapter

import (
	"context"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
)

type Controller interface {
	Create(ctx context.Context) (*entity.Workstatus, error)
	Desc(ctx context.Context) (*entity.DescWorkstatus, error)
}
