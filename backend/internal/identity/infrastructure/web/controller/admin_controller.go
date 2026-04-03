package controller

import (
	"net/http"
	"strconv"

	"github.com/MuriloFlores/order-manager/internal/common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/web/helper"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/web/middleware"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/admin"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	getUserInfo  admin.GetUsersInfo
	changeStatus admin.ChangeUserStatusUseCase
	changeRole   admin.ChangeUserRoleUseCase
	tokenManager auth.TokenManager
}

func NewAdminController(getUserInfo admin.GetUsersInfo) *AdminController {
	return &AdminController{
		getUserInfo: getUserInfo,
	}
}

func (h *AdminController) RegisterRoutes(engine *gin.Engine) {
	allowedRoles := []vo.Role{
		vo.ManagerRole,
		vo.AdminRole,
	}

	adminRoutes := engine.Group("/admin")
	adminRoutes.Use(middleware.RequireAuth(h.tokenManager))
	adminRoutes.Use(middleware.VerifyRole(allowedRoles...))
	{
		adminRoutes.GET("/users", h.GetUsersInfo)
		adminRoutes.PATCH("/:id/status", h.ChangeUserStatus)
		adminRoutes.PUT("/:id/roles", h.ChangeUserRoles)
	}
}

// GetUsersInfo returns paginated user information
// @Summary List Users
// @Description Returns a paginated list of users filtered by role
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param roles query []string false "Filter by roles"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Items per page (default 10)"
// @Param search query string false "Search query"
// @Param sort query string false "Sort field (default name)"
// @Param direction query string false "Sort direction (ASC/DESC, default DESC)"
// @Success 200 {object} _common.PaginatedResult[entity.User] "Paginated user data"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Router /admin/users [get]
func (h *AdminController) GetUsersInfo(c *gin.Context) {
	roles := c.QueryArray("roles")

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	search := c.DefaultQuery("search", "")
	sort := c.DefaultQuery("sort", "name")
	direction := c.DefaultQuery("direction", "DESC")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	pagination := common.NewPagination(page, pageSize, search, sort, direction)

	info, err := h.getUserInfo.Execute(c.Request.Context(), pagination, roles)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

// ChangeUserStatus activates or deactivates a user
// @Summary Change User Status
// @Description Activates or deactivates a user by ID
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param status body object{status=bool} true "New status"
// @Success 200 {object} map[string]string "status: status updated"
// @Failure 400 {object} map[string]string "error: invalid input"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Failure 404 {object} map[string]string "error: user not found"
// @Router /admin/{id}/status [patch]
func (h *AdminController) ChangeUserStatus(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Status bool `json:"status" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.changeStatus.Execute(c.Request.Context(), id, input.Status); err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "status updated"})
}

// ChangeUserRoles updates user's roles
// @Summary Change User Roles
// @Description Updates the list of roles for a user by ID
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param roles body object{roles=[]string} true "New roles list"
// @Success 200 {object} map[string]string "roles: roles updated"
// @Failure 400 {object} map[string]string "error: invalid input"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Failure 404 {object} map[string]string "error: user not found"
// @Router /admin/{id}/roles [put]
func (h *AdminController) ChangeUserRoles(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Roles []string `json:"roles" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.changeRole.Execute(c.Request.Context(), id, input.Roles); err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": "roles updated"})
}
