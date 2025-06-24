package validator


type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=3"`
	Email string `json:"email" binding:"required,email"`
	PASSWORD string `json:"password" binding:"required,strongpassword"`
}
