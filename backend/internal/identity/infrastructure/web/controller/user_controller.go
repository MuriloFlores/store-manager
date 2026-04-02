package controller

import (
	"net/http"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/web/helper"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/user"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/web/middleware"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	createUserUC user.CreateUserUseCase
	getMyInfoUC  user.MyInfoUseCase
	tokenManager auth.TokenManager
}

func NewUserHandle(createUser user.CreateUserUseCase, getMyInfoUC user.MyInfoUseCase, tokenManager auth.TokenManager) *UserController {
	return &UserController{
		createUserUC: createUser,
		getMyInfoUC:  getMyInfoUC,
		tokenManager: tokenManager,
	}
}

func (h *UserController) RegisterRoutes(router *gin.Engine) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/", h.CreateUser)
	}

	privateRoutes := router.Group("/private/user")
	privateRoutes.Use(middleware.RequireAuth(h.tokenManager))
	{
		privateRoutes.GET("/me", h.MyInfo)
	}
}

// CreateUser creates a new user
// @Summary Create User
// @Description Creates a new user with the provided information
// @Tags User
// @Accept json
// @Produce json
// @Param createUserInput body dto.CreateUserInput true "User information"
// @Success 201 {object} map[string]string "message: user created successfully"
// @Failure 400 {object} map[string]string "error: invalid input"
// @Failure 422 {object} map[string]string "error: validation failed"
// @Router /user [post]
func (h *UserController) CreateUser(c *gin.Context) {
	var input dto.CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.createUserUC.Execute(c.Request.Context(), input); err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

// MyInfo returns the current authenticated user's information
// @Summary Get Current User Info
// @Description Returns the profile information of the currently authenticated user
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserInfo "Current user information"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Failure 404 {object} map[string]string "error: user not found"
// @Router /private/user/me [get]
func (h *UserController) MyInfo(c *gin.Context) {
	claims, err := helper.ExtractUserClaims(c.Request.Context())
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	userData, err := h.getMyInfoUC.Execute(c.Request.Context(), claims.UserID)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, userData)
}
