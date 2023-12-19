package routes

import (
	"github.com/airchains-network/da-client/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/celestia", controllers.CelestiaController)
	r.POST("/avail", controllers.AvailController)

	return r
}
