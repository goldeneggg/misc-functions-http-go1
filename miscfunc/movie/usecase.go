package movie

import (
	"context"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/entity"
)

type Usecase interface {
	Create(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
}
