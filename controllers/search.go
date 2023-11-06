package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hbagdi/go-unsplash/unsplash"
	"github.com/joho/godotenv"
	"github.com/kosa3/pexels-go"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

type searchRequest struct {
	Color string `json:"color" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

func Search(c *gin.Context) {
	var colorsach []interface{}
	pexel := SearchPexel(c)
	colorsach = append(colorsach, pexel)
	c.JSON(http.StatusOK, colorsach)
}

func SearchPexel(c *gin.Context) *pexels.SearchPhotoResponse {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err loading .env file")
	}
	key := os.Getenv("PEXELS_API_KEY")
	var req searchRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return nil
	}

	px := pexels.NewClient(key)
	ctx := context.Background()
	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{Query: req.Name, Color: req.Color})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
	}
	return res

}

//func SearchPexel(c *gin.Context) {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("err loading .env file")
//	}
//	key := os.Getenv("PEXELS_API_KEY")
//	var req searchRequest
//	err = c.ShouldBindJSON(&req)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, errResponse(err))
//		return
//	}
//
//	px := pexels.NewClient(key)
//	ctx := context.Background()
//	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{Query: req.Name, Color: req.Color})
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, errResponse(err))
//	}
//	c.JSON(http.StatusOK, res)
//
//}

func UnsplashSSearch(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err loading .env file")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("UNSPLASH_KEY")})
	client := oauth2.NewClient(context.Background(), ts)
	sdk := unsplash.New(client)

	var req searchRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	searchOpt := unsplash.SearchOpt{Query: req.Name}
	photos, _, err := sdk.Search.Photos(&searchOpt)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
	}
	found := false

	for _, photo := range *photos.Results {
		if *photo.Color == req.Color {
			c.JSON(http.StatusOK, photo)
			found = true
		}
	}
	if !found {
		c.JSON(http.StatusOK, fmt.Sprintf("%s colored %s not available", req.Color, req.Name))
	}
	//c.JSON(http.StatusOK, photos)

}
