package routes

import (
	"github.com/gin-gonic/gin"
)

// Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

// AddRoutes adds routes to a router instance. If there are middlewares defined
// for a route, a new negroni app is created and wrapped as a http.Handler
func AddRoutes(routes []Route, router *gin.RouterGroup) {

	for _, route := range routes {
		// Add any specified middlewares
		if len(route.Middlewares) > 0 {

			for _, middleware := range route.Middlewares {
				router.Use(middleware)
			}

			// Wrap the handler in the negroni app with middlewares
			//router.Use(route.HandlerFunc)
		}
		router.Handle(route.Method, route.Pattern, route.HandlerFunc)
	}
}
