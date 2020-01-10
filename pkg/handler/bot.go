package handler

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/f110/discord-bot/pkg/amesh"
	"github.com/f110/discord-bot/pkg/storage"
)

type Bot struct {
	Name        string
	storageHost string
	bucket      string
	bucketHost  string

	client    *discordgo.Session
	doneCh    chan struct{}
	closeOnce sync.Once

	amesh   *amesh.Generator
	storage *storage.Storage
}

func New(name, token string, storageHost, bucket, bucketHost string) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	dg.UserAgent = "discord-bot"
	b := &Bot{
		Name:        name,
		storageHost: storageHost,
		bucket:      bucket,
		bucketHost:  bucketHost,
		client:      dg,
		doneCh:      make(chan struct{}),
	}
	dg.AddHandler(b.handleMessageCreate)

	if err := b.init(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bot) init() error {
	g, err := amesh.NewGenerator()
	if err != nil {
		return err
	}
	b.amesh = g

	b.storage = storage.New(b.storageHost, b.bucket)

	return nil
}

func (b *Bot) Run() error {
	if err := b.client.Open(); err != nil {
		return err
	}

	return nil
}

func (b *Bot) Wait() {
	<-b.doneCh
	return
}

func (b *Bot) Shutdown() error {
	if err := b.client.Close(); err != nil {
		return err
	}

	b.closeOnce.Do(func() {
		close(b.doneCh)
	})

	return nil
}

func (b *Bot) handleMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	// Ignore myself
	if e.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(e.Content, b.Name) {
		return
	}

	c := strings.Split(e.Content, " ")[1:]
	if len(c) == 0 {
		return
	}

	switch c[0] {
	case "amesh":
		b.ameshHandler(s, e)
	}
}

func (b *Bot) ameshHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	date := b.amesh.LatestTime()
	filename := fmt.Sprintf("amesh/%s.png", date)
	u := fmt.Sprintf("https://%s/%s", b.bucketHost, filename)

	if b.storage.Exist(filename) {
		// Already generated and stored
		_, err := s.ChannelMessageSend(e.ChannelID, u)
		if err != nil {
			log.Print(err)
		}
		return
	}

	img, err := b.amesh.Generate(date)
	if err != nil {
		_, err := s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("画像の生成でエラーだよ！: %v", err))
		if err != nil {
			log.Print(err)
		}
		return
	}

	_, err = b.storage.Store(filename, img)
	if err != nil {
		_, err := s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("画像の保存でエラーだよ！: %v", err))
		if err != nil {
			log.Print(err)
		}
		return
	}

	_, err = s.ChannelMessageSend(e.ChannelID, u)
	if err != nil {
		log.Print(err)
	}
}
