package ioc

import (
	"fmt"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func InitDB(l logger.Logger) *gorm.DB {

	type Config struct {
		Dsn string `yaml:"dsn"`
	}
	var cfg Config
	cfg = Config{
		Dsn: "root:root@tcp(localhost:3306)/webook",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", cfg, err))
	}
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  glogger.Info,
			ParameterizedQueries:      true,
		}),
	})
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

type gormLoggerFunc func(msg string, feilds ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{
		Key:   "args",
		Value: args,
	})
}
