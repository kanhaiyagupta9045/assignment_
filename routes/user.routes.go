package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/car_management/controllers"
)

func UserRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{

		api.POST("/register/user", controllers.SignUpController())
		api.POST("/signin/user", controllers.Login())
	}
}
