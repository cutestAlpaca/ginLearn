package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if context.Request.Method == http.MethodOptions {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
