package apis

import (

	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterPublicRoutes(router *gin.RouterGroup)
	RegisterProtectedRoutes(router *gin.RouterGroup)
	RegisterPrivateRoutes(router *gin.RouterGroup)
}

type ApiController struct {
	registrars []RouteRegistrar
}


func NewApiController(registrars ...RouteRegistrar) *ApiController {
	return &ApiController{
		registrars: registrars,
	}
}



func (api *ApiController) RegisterRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")

	public := apiV1.Group("/")
	protected := apiV1.Group("/")
	private := apiV1.Group("/admin")

	// Apply auth middlewares
	// protected.Use(middlewares.AuthMiddleware())
	// private.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())

	// Register routes from all controllers
	for _, registrar := range api.registrars {
		registrar.RegisterPublicRoutes(public)
		registrar.RegisterProtectedRoutes(protected)
		registrar.RegisterPrivateRoutes(private)
	}
}