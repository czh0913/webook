package retry

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/article"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
	"time"
)

type ArticleReaderRepoRetry struct {
	repo   article.ArticleReaderRepository
	l      logger.Logger
	maxTry int
	delay  time.Duration
}

func NewArticleReaderRepoRetry(repo article.ArticleReaderRepository, l logger.Logger, maxTry int, delay time.Duration) article.ArticleReaderRepository {
	return &ArticleReaderRepoRetry{
		repo:   repo,
		l:      l,
		maxTry: maxTry,
		delay:  delay,
	}
}

func (a ArticleReaderRepoRetry) Save(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  int64
		err error
	)

	for i := 0; i < a.maxTry; i++ {
		id, err = a.repo.Save(ctx, art)
		if err == nil {
			return id, err
		}

		a.l.Error("部分失败，Save 重试失败",
			logger.Int64("art_id", art.Id),
			logger.Int64("重试次数", int64(i)),
			logger.Error(err))

		if i < a.maxTry-1 {
			time.Sleep(a.delay)
		}

	}

	a.l.Error("部分失败，Save 重试全部失败",
		logger.Int64("art_id", art.Id),
		logger.Error(err))

	return id, err

}
