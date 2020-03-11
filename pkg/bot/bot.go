package bot

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"

	"github.com/f110/discord-bot/pkg/amesh"
	"github.com/f110/discord-bot/pkg/command"
	"github.com/f110/discord-bot/pkg/config"
	"github.com/f110/discord-bot/pkg/handler"
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
	handler   *handler.Handler
	conf      *config.Config

	amesh   *amesh.Generator
	storage *storage.Storage
}

func NewBot(name, token string, storageHost, bucket, bucketHost string, conf *config.Config) (*Bot, error) {
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
		handler:     handler.New(name),
		conf:        conf,
	}
	dg.AddHandler(b.handler.ReceiveEvent)

	for _, p := range conf.EnablePlugins {
		log.Printf("Enabling %s", p)
		b.LoadPlugin(p)
	}

	return b, nil
}

func (b *Bot) LoadPlugin(name string) {
	p, ok := command.Manager.Fetch(name)
	if !ok {
		log.Printf("Plugin %s is not found", name)
		return
	}

	if err := p.Enable(b.conf); err != nil {
		log.Printf("Failed enabling %s: %v", name, err)
		return
	}
	p.Subscribe(b.handler)
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
