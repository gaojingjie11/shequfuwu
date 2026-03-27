package middleware

import (
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireRole 权限校验中间件
// roles: 允许访问的角色列表 (只要满足其中一个即可)
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取当前用户角色 (由 JWT 中间件注入)
		userRole, exists := c.Get("role")
		if !exists {
			response.FailWithCode(c, 403, "无权限: 未获取到角色信息")
			c.Abort()
			return
		}

		roleStr := userRole.(string)

		// 2. 检查是否匹配
		for _, role := range roles {
			if role == roleStr {
				c.Next()
				return
			}
		}

		// 3. 都不匹配 (如果是 admin 角色，通常拥有所有权限，这里简单处理，必须在 roles 列表中)
		// 如果需要 "admin" 角色能访问所有接口，可以在这里加: if roleStr == "admin" { c.Next(); return }

		response.FailWithCode(c, 403, "无权限访问此资源")
		c.Abort()
	}
}
