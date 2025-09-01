package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConfilictsErrNo uint16 = 1062
		if mysqlErr.Number == 1062 {
			//邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, u User) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", u.Email).First(&user).Error
	return user, err
}

// User 直接对应数据库表结构
// 有些人叫 entity 有些人叫 model  有些人叫 PO (persistent object)
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// Email 全部用户唯一
	Email    string `gorm:"unique"`
	Password string

	//创建时间，毫秒
	Ctime int64

	//更新时间，毫秒
	Utime int64
}
