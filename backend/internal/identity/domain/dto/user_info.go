package dto

type UserInfo struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     []string `json:"role"`
}
