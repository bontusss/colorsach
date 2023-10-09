package api

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kosa3/pexels-go"
)

type searchRequest struct {
	Color string `json:"color" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// Gin handler for getting image from pexels api
func SearchPexel(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err loading .env file")
	}
	key := os.Getenv("PEXELS_API_KEY")
	var req searchRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return 
	}

	px := pexels.NewClient(key)
	ctx := context.Background()
	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{Query: req.Name, Color: req.Color})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
	}
	c.JSON(http.StatusOK, res)

}

