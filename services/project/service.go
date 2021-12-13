package project

import (
	"ibtools_server/config"
	"ibtools_server/services/oauth"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
	oauthService oauth.ServiceInterface
	stsClient    *sts.Client
	ossClient    *oss.Client
}

// NewService returns a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, cache *cache.Cache, logger *logrus.Logger, oauthService oauth.ServiceInterface) *Service {

	stsclient, _ := sts.NewClientWithAccessKey(cnf.Aliyun.RegionId, cnf.Aliyun.AccessKeyId, cnf.Aliyun.AccessKeySecret)
	ossclient, _ := oss.New(cnf.Aliyun.OssEndPoint, cnf.Aliyun.AccessKeyId, cnf.Aliyun.AccessKeySecret)
	return &Service{
		cnf:          cnf,
		db:           db,
		cache:        cache,
		allowedRoles: []string{},
		logger:       logger,
		oauthService: oauthService,
		stsClient:    stsclient,
		ossClient:    ossclient,
	}
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
