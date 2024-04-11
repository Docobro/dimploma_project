package middleware

import "github.com/gin-gonic/gin"

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from all origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow all the HTTP methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		// Allow all the headers
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		// Allow credentials (cookies, authorization headers, etc.)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// Move to the next middleware or handler
		c.Next()
	}
}
