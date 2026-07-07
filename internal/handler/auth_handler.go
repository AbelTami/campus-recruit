package handler

import (
	"employment-server/internal/dto/request"
	dto "employment-server/internal/dto/response"
	"employment-server/internal/service"
	"employment-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 用户登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请提供正确的用户名和密码")
		return
	}

	pair, user, err := h.authService.Login(c.Request.Context(), req.Username, req.Password, c.ClientIP())
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			response.Error(c, response.CodeUnauthorized, "用户名或密码错误")
		case service.ErrAccountLocked:
			response.Error(c, response.CodeUnauthorized, "账号已被锁定，请30分钟后重试")
		case service.ErrAccountDisabled:
			response.Error(c, response.CodeForbidden, "账号已被禁用，请联系管理员")
		default:
			response.InternalError(c, "登录失败，请稍后重试")
		}
		return
	}

	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles { roles[i] = r.Code }

	response.Success(c, dto.LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
		UserInfo: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Roles:    roles,
		},
	})
}

// RefreshToken 刷新 Access Token
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请提供 Refresh Token")
		return
	}

	pair, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Token已失效，请重新登录")
		return
	}

	response.Success(c, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"expires_in":    pair.ExpiresIn,
	})
}

// Logout 退出登录
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint64("user_id")
	_ = h.authService.Logout(c.Request.Context(), userID)
	response.Success(c, nil)
}

// GetUserInfo 获取当前用户信息
// GET /api/v1/auth/userinfo
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetUint64("user_id")
	username := c.GetString("username")
	roles := c.GetStringSlice("user_roles")

	// ponytail: nickname not in JWT claims; fetch from DB if needed later
	response.Success(c, dto.UserInfo{
		ID:       userID,
		Username: username,
		Nickname: username, // fallback
		Roles:    roles,
	})
}
