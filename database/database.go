package database

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"ibtools_server/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Drivers
	_ "github.com/lib/pq"
)

// NewDatabase returns a gorm.DB struct, gorm.DB.DB() returns a database handle
// see http://golang.org/pkg/database/sql/#DB
func NewDatabase(cnf *config.Config) (*gorm.DB, error) {
	// Postgres
	if cnf.Database.Type == "postgres" {
		// Connection args
		// see https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
		args := fmt.Sprintf(
			"sslmode=disable host=%s port=%d user=%s password='%s' dbname=%s",
			cnf.Database.Host,
			cnf.Database.Port,
			cnf.Database.User,
			cnf.Database.Password,
			cnf.Database.DatabaseName,
		)

		db, err := gorm.Open(postgres.Open(args), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			Logger:                                   initLogger(cnf),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			return db, err
		}
		// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		// Max idle connections
		sqlDB.SetMaxIdleConns(cnf.Database.MaxIdleConns)

		// Max open connections
		sqlDB.SetMaxOpenConns(cnf.Database.MaxOpenConns)

		return db, nil
	}

	// Database type not supported
	return nil, fmt.Errorf("Database type %s not suppported", cnf.Database.Type)
}

func initLogger(cnf *config.Config) logger.Interface {
	if !cnf.HasDBLog {
		return nil
	}
	var out *os.File
	var level logger.LogLevel
	var Colorful bool
	if cnf.IsDevelopment {
		out = os.Stdout
		level = logger.Info
		Colorful = true
	} else {
		if err := os.MkdirAll(cnf.LogLocation, os.ModePerm); err != nil {
			fmt.Println("err", err)
		}
		// 日志文件
		fileName := path.Join(cnf.LogLocation, "database.log")
		// 写入文件
		src, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("err", err)
		}
		out = src
		level = logger.Silent
		Colorful = false
	}

	newLogger := logger.New(
		log.New(out, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  level,       // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  Colorful,    // 禁用彩色打印
		},
	)
	return newLogger
}
