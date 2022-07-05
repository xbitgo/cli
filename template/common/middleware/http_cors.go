package middleware

import "github.com/gin-gonic/gin"

type CorsConfig struct {
	AllowOrigins string
	Credentials  string
	Headers      string
	Methods      string
}

func HTTPCors(corsCfg ...*CorsConfig) gin.HandlerFunc {
	var (
		origins     = "*"
		credentials = "true"
		headers     = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With"
		methods     = "POST, OPTIONS, GET, PUT, OPTIONS"
	)
	return func(c *gin.Context) {
		if corsCfg != nil && len(corsCfg) == 1 {
			cors := corsCfg[0]
			if cors.Credentials != "" {
				credentials = cors.Credentials
			}
			if cors.AllowOrigins != "" {
				origins = cors.AllowOrigins
			}
			if cors.Headers != "" {
				headers = cors.Headers
			}
			if cors.Methods != "" {
				methods = cors.Methods
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", credentials)
		c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		c.Writer.Header().Set("Access-Control-Allow-Methods", methods)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
