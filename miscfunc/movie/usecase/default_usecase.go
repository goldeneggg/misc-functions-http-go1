package usecase

import (
	"context"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/movie/adapter"
)

type DefaultUsecase struct {
	gw adapter.Gateway
}

func NewDefaultUsecase(gw adapter.Gateway) movie.Usecase {
	return &DefaultUsecase{gw}
}

func (du *DefaultUsecase) Create(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	// TODO
	return &entity.Movie{}, nil
}
