package main

import (
	"bytes"
	"errors"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
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

var options = []sensu.ConfigOption{
	&sensu.PluginConfigOption[string]{
		Path:      "api-token",
		Argument:  "api-token",
		Shorthand: "a",
		Default:   "",
		Usage:     "The API token to use when connecting to the Telegram service",
		Value:     &config.APIToken,
	},
	&sensu.PluginConfigOption[int64]{
		Path:      "chatid",
		Argument:  "chatid",
		Shorthand: "c",
		Default:   int64(0),
		Usage:     "The Chat ID to use when connecting to the Telegram service",
		Value:     &config.ChatID,
	},
	&sensu.PluginConfigOption[string]{
		Path:      "template",
		Argument:  "template",
		Shorthand: "t",
		Usage:     "The default message template",
		Default:   defaultMessageTemplate,
		Value:     &config.MessageTemplate,
	},
}

func main() {
	telegramHandler := sensu.NewGoHandler(&config.PluginConfig, options, validateInput, executeHandler)
	telegramHandler.Execute()
}

func validateInput(_ *corev2.Event) error {
	if config.APIToken == "" {
		return errors.New("missing api token")
	}
	if config.ChatID == 0 {
		return errors.New("missing chat id")
	}
	return nil
}

func executeHandler(event *corev2.Event) error {
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
	message := tgbotapi.NewMessage(config.ChatID, text.String())
	message.ParseMode = tgbotapi.ModeMarkdown
	_, err = bot.Send(message)
	if err != nil {
		return err
	}
	return nil
}
