package services

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kosa3/pexels-go"
)

type searchRequest struct {
	Color       string `json:"color" binding:"required"`
	Query       string `json:"query" binding:"required"` // The search query. Ocean, Tigers, Pears, etc.
	Page        int    `json:"page"`                     // The page number you are requesting. Default: 1
	Size        string `json:"size"`                     // Minimum photo size. The current supported sizes are: large(24MP), medium(12MP) or small(4MP).
	PerPage     int    `json:"per-page"`                 // The number of results you are requesting per page. Default: 15 Max: 80
	Orientation string `json:"orientation"`              // Desired photo orientation. The current supported orientations are: landscape, portrait or square.
}

func SearchPexel(c *gin.Context) {
	var req searchRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	px := pexels.NewClient(os.Getenv("PEXELS_API_KEY"))
	ctx := context.Background()
	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{
		Query:       req.Query,
		Color:       req.Color,
		Page:        req.Page,
		Orientation: req.Orientation,
		Size:        req.Size,
		PerPage:     req.PerPage,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, res)
	//return res

}

func randomize[T any] (list []T) T {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rand.Intn(len(list))
	return list[index]
}

func GetRandomImage() (*pexels.Photo, error) {
	queryList := []string{"car", "house", "animal", "leggo", "wildlife", "forest", "ship", "africa"}
	colors := []string{"red", "blue", "green", "orange", "purple", "black"}
	randomQuery := randomize(queryList)
	randomColor := randomize(colors)

	px := pexels.NewClient(os.Getenv("PEXELS_API_KEY"))
	res, err := px.PhotoService.Search(context.Background(), &pexels.PhotoParams{
		Query: randomQuery,
		Color: randomColor,
	})
	if err != nil {
		return nil, err
	}
	photo := randomize(res.Photos)
	return photo, nil
}
