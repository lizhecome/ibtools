package project

import (
	"ibtools_server/config"

	routes "ibtools_server/util/routers"

	"github.com/gin-gonic/gin"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	RestrictToRoles(allowedRoles ...string)
	IsRoleAllowed(role string) bool
	RegisterRoutes(router *gin.RouterGroup, prefix string)
	GetRoutes() []routes.Route
	Close()
}
