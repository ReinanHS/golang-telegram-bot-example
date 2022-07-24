package command

import (
	"fmt"
	telegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// Command basic structure to create a command
type Command struct {
	Name        string
	Description string
}

// StartCommand this is the initial command
type StartCommand struct {
	Command
}

// HelpCommand this is the help command
type HelpCommand struct {
	Command
}

// ActionCommand this is an interface for having the special methods of the commands
type ActionCommand interface {
	CommandAction(update *telegramApi.Update, bot *telegramApi.BotAPI) error
	GetCommandName() string
}

// GetEnabledCommands method responsible for returning all enabled commands
func GetEnabledCommands() []ActionCommand {
	return []ActionCommand{
		StartCommand{
			Command: Command{
				Name:        "Start",
				Description: "Initial command",
			},
		},
		HelpCommand{
			Command: Command{
				Name:        "Commands",
				Description: "Lists all available commands",
			},
		},
	}
}

// GetCommandName method responsible for returning the command name
func (c Command) GetCommandName() string {
	return strings.ToLower(c.Name)
}

// CommandAction method responsible for showing a welcome message
func (c StartCommand) CommandAction(update *telegramApi.Update, bot *telegramApi.BotAPI) error {
	msgText := fmt.Sprintf("Hello welcome, %s", update.Message.From.UserName)

	msg := telegramApi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(msg)

	return err
}

// CommandAction method responsible for showing the help message
func (c HelpCommand) CommandAction(update *telegramApi.Update, bot *telegramApi.BotAPI) error {
	var msgText strings.Builder

	msgText.WriteString("I can help you create and manage <b>websites</b> and <b>landing page</b>.")
	msgText.WriteString("If you're new here, check out this list of commands you can use to interact with the bot.")
	msgText.WriteString("\n\n")
	msgText.WriteString("You can control me by sending these commands: \n")
	msgText.WriteString("\n/newpage - create a new page")
	msgText.WriteString("\n/mypages - edit your pages [beta]")

	text := msgText.String()

	msg := telegramApi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = telegramApi.ModeHTML
	_, err := bot.Send(msg)

	return err
}
