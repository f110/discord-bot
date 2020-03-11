package handler

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CreateEventHandler func(s *discordgo.Session, e *discordgo.MessageCreate)

type Handler struct {
	Name              string
	createSubscribers []CreateEventHandler
}

func New(name string) *Handler {
	return &Handler{
		Name:              name,
		createSubscribers: make([]CreateEventHandler, 0),
	}
}

func (h *Handler) SubscribeCreate(handler CreateEventHandler) {
	h.createSubscribers = append(h.createSubscribers, handler)
}

func (h *Handler) ReceiveEvent(s *discordgo.Session, e interface{}) {
	switch v := e.(type) {
	case *discordgo.MessageCreate:
		if v.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(v.Content, h.Name) {
			return
		}

		for _, subscriber := range h.createSubscribers {
			subscriber(s, v)
		}
	}
}
