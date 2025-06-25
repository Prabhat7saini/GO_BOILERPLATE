package user

import (
	"boiler-platecode/src/apis/user/domain"
	"boiler-platecode/src/apis/user/validator"
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/utils"
	"fmt"
	"net/http"

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

func (u *UserController) RegisterPublicRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.POST("/", u.CreateUser)
}

func (u *UserController) RegisterProtectedRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.GET("/profile", u.GetUserProfile)
}

func (u *UserController) RegisterPrivateRoutes(router *gin.RouterGroup) {
	// You can leave this empty for now or add admin-only routes
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


func (u *UserController) GetUserProfile(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	// Type assert to int
	id, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	output := u.userService.GetUserProfile(ctx, id)
	utils.SendRestResponse(ctx, output)
}
