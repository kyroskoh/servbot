package bot

import (
	"log"
	"net"

	"gopkg.in/irc.v2"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/repos"
)

type chatClient struct {
	Client *irc.Client
	Ready  bool
}

// Start function dials up connection for chathandler
func Start() {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Fatalln(err)
	}

	config := irc.ClientConfig{
		Nick:    repos.Config.BotUserName,
		Pass:    repos.Config.OauthKey,
		User:    repos.Config.BotUserName,
		Name:    repos.Config.BotUserName,
		Handler: chatHandler}
	chatClient := irc.NewClient(conn, config)
	log.Println("Bot is starting...")

	clientError := chatClient.Run()
	log.Println(clientError)
	log.Println("Bot died...")
	IrcClientInstance = &ircClient.IrcClient{Ready: false, MessageQueue: []string{}}
	conn.Close()
}
