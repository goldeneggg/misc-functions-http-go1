package usecase

import (
	"context"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/workstatus/adapter"
)

type DefaultUsecase struct {
	gw adapter.Gateway
}

func NewDefaultUsecase(gw adapter.Gateway) workstatus.Usecase {
	return &DefaultUsecase{gw}
}

func (du *DefaultUsecase) Create(ctx context.Context, workstatus *entity.Workstatus) (*entity.Workstatus, error) {
	err := du.gw.Create(ctx, workstatus)
	return workstatus, err
}
