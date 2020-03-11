package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/f110/discord-bot/pkg/amesh"
	"github.com/f110/discord-bot/pkg/config"
	"github.com/f110/discord-bot/pkg/handler"
	"github.com/f110/discord-bot/pkg/storage"
)

func init() {
	Manager.Register("amesh", NewAmesh())
}

type Amesh struct {
	bucketHost string

	enable  bool
	amesh   *amesh.Generator
	storage *storage.Storage
}

func NewAmesh() *Amesh {
	return &Amesh{}
}

func (a *Amesh) Enable(conf *config.Config) error {
	g, err := amesh.NewGenerator()
	if err != nil {
		return err
	}
	a.amesh = g

	a.storage = storage.New(conf.StorageHost, conf.Bucket)
	a.bucketHost = conf.BucketHost

	a.enable = true

	return nil
}

func (a *Amesh) Subscribe(handler *handler.Handler) {
	handler.SubscribeCreate(a.Execute)
}

func (a *Amesh) Execute(s *discordgo.Session, e *discordgo.MessageCreate) {
	if !a.enable {
		log.Print("Call amesh plugin but amesh plugin is not enabled.")
		return
	}
	c := strings.Split(e.Content, " ")
	if len(c) < 2 || c[1] != "amesh" {
		log.Printf("Unknown command: %s", e.Content)
		return
	}

	date := a.amesh.LatestTime()
	filename := fmt.Sprintf("amesh/%s.png", date)
	u := fmt.Sprintf("https://%s/%s", a.bucketHost, filename)

	if a.storage.Exist(filename) {
		// Already generated and stored
		_, err := s.ChannelMessageSend(e.ChannelID, u)
		if err != nil {
			log.Print(err)
		}
		return
	}

	img, err := a.amesh.Generate(date)
	if err != nil {
		_, err := s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("画像の生成でエラーだよ！: %v", err))
		if err != nil {
			log.Print(err)
		}
		return
	}

	_, err = a.storage.Store(filename, img)
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
