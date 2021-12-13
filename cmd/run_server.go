package cmd

import (
	"fmt"
	"ibtools_server/config"
	"ibtools_server/middleware"
	"ibtools_server/services"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RunServer runs the app
func RunServer(configBackend string) error {
	cnf, db, cache, log, redis, err := initConfigDB(true, true, configBackend)
	if err != nil {
		return err
	}

	// start the services
	if err := services.Init(cnf, db, cache, log, redis); err != nil {
		return err
	}
	defer services.Close()

	//TODO config support
	if cnf.IsDevelopment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	logger := initAccessLogger(cnf)
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	// Version: v1
	v1 := router.Group("/v1")
	{
		services.OauthService.RegisterRoutes(v1, "oauth")
		services.ProjectService.RegisterRoutes(v1, "project")
	}
	url := ginSwagger.URL(cnf.SwaggerURL) // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	log.Info("server running")
	router.Run(":8080")
	return nil
}

func initAccessLogger(cnf *config.Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	if !cnf.IsDevelopment {
		if err := os.MkdirAll(cnf.LogLocation, os.ModePerm); err != nil {
			fmt.Println("err", err)
		}
		// 日志文件
		fileName := path.Join(cnf.LogLocation, "access.log")
		// 写入文件
		src, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("err", err)
		}
		logger.Out = src
	}

	return logger
}
