package service

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/article"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	//repo repository.ArticleRepository
	author article.ArticleAutherRepository
	reader article.ArticleReaderRepository
	l      logger.Logger
}

func NewArticleService(author article.ArticleAutherRepository, reader article.ArticleReaderRepository, l logger.Logger) ArticleService {
	return articleService{
		reader: reader,
		author: author,
		l:      l,
	}
}

func (a articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	id := art.Id
	var err error
	if id > 0 {
		err = a.author.Update(ctx, art)
		if err != nil {
			a.l.Error("更新制作库失败")
		}

	} else {
		id, err = a.author.Creat(ctx, art)
		if err != nil {
			a.l.Error("新建制作库失败")
		}
	}
	if err != nil {
		return 0, err
	}

	art.Id = id

	for i := 0; i < 3; i++ {
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败，保存到线上库失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败，保存到线上库重试彻底失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}

	return id, err
}

func (a articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id != 0 {
		err := a.author.Update(ctx, art)
		return art.Id, err
	}

	return a.author.Creat(ctx, art)
}
