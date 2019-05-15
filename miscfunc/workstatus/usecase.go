package workstatus

import (
	"context"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
)

type Usecase interface {
	Create(ctx context.Context, workstatus *entity.Workstatus) (*entity.Workstatus, error)
}
