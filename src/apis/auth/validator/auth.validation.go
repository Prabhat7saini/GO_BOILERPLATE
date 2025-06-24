package validator
type Login struct {
	Email    string `json:"email" binding:"required,email"`
	PASSWORD string `json:"password" binding:"required"`
}