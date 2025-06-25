package apis

import (
	"boiler-platecode/src/apis/auth"
	"boiler-platecode/src/apis/user"
	middleware "boiler-platecode/src/common/middlewares"
	"boiler-platecode/src/core/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RouteRegistrar interface {
	RegisterPublicRoutes(router *gin.RouterGroup)
	RegisterProtectedRoutes(router *gin.RouterGroup)
	RegisterPrivateRoutes(router *gin.RouterGroup)
}

type ApiController struct {
	registrars []RouteRegistrar
	redisService redis.RedisService
}

func NewApiController(redisService redis.RedisService, registrars ...RouteRegistrar) *ApiController {
	return &ApiController{
		redisService: redisService,
		registrars:   registrars,
	}
}
func InitApiController(db *gorm.DB, redisService redis.RedisService) *ApiController {
	userController := user.InitUserController(db)
	authController := auth.InitAuthController(db, redisService)

	return NewApiController(redisService, userController, authController)
}

func (api *ApiController) RegisterRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")

	public := apiV1.Group("/")
	protected := apiV1.Group("/")
	private := apiV1.Group("/admin")

	// Apply auth middlewares
	protected.Use(middleware.AuthMiddleware(api.redisService))
	// private.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())

	// Register routes from all controllers
	for _, registrar := range api.registrars {
		registrar.RegisterPublicRoutes(public)
		registrar.RegisterProtectedRoutes(protected)
		registrar.RegisterPrivateRoutes(private)
	}
}
