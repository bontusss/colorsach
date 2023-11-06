package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kosa3/pexels-go"
	"log"
	"net/http"
	"os"
)

type searchRequest struct {
	Color string `json:"color" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

//func Search(c *gin.Context) {
//	var colorsach []interface{}
//	pexel := SearchPexel(c)
//	splash := GetUnsplash(c)
//
//	colorsach = append(colorsach, pexel)
//	colorsach = append(colorsach, splash)
//	c.JSON(http.StatusOK, colorsach)
//}

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

	}

	px := pexels.NewClient(key)
	ctx := context.Background()
	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{Query: req.Name, Color: req.Color})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
	}
	c.JSON(http.StatusOK, res)
	//return res

}
