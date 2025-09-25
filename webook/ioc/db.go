package ioc

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	type Config struct {
		Dsn string `yaml:"dsn"`
	}
	var cfg Config
	cfg = Config{
		Dsn: "root:root@tcp(localhost:3306)/webook",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Dsn))
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		//只会在初始化过程种 panic
		//整个 goroutine 结束
		//一旦初始化过程出错，应用就不要启动了
		panic(err)
		//直接退出
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
