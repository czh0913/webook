package repository

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
)

var (
	ErrAsyncNotFound = dao.ErrAsyncNotFound
)

type AsyncSMSRepository interface {
	Store(ctx context.Context, biz string, args []string, number string, status string) error
	Find(ctx context.Context) ([]domain.FailedSMS, error)
}

type Async struct {
	AsyncDAO dao.AsyncDAO
}

func NewAsync(as dao.AsyncDAO) AsyncSMSRepository {
	return &Async{
		AsyncDAO: as,
	}
}

func (a Async) Store(ctx context.Context, biz string, args []string, number string, status string) error {
	return a.AsyncDAO.StoreAsync(ctx, biz, args, number, status)
}

func (a Async) Find(ctx context.Context) ([]domain.FailedSMS, error) {
	return a.AsyncDAO.FindAsync(ctx)
}
