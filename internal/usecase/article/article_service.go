package article

import (
	"context"
	sv "github.com/core-go/core"
)

type ArticleService interface {
	Load(ctx context.Context, id string) (*Article, error)
}

func NewArticleService(repository sv.Repository) ArticleService {
	return &articleService{repository: repository}
}

type articleService struct {
	repository sv.Repository
}

func (s *articleService) Load(ctx context.Context, id string) (*Article, error) {
	var Article Article
	ok, err := s.repository.Get(ctx, id, &Article)
	if !ok {
		return nil, err
	} else {
		return &Article, err
	}
}
