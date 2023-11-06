package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

type Unsplash struct {
	Total      int       `json:"total"`
	TotalPages int       `json:"total_pages"`
	Results    []Results `json:"results"`
}
type ProfileImage struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}
type Liinks struct {
	Self   string `json:"self"`
	HTML   string `json:"html"`
	Photos string `json:"photos"`
	Likes  string `json:"likes"`
}
type User struct {
	ID                string       `json:"id"`
	Username          string       `json:"username"`
	Name              string       `json:"name"`
	FirstName         string       `json:"first_name"`
	LastName          string       `json:"last_name"`
	InstagramUsername string       `json:"instagram_username"`
	TwitterUsername   string       `json:"twitter_username"`
	PortfolioURL      string       `json:"portfolio_url"`
	ProfileImage      ProfileImage `json:"profile_image"`
	Links             Liinks       `json:"links"`
}
type Urls struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}
type Links struct {
	Self     string `json:"self"`
	HTML     string `json:"html"`
	Download string `json:"download"`
}
type Results struct {
	ID                     string `json:"id"`
	CreatedAt              string `json:"created_at"`
	Width                  int    `json:"width"`
	Height                 int    `json:"height"`
	Color                  string `json:"color"`
	BlurHash               string `json:"blur_hash"`
	Likes                  int    `json:"likes"`
	LikedByUser            bool   `json:"liked_by_user"`
	Description            string `json:"description"`
	User                   User   `json:"user"`
	CurrentUserCollections []any  `json:"current_user_collections"`
	Urls                   Urls   `json:"urls"`
	Links                  Links  `json:"links"`
}

func GetUnsplash(ctx *gin.Context) {
	fmt.Println("1. Performing Http GET")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err loading .env file")
	}
	key := os.Getenv("UNSPLASH_KEY")
	var r searchRequest
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}
	url := fmt.Sprintf("https://api.unsplash.com/search/photos?query=%s&color=%s", r.Name, r.Color)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Close = true
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		fmt.Errorf("errror: %s")
	}
	req.Header.Set("Authorization", key)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: &http.Transport{}}
	res, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	var sr Unsplash
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}
	ctx.JSON(http.StatusOK, sr)
	//return &sr
}
