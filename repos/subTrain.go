package repos

import (
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/khades/servbot/models"
)

func GetChannelsWithSubtrainNotification() (*[]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := Db.C(channelInfoCollection).Find(
		bson.M{
			"subtrain.enabled": true,
			"subtrain.notificationshown" : false,
			"subtrain.currentstreak": bson.M{
				"$ne":0},
			"subtrain.notificationtime":bson.M{
				"$lt": time.Now()}}).All(&result)
	return &result, error
}

func GetChannelsWithExpiredSubtrain() (*[]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := Db.C(channelInfoCollection).Find(
		bson.M{
			"subtrain.enabled": true,
			"subtrain.currentstreak": bson.M{
				"$ne":0},
			"subtrain.expirationtime":bson.M{
				"$lt": time.Now()}}).All(&result)
	return &result, error
}

func PutChannelSubtrain(channelID *string, subTrain *models.SubTrain) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.SubTrain = *subTrain
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, SubTrain: *subTrain})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subtrain": *subTrain}})
}

func PutChannelSubtrainWeb(channelID *string, subTrain *models.SubTrain) {
	channelInfo, _ := GetChannelInfo(channelID)
	localSubtrain := channelInfo.SubTrain
	if (subTrain.Enabled == true && localSubtrain.Enabled == true && localSubtrain.ExpirationLimit == subTrain.ExpirationLimit && localSubtrain.NotificationLimit == subTrain.NotificationLimit) {
		subTrain.ExpirationTime = localSubtrain.ExpirationTime
		subTrain.NotificationTime = localSubtrain.NotificationTime
		subTrain.CurrentStreak = localSubtrain.CurrentStreak
		subTrain.Users = localSubtrain.Users
		subTrain.NotificationShown = localSubtrain.NotificationShown
	}
	if channelInfo != nil {
		channelInfo.SubTrain = *subTrain
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, SubTrain: *subTrain})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subtrain": *subTrain}})
}

func SetSubtrainNotificationShown(channelInfo *models.ChannelInfo) {
	subTrain := channelInfo.SubTrain
	subTrain.NotificationShown = true
	PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}	

func IncrementSubtrainCounterByChannelID(channelID *string, user *string) {
	log.Println("Got to incrementSubTrainByID")
	channelInfo, error := GetChannelInfo(channelID)
	if error == nil {
		IncrementSubtrainCounter(channelInfo, user)
		return
	} 
	log.Println(error)
}

func IncrementSubtrainCounter(channelInfo *models.ChannelInfo, user *string) {
	subTrain := channelInfo.SubTrain
	if subTrain.Enabled == false {
		log.Println("Got to incrementSubTrainByID FALSE")

		return
	}
	subTrain.ExpirationTime = time.Now().Add(time.Second * time.Duration(subTrain.ExpirationLimit))
	subTrain.NotificationTime = time.Now().Add(time.Second * time.Duration(subTrain.NotificationLimit))
	subTrain.CurrentStreak = subTrain.CurrentStreak +1
	subTrain.Users = append(subTrain.Users, *user)
	PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}

func ResetSubtrainCounter(channelInfo *models.ChannelInfo) {
	subTrain := channelInfo.SubTrain
	subTrain.CurrentStreak = 0
	subTrain.NotificationShown = false
	subTrain.Users = []string{}
	PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}