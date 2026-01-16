package config

type SlackConfig struct {
	BotToken string `env:"SLACK_BOT_TOKEN, required"`
	AppToken string `env:"SLACK_APP_TOKEN, required"`
}
