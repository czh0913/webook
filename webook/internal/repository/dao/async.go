package dao

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"gorm.io/gorm"
	"time"
)

var (
	ErrAsyncNotFound = gorm.ErrRecordNotFound
)

type AsyncDAO interface {
	StoreAsync(ctx context.Context, biz string, args []string, number string, status string) error
	FindAsync(ctx context.Context) ([]domain.FailedSMS, error)
}

type Async struct {
	db *gorm.DB
}

func NewAsyncDAO(db *gorm.DB) AsyncDAO {
	return &Async{
		db: db,
	}
}

func (a Async) StoreAsync(ctx context.Context, biz string, args []string, number string, status string) error {
	//一个手机号一个存储
	failedSMS := FailedSMS{}
	var argsErr error
	if argsErr = a.db.WithContext(ctx).Where("number=? AND biz=? AND args=?", number, biz, args).First(&failedSMS).Error; argsErr != nil && !errors.Is(argsErr, ErrAsyncNotFound) {
		return argsErr
	}
	if !errors.Is(argsErr, gorm.ErrRecordNotFound) && failedSMS.Status != status {
		err := a.db.Where("number=? AND biz=? AND args=?", number, biz, args).Update("status=?", status).Error
		if err != nil {
			return err
		}
	}

	// 说明这个请求是新的请求

	failedSMS = FailedSMS{
		Biz:        biz,
		Args:       args,
		Number:     number,
		Status:     "Pending",
		Ctime:      time.Now(),
		UpdateTime: time.Now(),
	}
	if err := a.db.Create(&failedSMS).Error; err != nil {
		return err
	}

	return nil
}

func (a Async) FindAsync(ctx context.Context) ([]domain.FailedSMS, error) {
	var failedSMS []domain.FailedSMS
	err := a.db.Where("status=?", "Pending").Find(&failedSMS).Error
	if err != nil {
		return make([]domain.FailedSMS, 0), err
	}
	return failedSMS, nil
}

type FailedSMS struct {
	ID     int64 `gorm:"primaryKey"`
	Biz    string
	Args   []string
	Number string
	Status string // Waiting Success Failed

	Ctime      time.Time
	UpdateTime time.Time
}
