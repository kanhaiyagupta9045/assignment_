package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/car_management/controllers"
	"github.com/kanhaiyagupta9045/car_management/middleware"
)

func CarRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/add/car", middleware.Authenticate(), controllers.AddCar())
		api.GET("/get/car/added_by_user", middleware.Authenticate(), controllers.GetCarAddedByUser())
		api.GET("/search/cars", middleware.Authenticate(), controllers.SearchCarsByKeyowrd()) ///api/v1/search/cars?keyword=BMW
		api.GET("/cars/:car_id", middleware.Authenticate(), controllers.GetCarDetails())
		api.DELETE("/car/:car_id", middleware.Authenticate(), controllers.DeleteCar())
		api.PUT("/update/car/:car_id", middleware.Authenticate(), controllers.UpdateCar())
	}
}
