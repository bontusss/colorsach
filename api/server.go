package api

import (

	"github.com/gin-gonic/gin"
)

func Start(address string) error {
	router := gin.Default()

	router.POST("/pexel", searchPexel)
	// router.POST()

	return router.Run(address)
}

func errResponse(err error) interface{} {
	return gin.H{ "error": err.Error()}
}