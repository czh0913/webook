package article

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
)

type ArticleReaderRepository interface {
	// 有就更新，没有就创建
	Save(ctx context.Context, art domain.Article) (int64, error)
}
