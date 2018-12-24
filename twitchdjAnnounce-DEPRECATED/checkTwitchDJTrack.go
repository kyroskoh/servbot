package services

import (
	"encoding/json"
	"html"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

type tdjTrack struct {
	Title string
}

// CheckTwitchDJTrack checks last playing track
func CheckTwitchDJTrack() {
	channels, error := repos.GetTwitchDJEnabledChannels()
	if error != nil {
		return
	}
	for _, channel := range channels {
		checkOneTwitchDJTrack(&channel)
	}
}

func checkOneTwitchDJTrack(channel *models.ChannelInfo) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "twitchdj",
		"action":  "checkOneTwitchDJTrack"})
	status := models.TwitchDJ{ID: channel.TwitchDJ.ID}
	defer repos.PushTwitchDJ(&channel.ChannelID, &status)
	logger.Debugf("Checking %s twitchDj track", channel.Channel)
	resp, error := httpclient.Get("https://twitch-dj.ru/includes/back.php?func=get_track&channel=" + channel.TwitchDJ.ID)

	if error != nil {
		logger.Debug(error)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	track := tdjTrack{}
	marshallError := json.NewDecoder(resp.Body).Decode(&track)
	if marshallError != nil {
		logger.Debug(marshallError)
		return
	}
	if track.Title != "" {
		status.Playing = true
		status.Track = html.UnescapeString(track.Title)
	}
	if status.Playing == false {
		return
	}
	if channel.TwitchDJ.NotifyOnChange == true {
		if status.Playing == true && channel.TwitchDJ.Track != status.Track {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channel.Channel,
				Body:    "[TwitchDJ] Now Playing: " + status.Track})
		}
	}
}