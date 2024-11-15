package middleware

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/car_management/helpers"
)

func Middleware(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "Upgrade", "Connection"},
		AllowCredentials: true,
	}))
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(200)
	})
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helpers.GetUserFromCookie(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("user", user)

		c.Next()
	}
}
