package oauth

import (
	"ibtools_server/config"
	"ibtools_server/models"
	routes "ibtools_server/util/routers"

	"github.com/gin-gonic/gin"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	RestrictToRoles(allowedRoles ...string)
	IsRoleAllowed(role string) bool
	FindRoleByID(id uint) (*models.OauthRole, error)
	GetRoutes() []routes.Route
	RegisterRoutes(router *gin.RouterGroup, prefix string)
	UserPhoneExists(phone string) bool
	UserEmailExists(email string) bool
	FindUserByUserID(userid string) (*models.OauthUser, error)
	FindUserByEmail(email string) (*models.OauthUser, error)
	FindUserByPhone(phone string) (*models.OauthUser, error)
	CreateUser(phone, invitationCode, roleName, username, password, company, title, email string) (*models.OauthUser, error)
	SetPassword(user *models.OauthUser, password string) error
	UpdateUseremail(user *models.OauthUser, email string) error
	UpdateUserphone(user *models.OauthUser, phone string) error
	AuthUserByEmail(email, thePassword string) (*models.OauthUser, error)
	AuthUserByPhone(phone, thePassword string) (*models.OauthUser, error)
	Login(user *models.OauthUser, scope string) (*models.OauthAccessToken, *models.OauthRefreshToken, error)
	GrantAccessToken(user *models.OauthUser, expiresIn int, scope string) (*models.OauthAccessToken, error)
	GetOrCreateRefreshToken(user *models.OauthUser, expiresIn int, scope string) (*models.OauthRefreshToken, error)
	GetValidRefreshToken(token string) (*models.OauthRefreshToken, error)
	Authenticate(token string) (*models.OauthAccessToken, error)
	ClearUserTokens(userrefreshToken, useraccessToken string)
	SendPhoneNumValidateMessage(phone string) error
	AuthenticateMiddleWare() gin.HandlerFunc
	Close()
}
