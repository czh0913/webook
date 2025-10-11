package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateByIdWithAuthor(ctx context.Context, art Article) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (g *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error

	return art.Id, err
}

func (g *GORMArticleDAO) UpdateByIdWithAuthor(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	// 指定要更新的字段
	res := g.db.WithContext(ctx).Model(&art).Where("id=? AND author_id=?", art.Id, art.AuthorID).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，创作者非法， id: %d author_id %d ",
			art.Id, art.AuthorID)
	}

	return res.Error
}

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"type=varchar(1024)"`
	Content  string `gorm:"BLOB"`
	AuthorID int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}
