package api

import (
	"postcodes/area"
	"postcodes/service/postcodesio"

	"github.com/gin-gonic/gin"
)

func New(areas *area.Areas) *gin.Engine {
	router := gin.New()
	api := api{areas: areas, api: postcodesio.Client()}
	v1 := router.Group("postcodes/v1")

	v1.POST("/postcodes", api.Postcodes)

	router.GET("/postcodes/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
