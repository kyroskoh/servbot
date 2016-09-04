package commandHandlers

import (
	"strings"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// New works with new command that adds template
func New(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	if chatMessage.IsMod {
		commandName := ""
		template := ""
		separator := strings.Index(chatCommand.Body, "=")
		if separator != -1 {
			commandName = chatCommand.Body[:separator]
			template = strings.TrimSpace(chatCommand.Body[separator+1:])
		} else {
			commandName = chatCommand.Body
		}
		if strings.HasPrefix(template, "!") || strings.HasPrefix(template, ".") || strings.HasPrefix(template, "/") {
			ircClient.SendPublic(models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Запрещено зацикливать команды",
				User:    chatMessage.User})
		} else {
			repos.PutChannelTemplate(chatMessage.User, chatMessage.Channel, commandName, commandName, template)
			Template.updateTemplate(chatMessage.Channel, commandName, template)

			Template.updateAliases(chatMessage.Channel, commandName, template)
			ircClient.SendPublic(models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Ну в принципе готово VoHiYo",
				User:    chatMessage.User})
		}
	} else {
		ircClient.SendPublic(models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание алиaса: Вы не модер SMOrc",
			User:    chatMessage.User})
	}
}
