package project

import (
	routes "ibtools_server/util/routers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers route handlers for the oauth service
func (s *Service) RegisterRoutes(router *gin.RouterGroup, prefix string) {
	oauth := router.Group(prefix)
	{
		routes.AddRoutes(s.GetRoutes(), oauth)
	}
}

// GetRoutes returns []routes.Route slice for the project service
func (s *Service) GetRoutes() []routes.Route {
	return []routes.Route{
		{
			Name:        "CreateProjectByTemplate",
			Method:      "POST",
			Pattern:     "createprojectbytemplate",
			HandlerFunc: s.createPrjByTemplateHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
		{
			Name:        "GetAllTemplates",
			Method:      "POST",
			Pattern:     "getalltemplates",
			HandlerFunc: s.getAllTemplateHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
		{
			Name:        "GetProjectFilePermission",
			Method:      "POST",
			Pattern:     "getprojectfilepermission",
			HandlerFunc: s.getProjectFilePermissionHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
		{
			Name:        "GetMyProjects",
			Method:      "POST",
			Pattern:     "getmyprojects",
			HandlerFunc: s.getMyProjectsHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
		{
			Name:        "GetFullProject",
			Method:      "POST",
			Pattern:     "getfullproject",
			HandlerFunc: s.getFullProjectHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
		{
			Name:        "CreateDDModel",
			Method:      "POST",
			Pattern:     "createDDModel",
			HandlerFunc: s.createDDModelHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
		{
			Name:        "createDDModelByTemplate",
			Method:      "POST",
			Pattern:     "createDDModelByTemplate",
			HandlerFunc: s.createDDModelByTemplateHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
		{
			Name:        "CreateDDItem",
			Method:      "POST",
			Pattern:     "createDDItem",
			HandlerFunc: s.createDDItemHandler,
			Middlewares: []gin.HandlerFunc{s.oauthService.AuthenticateMiddleWare()},
		},
	}
}
