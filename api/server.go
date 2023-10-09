package api

import (
	"fmt"

	_ "github.com/bontusss/colosach/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	fmt.Println("Starting server at", address)
	router := gin.Default()

	router.POST("/pexels", SearchPexel)

	router.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(cors.Default())
	return router.Run(address)
}

func errResponse(err error) interface{} {
	return gin.H{"error": err.Error()}
}
