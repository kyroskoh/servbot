package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// GetChannelActiveTemplates is setting
func GetChannelActiveTemplates(channel string) ([]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := Db.C("templates").Find(bson.M{"channel": channel, "template": bson.M{"$ne": ""}}).All(&result)
	return result, error
}