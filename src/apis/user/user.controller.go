package user

import (
	"boiler-platecode/src/apis/user/domain"
	"boiler-platecode/src/apis/user/validator"
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService UserService
}

// Use constructor-style init
func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

// Register Gin routes here
func (u *UserController) InitUserRoutes(router *gin.Engine) {

	users := router.Group("/users")
	users.POST("/", u.CreateUser)
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	fmt.Println("Initializing user routes")

	var req validator.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		exc, _ := utils.CreateUserValidationErrors(err)

		resp := common.ServiceOutput[*domain.User]{
			Exception: exc,
		}
		utils.SendRestResponse(ctx, resp)
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.PASSWORD,
	}

	output := u.userService.CreateUser(ctx, user)
	utils.SendRestResponse(ctx, output)
}
