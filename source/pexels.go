package source

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bontusss/colosach/utils"
)

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
	// Src             Source `json:"src"`
	Alt string `json:"alt"`
}

func GetPexels(image, color string) (photo Photo, err error) {
	fmt.Println("Calling pexels API...")
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pexels.com/v1/search?query=%s&color=%s", image, color), nil)
	fmt.Println(req) //remove 
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("Authorization", config.PEXELSAPIKEY)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(resp)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var resObj Photo
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		log.Fatal(err.Error())
	}
	return resObj, err
}