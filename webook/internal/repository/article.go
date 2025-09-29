package repository

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Creat(ctx context.Context, art domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	da dao.ArticleDAO
}

func NewCachedArticleRepository(da dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		da: da,
	}
}

func (c CachedArticleRepository) Creat(ctx context.Context, art domain.Article) (int64, error) {
	return c.da.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorID: art.Author.Id,
	})
}
