package main

import (
	"bytes"
	"errors"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config holds configuration data for the handler
type Config struct {
	sensu.PluginConfig
	APIToken        string
	ChatID          int64
	MessageTemplate string
}

var config = Config{
	PluginConfig: sensu.PluginConfig{
		Name:     "sensu-telegram-handler",
		Short:    "Sensu Go handler for sending telegram notifications",
		Keyspace: "sensu.io/plugins/sensu-telegram-handler/config",
	},
}

// handler options
const defaultMessageTemplate = "*{{.Entity.Name}}/{{.Check.Name}}*: {{.Check.State}}\n`{{.Check.Output}}`"

var options = []*sensu.PluginConfigOption{
	{
		Path:      "api-token",
		Argument:  "api-token",
		Shorthand: "a",
		Default:   "",
		Usage:     "The API token to use when connecting to the Telegram service",
		Value:     &config.APIToken,
	},
	{
		Path:      "chatid",
		Argument:  "chatid",
		Shorthand: "c",
		Default:   int64(0),
		Usage:     "The Chat ID to use when connecting to the Telegram service",
		Value:     &config.ChatID,
	},
	{
		Path:      "template",
		Argument:  "template",
		Shorthand: "t",
		Usage:     "The default message template",
		Default:   defaultMessageTemplate,
		Value:     &config.MessageTemplate,
	},
}

func main() {
	telegramHandler := sensu.NewGoHandler(&config.PluginConfig, options, checkArgs, executeHandler)
	telegramHandler.Execute()
}

func checkArgs(_ *types.Event) error {
	if config.APIToken == "" {
		return errors.New("missing api token")
	}
	if config.ChatID == 0 {
		return errors.New("missing chat id")
	}
	return nil
}

func executeHandler(event *types.Event) error {
	// initialize bot
	bot, err := tgbotapi.NewBotAPI(config.APIToken)
	if err != nil {
		return err
	}

	// render template
	messageTemplate, err := template.New("message").Parse(config.MessageTemplate)
	if err != nil {
		return err
	}
	var text bytes.Buffer
	err = messageTemplate.Execute(&text, event)
	if err != nil {
		return err
	}

	// send message
	message := tgbotapi.NewMessage(int64(config.ChatID), text.String())
	message.ParseMode = tgbotapi.ModeMarkdown
	_, err = bot.Send(message)
	if err != nil {
		return err
	}
	return nil
}
