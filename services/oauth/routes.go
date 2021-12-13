package oauth

import (
	routes "ibtools_server/util/routers"

	"github.com/gin-gonic/gin"
)

const (
	registerResource                    = "register"
	registerPath                        = "/" + registerResource
	loginHandlerResource                = "login"
	loginHandlerPath                    = "/" + loginHandlerResource
	adminloginHandlerResource           = "adminlogin"
	adminloginHandlerPath               = "/" + adminloginHandlerResource
	forgetpasswordHandlerResource       = "forgetpassword"
	forgetpasswordHandlerPath           = "/" + forgetpasswordHandlerResource
	sendverificationcodeHandlerResource = "sendverificationcode"
	sendverificationcodeHandlerPath     = "/" + sendverificationcodeHandlerResource
	verificationcodeHandlerResource     = "verificationcode"
	verificationcodeHandlerPath         = "/" + verificationcodeHandlerResource
	refreshtokenHandlerResource         = "refreshtoken"
	refreshtokenHandlerPath             = "/" + refreshtokenHandlerResource

	uploaduserimageResource = "uploaduserimage"
	uploaduserimagePath     = "/" + uploaduserimageResource
)

// RegisterRoutes registers route handlers for the oauth service
func (s *Service) RegisterRoutes(router *gin.RouterGroup, prefix string) {
	oauth := router.Group(prefix)
	{
		routes.AddRoutes(s.GetRoutes(), oauth)
	}
}

// GetRoutes returns []routes.Route slice for the oauth service
func (s *Service) GetRoutes() []routes.Route {
	return []routes.Route{
		{
			Name:        "register",
			Method:      "POST",
			Pattern:     registerPath,
			HandlerFunc: s.registerHandler,
		},
		{
			Name:        "login",
			Method:      "POST",
			Pattern:     loginHandlerPath,
			HandlerFunc: s.loginHandler,
		},
		{
			Name:        "refreshtoken",
			Method:      "POST",
			Pattern:     refreshtokenHandlerPath,
			HandlerFunc: s.refreshtokenHandler,
		},
		{
			Name:        "checkphonenumberavailable",
			Method:      "POST",
			Pattern:     "checkphonenumberavailable",
			HandlerFunc: s.checkphonenumberavailableHandler,
		},
		{
			Name:        "forgetpassword",
			Method:      "POST",
			Pattern:     forgetpasswordHandlerPath,
			HandlerFunc: s.forgetpasswordHandler,
		},
		{
			Name:        "sendverificationcode",
			Method:      "POST",
			Pattern:     sendverificationcodeHandlerPath,
			HandlerFunc: s.sendverificationcodeHandler,
		},
		{
			Name:        "createinvitecode",
			Method:      "POST",
			Pattern:     "createinvitecode",
			HandlerFunc: s.createInviteCodeHandler,
			Middlewares: []gin.HandlerFunc{s.AuthenticateMiddleWare()},
		},
		{
			Name:        "createinvitebyphone",
			Method:      "POST",
			Pattern:     "createinvitebyphone",
			HandlerFunc: s.createInviteByPhoneHandler,
			Middlewares: []gin.HandlerFunc{s.AuthenticateMiddleWare()},
		},
		{
			Name:        "createinvitebyusercode",
			Method:      "POST",
			Pattern:     "createinvitebyusercode",
			HandlerFunc: s.createInviteByUserCodeHandler,
			Middlewares: []gin.HandlerFunc{s.AuthenticateMiddleWare()},
		},
		{
			Name:        "processInvite",
			Method:      "POST",
			Pattern:     "processInvite",
			HandlerFunc: s.processInviteHandler,
			Middlewares: []gin.HandlerFunc{s.AuthenticateMiddleWare()},
		},
		{
			Name:        "getInviteListToMe",
			Method:      "POST",
			Pattern:     "getInviteListToMe",
			HandlerFunc: s.getInviteListToMeHandler,
			Middlewares: []gin.HandlerFunc{s.AuthenticateMiddleWare()},
		},
		{
			Name:        "getInviteListFromMe",
			Method:      "POST",
			Pattern:     "getInviteListFromMe",
			HandlerFunc: s.getInviteListFromMeHandler,
			Middlewares: []gin.HandlerFunc{s.AuthenticateMiddleWare()},
		},
	}
}
