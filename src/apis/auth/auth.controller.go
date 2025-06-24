package auth

import (
	"boiler-platecode/src/apis/auth/domain"
	"boiler-platecode/src/apis/auth/validator"
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (auth *AuthController) InitAuthRoutes(router *gin.Engine) {
	auths := router.Group("/auth")
	auths.POST("/login", auth.Login)
}

func (auth *AuthController) Login(ctx *gin.Context) {
	var req validator.Login

	if err := ctx.ShouldBindJSON(&req); err != nil {
		exc, _ := utils.CreateUserValidationErrors(err)

		resp := common.ServiceOutput[*domain.LoginResponse]{
			Exception: exc,
		}
		utils.SendRestResponse(ctx, resp)
		return
	}

	reqData := &domain.Login{
		Email:    req.Email,
		Password: req.PASSWORD,
	}

	output := auth.authService.Login(ctx, reqData)
	utils.SendRestResponse(ctx, output)

}
