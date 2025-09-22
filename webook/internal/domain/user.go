package domain

import (
	"time"
)

type User struct {
	Id int64
	// sql.NullString 可以为nil的数据库变量类型
	Email      string
	Password   string
	Phone      string
	WeChatInfo WeChatInfo

	// UTC 0 的时区
	Ctime time.Time

	//Addr Address
}

//type Address struct {
//	Province string
//	Region   string
//}

//func (u User) ValidateEmail() bool {
// 在这里用正则表达式校验
//return u.Email
//}
