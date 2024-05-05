package utils

import (
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
