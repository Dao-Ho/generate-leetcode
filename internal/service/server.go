package service

import (
	"context"

	"generate-leetcode/config"
	"generate-leetcode/internal/handler"
	"generate-leetcode/internal/service/link"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type App struct {
	SlackClient  *slack.Client
	SocketClient *socketmode.Client
}

func InitApp(ctx context.Context, cfg *config.Config) (*App, error) {
	slackClient := slack.New(
		cfg.Slack.BotToken,
		slack.OptionAppLevelToken(cfg.Slack.AppToken),
	)

	socketClient := socketmode.New(slackClient)

	linkService := link.NewLinkService()
	slackHandler := handler.NewSlackHandler(slackClient, linkService)

	go slackHandler.HandleEvents(socketClient)

	return &App{
		SlackClient:  slackClient,
		SocketClient: socketClient,
	}, nil
}
