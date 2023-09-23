package api

import (
	"context"
	"log"
	"net/http"

	"github.com/bontusss/colosach/utils"
	"github.com/gin-gonic/gin"
	"github.com/kosa3/pexels-go"
)

type searchRequest struct {
	Color string `json:"color" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// Gin handler for getting image from pexels api
func searchPexel(c *gin.Context) {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	var req searchRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return 
	}

	px := pexels.NewClient(config.PEXELAPIKEY)
	ctx := context.Background()
	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{Query: req.Name, Color: req.Color})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
	}
	c.JSON(http.StatusOK, res)

}

