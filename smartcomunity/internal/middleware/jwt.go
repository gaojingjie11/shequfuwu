package middleware

import (
	"context"
	"fmt"
	"smartcommunity/internal/global"
	"smartcommunity/pkg/response"
	"smartcommunity/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Header 中的 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FailWithCode(c, 401, "请先登录")
			c.Abort() // 阻止后续处理
			return
		}

		// 2. 格式校验 (通常是 "Bearer <token>")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.FailWithCode(c, 401, "Token格式错误")
			c.Abort()
			return
		}

		//// 3. 解析 Token
		//claims, err := utils.ParseToken(parts[1])
		//if err != nil {
		//	response.FailWithCode(c, 401, "Token无效或已过期")
		//	c.Abort()
		//	return
		//}
		//
		//// 4. 将当前用户ID存入 Context，供后续 Handler 使用
		//// 【关键】后续在 Handler 里用 c.GetInt64("userID") 取出来
		//c.Set("userID", claims.UserID)
		//c.Set("role", claims.Role)
		// 3. 解析 Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			response.FailWithCode(c, 401, "Token无效或已过期")
			c.Abort()
			return
		}

		// --- [新增 Redis 校验逻辑] ---
		// 1. 拼装 Key
		redisKey := fmt.Sprintf("login:token:%d", claims.UserID)

		// 2. 从 Redis 获取该用户当前的 Token
		cachedToken, err := global.RDB.Get(context.Background(), redisKey).Result()

		// 3. 校验
		// 情况A: Redis里没数据 (可能过期了，或者用户已注销)
		// 情况B: Redis里的Token 和 请求头里的Token 不一致 (说明用户在别处登录了)
		if err != nil || cachedToken != parts[1] {
			response.FailWithCode(c, 401, "登录已失效，请重新登录")
			c.Abort()
			return
		}
		// ---------------------------

		// 4. 将当前用户ID存入 Context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next() // 放行
	}
}
