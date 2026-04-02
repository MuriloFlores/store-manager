package dto

type CreateUserInput struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Roles    []string `json:"roles" binding:"required"`
}
