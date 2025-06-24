package apis

import (
	"boiler-platecode/src/apis/auth"
	"boiler-platecode/src/apis/user"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	UserController *user.UserController
	AuthController *auth.AuthController
}

func NewApiController(userController *user.UserController,authController *auth.AuthController) *ApiController {
	return &ApiController{
		UserController: userController,
		AuthController: authController,
	}
}

func (api *ApiController) RegisterRoutes(router *gin.Engine) {
	api.UserController.InitUserRoutes(router)
	api.AuthController.InitAuthRoutes(router)
}
