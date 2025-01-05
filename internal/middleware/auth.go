package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tihai/global"
	"tihai/utils"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			ctx.Abort()
			return
		}
		// 查询token是否在黑名单中
		_, err := global.RedisDB.Get(token).Result()
		if err == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token失效，重新登陆"})
			ctx.Abort()
			return
		}

		authMap, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("uid", authMap["uid"])
		ctx.Set("role", authMap["role"])
		ctx.Next()
	}
}

func TeacherMiddle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, _ := ctx.Get("role")
		if role != "teacher" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "该操作仅支持老师身份"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
