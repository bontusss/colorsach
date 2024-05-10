package utils

import (
	"math/rand"
	"time"

	"github.com/bontusss/colosach/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}
	doc = &bson.D{}
	err = bson.Unmarshal(data, doc)
	return
}

func GetCurrentUser(c *gin.Context) *models.DBResponse {
	userData := c.MustGet("currentUser").(*models.DBResponse)
	return userData
}

func GeneratePassword(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABZDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
