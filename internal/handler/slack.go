package handler

import (
	"generate-leetcode/internal/service/link"
	"log/slog"

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
				continue
			}
			socket.Ack(*evt.Request)

			h.handleEventsAPI(eventsAPIEvent)
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
	link, err := h.linkService.GetLink(event.Text)
	if err != nil {
		slog.Error("failed to get link", "error", err)
		return
	}

	_, _, err = h.client.PostMessage(
		event.Channel,
		slack.MsgOptionText(link, false),
		slack.MsgOptionTS(event.TimeStamp),
	)
	if err != nil {
		slog.Error("failed to post message", "error", err)
	}
}
