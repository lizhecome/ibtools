package oauth

import (
	"ibtools_server/drerror"

	"github.com/gin-gonic/gin"
)

// AuthenticateMiddleWare 认证中间件
func (s *Service) AuthenticateMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		at := c.GetHeader("at")
		if at == "" {
			drerror.ResponseError(c, drerror.APIErrAccessTokenIsNull, nil)
			c.Abort()
		}

		if accesstoken, err := s.Authenticate(at); err != nil {
			drerror.ResponseError(c, drerror.APIErrAuthFailed, nil)
			c.Abort()
		} else {
			c.Set("user", accesstoken.User)
			c.Next()
		}

	}
}
