package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (g GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error

	return art.ID, err
}

type Article struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"type=varchar(1024)"`
	Content  string `gorm:"BLOB"`
	AuthorID int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}
