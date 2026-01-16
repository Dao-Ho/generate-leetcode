package handler

import (
	"generate-leetcode/internal/service/link"
	"log/slog"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackHandler struct {
	client      *slack.Client
	linkService *link.LinkService
}

func NewSlackHandler(client *slack.Client, linkService *link.LinkService) *SlackHandler {
	return &SlackHandler{
		client:      client,
		linkService: linkService,
	}
}

func (h *SlackHandler) HandleEvents(socket *socketmode.Client) {
	for evt := range socket.Events {
		switch evt.Type {
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				slog.Warn("failed to cast event data")
				continue
			}
			socket.Ack(*evt.Request)
			h.handleEventsAPI(eventsAPIEvent)

		case socketmode.EventTypeConnectionError:
			slog.Error("connection error", "data", evt.Data)

		case socketmode.EventTypeConnecting:
			slog.Info("connecting to Slack...")

		case socketmode.EventTypeConnected:
			slog.Info("connected to Slack")
		}
	}
}

func (h *SlackHandler) handleEventsAPI(event slackevents.EventsAPIEvent) {
	switch event.Type {
	case slackevents.CallbackEvent:
		switch ev := event.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			h.handleAppMention(ev)
		}
	}
}

func (h *SlackHandler) handleAppMention(event *slackevents.AppMentionEvent) {
	cmd := ParseCommand(event.Text)

	slog.Debug("received app mention",
		"user", event.User,
		"channel", event.Channel,
		"text", event.Text,
		"command", cmd.String(),
	)

	switch cmd {
	case CommandRandom:
		link, err := h.linkService.GetLink(cmd.String())
		if err != nil {
			slog.Error("failed to get link", "error", err)
			return
		}
		h.reply(event, link)

	case CommandHelp:
		helpText := "Supported commands: `" + strings.Join(ListCommands(), "`, `") + "`"
		h.reply(event, helpText)

	case CommandUnknown:
		h.reply(event, "Unknown command. Try `help` for available commands.")
	}
}

func (h *SlackHandler) reply(event *slackevents.AppMentionEvent, text string) {
	_, _, err := h.client.PostMessage(
		event.Channel,
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(event.TimeStamp),
	)
	if err != nil {
		slog.Error("failed to post message", "error", err)
	}
}
