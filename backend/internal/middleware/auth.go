package middleware

import (
	"employment-server/pkg/jwt"
	"employment-server/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证信息")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := jwtService.ParseAccessToken(parts[1])
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Unauthorized(c, "Token已过期，请刷新")
			} else {
				response.Unauthorized(c, "Token无效")
			}
			c.Abort()
			return
		}

		// 注入用户信息到 Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_roles", claims.Roles)

		c.Next()
	}
}
