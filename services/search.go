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

type Photo struct {
	ID              int    `json:"id"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	URL             string `json:"url"`
	Photographer    string `json:"photographer"`
	PhotographerURL string `json:"photographer_url"`
	PhotographerID  int    `json:"photographer_id"`
	AvgColor        string `json:"avg_color"`
	Liked           bool   `json:"liked"`
	Src             Source `json:"src"`
	Alt             string `json:"alt"`
}

type Source struct {
	Original  string `json:"original"`
	Large2X   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type SearchPhotoResponse struct {
	TotalResults int      `json:"total_results"`
	Page         int      `json:"page"`
	PerPage      int      `json:"per_page"`
	Photos       []*Photo `json:"photos"`
	NextPage     string   `json:"next_page"`
	PrevPage     string   `json:"prev_page"`
}

func getRandomAfroSymns() string {
	symns := []string{"africa", "black", "afro"}
	randomSymn := randomize(symns)
	return randomSymn
}

// @Summary Search Pexel photos
// @Description Requires a color and query and returns a list of photos
// @Tags search
// @Accept json
// @Produce json
// @Param searchRequest body searchRequest true "Search request"
// @Success 200 {object} SearchPhotoResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/search [post]
func SearchPexel(c *gin.Context) {
	var req searchRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	

	px := pexels.NewClient(os.Getenv("PEXELS_API_KEY"))
	ctx := context.Background()
	queryPrefix := getRandomAfroSymns()
	res, err := px.PhotoService.Search(ctx, &pexels.PhotoParams{
		Query:       queryPrefix + req.Query,
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
