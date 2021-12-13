package cmd

import (
	"fmt"
	"ibtools_server/config"
	"ibtools_server/database"
	"os"
	"path"
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// initConfigDB loads the configuration and connects to the database
func initConfigDB(mustLoadOnce, keepReloading bool, configBackend string) (*config.Config, *gorm.DB, *cache.Cache, *logrus.Logger, *redis.Client, error) {
	// Config
	cnf := config.NewConfig(mustLoadOnce, keepReloading, configBackend)
	c := cache.New(5*time.Minute, 10*time.Minute)
	// Database
	db, err := database.NewDatabase(cnf)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     cnf.Redis.Addr,
	// 	Password: cnf.Redis.Password,
	// 	DB:       cnf.Redis.DB,
	// })

	// pong, err := rdb.Ping().Result()
	// if err != nil {
	// 	return nil, nil, nil, nil, nil, err
	// }
	// fmt.Println(pong, err)

	return cnf, db, c, initApplicationLogger(cnf), nil, nil
}

func initApplicationLogger(cnf *config.Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	if !cnf.IsDevelopment {
		if err := os.MkdirAll(cnf.LogLocation, os.ModePerm); err != nil {
			fmt.Println("err", err)
		}
		//日志文件
		fileName := path.Join(cnf.LogLocation, "server.log")
		//写入文件
		src, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("err", err)
		}
		logger.Out = src
	}

	return logger
}
