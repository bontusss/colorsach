package services

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

func SearchPexel(c *gin.Context) {
	//todo load env via config package
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("err loading app.env file")
	}
	key := os.Getenv("PEXELS_API_KEY")
	var req searchRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})

	}

	px := pexels.NewClient(key)
	ctx := context.Background()
	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{Query: req.Name, Color: req.Color})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, res)
	//return res

}
