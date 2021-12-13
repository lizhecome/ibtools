package oauth

import (
	"ibtools_server/config"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Service struct keeps objects to avoid passing them around
type Service struct {
	cnf          *config.Config
	db           *gorm.DB
	cache        *cache.Cache
	allowedRoles []string
	logger       *logrus.Logger
	redis        *redis.Client
	smsClient    *dysmsapi20170525.Client
}

// NewService returns a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, cache *cache.Cache, logger *logrus.Logger, redis *redis.Client) *Service {
	client, _ := CreateClient(&cnf.Aliyun.AccessKeyId, &cnf.Aliyun.AccessKeySecret)
	// err := db.SetupJoinTable(&models.OauthUser{}, "Projects", &models.ProjectUser{})
	// db.SetupJoinTable(&models.Project{}, "Users", &models.ProjectUser{})
	// if err != nil {
	// 	return nil
	// }

	return &Service{
		cnf:          cnf,
		db:           db,
		cache:        cache,
		allowedRoles: []string{},
		logger:       logger,
		redis:        redis,
		smsClient:    client,
	}
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// GetConfig returns config.Config instance
func (s *Service) GetConfig() *config.Config {
	return s.cnf
}

// RestrictToRoles restricts this service to only specified roles
func (s *Service) RestrictToRoles(allowedRoles ...string) {
	s.allowedRoles = allowedRoles
}

// IsRoleAllowed returns true if the role is allowed to use this service
func (s *Service) IsRoleAllowed(role string) bool {
	for _, allowedRole := range s.allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}

// Close stops any running services
func (s *Service) Close() {}
