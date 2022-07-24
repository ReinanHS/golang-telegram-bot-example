package handler

import (
	"errors"
	"fmt"
	telegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/reinanhs/golang-telegram-bot-example/pkg/handler/command"
	"log"
	"strings"
)

// ActionHandler method responsible for deciding the action that the received message will perform
func ActionHandler(update *telegramApi.Update, bot *telegramApi.BotAPI) error {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	if update.Message.IsCommand() {
		return ActionHandlerCommands(update, bot)
	}

	return ActionHandlerDialog(update, bot)
}

// ActionHandlerCommands method responsible for rendering the commands
func ActionHandlerCommands(update *telegramApi.Update, bot *telegramApi.BotAPI) error {
	commandName := strings.ToLower(update.Message.Command())

	for _, c := range command.GetEnabledCommands() {
		if commandName == strings.ToLower(c.GetCommandName()) {
			return c.CommandAction(update, bot)
		}
	}

	return replyMessageByUpdate("Sorry this command doesn't exist, see the list of all available commands using /commands", update, bot)
}

// ActionHandlerDialog method responsible for rendering the dialogs
func ActionHandlerDialog(update *telegramApi.Update, bot *telegramApi.BotAPI) error {
	return replyMessageByUpdate(update.Message.Text, update, bot)
}

// replyMessageByUpdate method responsible for sending a text message back to the user
func replyMessageByUpdate(message string, update *telegramApi.Update, bot *telegramApi.BotAPI) error {
	msg := telegramApi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Could not complete message send: %s", err.Error())
		return errors.New(fmt.Sprintf("Could not complete message send to chat id %d", update.Message.Chat.ID))
	}

	return nil
}
