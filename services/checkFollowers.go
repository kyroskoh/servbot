package services

import (
	"encoding/json"
	"strings"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type links struct {
	Next string `json:"next"`
}
type follows struct {
	User user `json:"user"`
}
type user struct {
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}
type followerResponse struct {
	Cursor  string    `json:"_cursor"`
	Links   links     `json:"_links"`
	Follows []follows `json:"follows"`
}

// https://api.twitch.tv/kraken/channels/areyouready_88/follows?client_id=1icafjrio64l3n0uwg5qlkdseu0j8d6&direction=ASC&cursor=1494757816579821000
func CheckChannelsFollowers() {
	userIDs, error := repos.GetUsersID(&repos.Config.Channels)
	if error != nil {
		return
	}
	for key, value := range *userIDs {
		checkOneChannelFollowers(&key, &value)
	}

}
func checkOneChannelFollowers(channel *string, channelID *string) {
	cursor := ""
	cursorObject, error := repos.GetFollowerCursor(channelID)
	if error != nil && error.Error() != "not found" {
		return
	}

	if error != nil && error.Error() == "not found" || cursorObject.Cursor == "" {
		cursor = getInitialCursor(channelID)
		if cursor != "" {
			repos.SetFollowerCursor(channelID, &cursor)
		}

	} else {
		cursor = cursorObject.Cursor
	}
	if cursor == "" {
		return
	}

	followers, followersError := getFollowers(channelID, &cursor)
	if followersError != nil || followers.Cursor == "" || len(followers.Follows) == 0 {
		return
	}

	repos.SetFollowerCursor(channelID, &followers.Cursor)
	followersToGreet := []string{}
	for _, follow := range followers.Follows {
		alreadyGreeted, _ := repos.CheckIfFollowerGreeted(channelID, &follow.User.Name)
		if alreadyGreeted == false {
			followersToGreet = append(followersToGreet, follow.User.Name)

			repos.AddFollowerToList(channelID, &follow.User.Name)
			// alertInfo, alertError := repos.GetSubAlert(channelID)
			// if alertError == nil && alertInfo.Enabled == true && alertInfo.FollowerMessage != "" {
			// 	bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
			// 		Channel: *channel,
			// 		Body:    alertInfo.FollowerMessage,
			// 		User:    follow.User.Name})
			// }

		}

	}
	if len(followersToGreet) > 0 {
		alertInfo, alertError := repos.GetSubAlert(channelID)
		if alertError == nil && alertInfo.Enabled == true && alertInfo.FollowerMessage != "" {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: *channel,
				Body:    "@" + strings.Join(followersToGreet, " @") + " " + alertInfo.FollowerMessage})
		}
	}

}
func getFollowers(channelID *string, cursor *string) (*followerResponse, error) {
	url := "https://api.twitch.tv/kraken/channels/" + *channelID + "/follows?direction=ASC&cursor=" + *cursor

	resp, respError := httpclient.TwitchV5(repos.Config.ClientID, "GET", url, nil)
	if respError != nil {
		return nil, respError
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var responseBody = new(followerResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)

	if marshallError != nil {
		return nil, marshallError
	}
	return responseBody, nil

}
func getInitialCursor(channelID *string) string {
	url := "https://api.twitch.tv/kraken/channels/" + *channelID + "/follows?direction=DESC&limit=1"
	resp, respError := httpclient.TwitchV5(repos.Config.ClientID, "GET", url, nil)
	if respError != nil {
		return ""
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var responseBody = new(followerResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)
	if marshallError != nil {

		return ""
	}
	return responseBody.Cursor
}
