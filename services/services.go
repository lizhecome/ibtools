package services

import (
	"ibtools_server/config"
	"ibtools_server/services/oauth"
	"ibtools_server/services/project"
	"reflect"

	//"github.com/coreos/etcd/proxy/grpcproxy/cache"
	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func init() {

}

var (
	// OauthService ...
	OauthService   oauth.ServiceInterface
	ProjectService project.ServiceInterface
)

// UseOauthService sets the oAuth service
func UseOauthService(o oauth.ServiceInterface) {
	OauthService = o
}

// UseOauthService sets the oAuth service
func UseProjectService(o project.ServiceInterface) {
	ProjectService = o
}

// Init starts up all services
func Init(cnf *config.Config, db *gorm.DB, cache *cache.Cache, logger *logrus.Logger, redis *redis.Client) error {
	// if nil == reflect.TypeOf(HealthService) {
	// 	HealthService = health.NewService(db)
	// }
	var err error
	if nil == reflect.TypeOf(OauthService) {
		OauthService = oauth.NewService(cnf, db, cache, logger, redis)
	}

	if nil == reflect.TypeOf(ProjectService) {
		ProjectService = project.NewService(cnf, db, cache, logger, OauthService)
	}

	return err
}

// Close closes any open services
func Close() {
	OauthService.Close()
	ProjectService.Close()
}
