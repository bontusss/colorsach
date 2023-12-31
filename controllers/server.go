package controllers

import (
	"fmt"

	_ "github.com/bontusss/colosach/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title Colosach API

// @BasePath /:3000

// Start colosach pexels godoc
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
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AddAllowMethods("OPTIONS")
	router.Use(cors.New(corsConfig))

	router.POST("/pexels", SearchPexel)
	router.POST("/unsplash", GetUnsplash)

	//router.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router.Run(address)
}

func errResponse(err error) interface{} {
	return gin.H{"error": err.Error()}
}
