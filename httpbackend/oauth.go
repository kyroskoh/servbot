package httpbackend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type tokenResponse struct {
	Token string `json:"access_token"`
}

func oauth(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Incoming Twitch code is missing", http.StatusUnprocessableEntity)
		return
	}
	fmt.Fprintf(w, "K, we parsed %s", code)

	resp, err := http.PostForm("https://api.twitch.tv/kraken/oauth2/token",
		url.Values{
			"client_id":     {repos.Config.ClientID},
			"client_secret": {repos.Config.ClientSecret},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {repos.Config.AppOauthURL},
			"code":          {code}})
	//"state":         {}
	if err != nil {
		log.Println(err)
		http.Error(w, "Twitch Error", http.StatusUnprocessableEntity)
		return
	}
	var responseBody = new(tokenResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)
	if marshallError != nil {
		log.Println(marshallError)
		http.Error(w, "Twitch Error", http.StatusUnprocessableEntity)
		return
	}

	defer resp.Body.Close()
}
