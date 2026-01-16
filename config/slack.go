package config

type SlackConfig struct {
	BotToken string `env:"BOT_TOKEN, required"`
	AppToken string `env:"APP_TOKEN, required"`
}
