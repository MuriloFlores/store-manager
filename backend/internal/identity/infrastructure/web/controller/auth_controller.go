package controller

import (
	"net/http"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/web/helper"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/web/middleware"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/security"
	"github.com/gin-gonic/gin"
)

const (
	AccessTokenKey       = "access_token"
	RefreshTokenKey      = "refresh_token"
	AccessTokenDuration  = 900
	RefreshTokenDuration = 604800
)

type AuthController struct {
	tokenManager     security.TokenManager
	loginUC          security.LoginUseCase
	refreshTokenUC   security.RotateRefreshTokenUseCase
	forgotPassword   security.ForgotPasswordUseCase
	changePasswordUC security.ChangePasswordUseCase
	logoutUC         security.LogoutUseCase
	rateLimit        ports.RateLimiterRepository
}

func NewLoginController(
	tokenManager security.TokenManager,
	loginUC security.LoginUseCase,
	refreshToken security.RotateRefreshTokenUseCase,
	forgotPassword security.ForgotPasswordUseCase,
	changePasswordUseCase security.ChangePasswordUseCase,
	logout security.LogoutUseCase,
	rateLimit ports.RateLimiterRepository,
) *AuthController {
	return &AuthController{
		tokenManager:     tokenManager,
		loginUC:          loginUC,
		refreshTokenUC:   refreshToken,
		forgotPassword:   forgotPassword,
		changePasswordUC: changePasswordUseCase,
		logoutUC:         logout,
		rateLimit:        rateLimit,
	}
}

func (h *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	authRoutes := router.Group("/auth")
	authRoutes.Use(middleware.RateLimit(h.rateLimit))
	{
		authRoutes.POST("/login", h.Login)
		authRoutes.GET("/refresh", h.RefreshToken)
		authRoutes.POST("/forgot-password", h.ForgotPassword)
	}

	authPrivates := authRoutes.Group("/private/auth")
	authPrivates.Use(middleware.RequireAuth(h.tokenManager))
	{
		authPrivates.GET("/logout", h.Logout)
		authPrivates.POST("/change-password", h.ChangePassword)
	}
}

// Login handles user authentication
// @Summary User Login
// @Description Authenticates user and sets session cookies
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} map[string]string "message: login successfully"
// @Failure 400 {object} map[string]string "error: invalid input"
// @Failure 401 {object} map[string]string "error: invalid credentials"
// @Router /auth/login [post]
func (h *AuthController) Login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResult, err := h.loginUC.Execute(c.Request.Context(), &input)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.SetCookie(AccessTokenKey, loginResult.AccessToken, AccessTokenDuration, "/", "", false, true)
	c.SetCookie(RefreshTokenKey, loginResult.RefreshToken, RefreshTokenDuration, "/auth/refresh", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "login successfully"})
}

// RefreshToken rotates tokens using refresh cookie
// @Summary Refresh Tokens
// @Description Rotates Access and Refresh tokens
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string "message: refresh token successfully"
// @Failure 400 {object} map[string]string "error: refresh_token cookie not found"
// @Failure 401 {object} map[string]string "error: invalid session"
// @Router /auth/refresh [get]
func (h *AuthController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token cookie not found"})
		return
	}

	tokens, err := h.refreshTokenUC.Execute(c.Request.Context(), refreshToken)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.SetCookie(AccessTokenKey, tokens.AccessToken, AccessTokenDuration, "/", "", false, true)
	c.SetCookie(RefreshTokenKey, tokens.RefreshToken, RefreshTokenDuration, "/auth/refresh", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "refresh token successfully"})
}

// Logout clears session cookies
// @Summary User Logout
// @Description Logs out the user and clears session cookies
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "message: logout successfully"
// @Router /private/auth/logout [get]
func (h *AuthController) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(RefreshTokenKey)
	if err == nil {
		if err := h.logoutUC.Execute(c.Request.Context(), refreshToken); err != nil {
			helper.HandleError(c, err)
			return
		}
	}

	c.SetCookie(AccessTokenKey, "", -1, "/", "", false, true)
	c.SetCookie(RefreshTokenKey, "", -1, "/auth/refresh", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logout successfully"})
}

// ForgotPassword requests a password reset email
// @Summary Forgot Password
// @Description Sends a password reset email to the provided address
// @Tags Auth
// @Accept json
// @Produce json
// @Param forgotPasswordInput body object{email=string} true "Email for password reset"
// @Success 200 {object} map[string]string "message: email sent successfully"
// @Failure 400 {object} map[string]string "error: invalid email"
// @Router /auth/forgot-password [post]
func (h *AuthController) ForgotPassword(c *gin.Context) {
	var forgotPasswordInput struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBind(&forgotPasswordInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPassword.Execute(c.Request.Context(), forgotPasswordInput.Email); err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email sent successfully"})
}

// ChangePassword changes user's password
// @Summary Change Password
// @Description Updates user's password with a new one
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param changePasswordInput body object{old_password=string,new_password=string} true "Old and new password"
// @Success 200 {object} map[string]string "message: password successfully changed"
// @Failure 400 {object} map[string]string "error: invalid input"
// @Failure 401 {object} map[string]string "error: unauthorized or invalid old password"
// @Router /private/auth/change-password [post]
func (h *AuthController) ChangePassword(c *gin.Context) {
	var changePasswordInput struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBind(&changePasswordInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := helper.ExtractUserClaims(c.Request.Context())
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	err = h.changePasswordUC.Execute(c.Request.Context(), claims.UserID, changePasswordInput.OldPassword, changePasswordInput.NewPassword)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password successfully changed"})
}
