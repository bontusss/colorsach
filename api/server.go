package api

import (
	_ "github.com/bontusss/colosach/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-contrib/cors"
)

// @title Colosach API

// @BasePath /:3000

// colosach pexels godoc
// @Summary colosach pexels
// @param name string
//
// @Description Get an image from pexels
// @Tags example
// @Accept json
// @Produce json
// @Success 200
// @Router /pexel [post]
func Start(address string) error {
	router := gin.Default()

	router.POST("/pexel", SearchPexel)

	router.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(cors.Default())
	return router.Run(address)
}

func errResponse(err error) interface{} {
	return gin.H{"error": err.Error()}
}
