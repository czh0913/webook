package article

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
)

type ArticleAutherRepository interface {
	Creat(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
