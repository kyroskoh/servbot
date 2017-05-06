package commandHandlers

import (
	"strings"

	"html"

	"github.com/hoisie/mustache"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// Custom handler checks if input command has template and then fills it in with mustache templating and sends to a specified/user
func Custom(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	template, err := repos.GetChannelTemplate(&chatMessage.ChannelID, &chatCommand.Command)
	if err == nil {
		values, _ := repos.GetChannelInfo(&chatMessage.ChannelID)
		message := mustache.Render(template.Template, values)
		user := chatMessage.User
		redirectTo := chatMessage.User
		if chatCommand.Body != "" {
			if strings.HasPrefix(chatCommand.Body, "@") {
				redirectTo = chatCommand.Body[1:]
			} else {
				redirectTo = chatCommand.Body

			}
		}
		ircClient.SendDebounced(models.OutgoingDebouncedMessage{
			Message: models.OutgoingMessage{
				Channel: chatMessage.Channel,
				User:    user,
				Body:    html.UnescapeString(message)},
			Command:    chatCommand.Command,
			RedirectTo: redirectTo})
	}
}
