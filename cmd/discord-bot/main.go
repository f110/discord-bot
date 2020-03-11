package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"

	"github.com/f110/discord-bot/pkg/bot"
	"github.com/f110/discord-bot/pkg/config"
)

var shutdownSem = make(chan struct{}, 1)

func process(args []string) error {
	name := ""
	confPath := ""
	storageHost := ""
	bucket := ""
	bucketHost := ""
	fs := pflag.NewFlagSet("discord-bot", pflag.ContinueOnError)
	fs.StringVarP(&confPath, "conf", "c", confPath, "Config file path")
	fs.StringVarP(&name, "name", "n", name, "Bot name")
	fs.StringVar(&storageHost, "storage", storageHost, "Storage endpoint")
	fs.StringVar(&bucket, "bucket", bucket, "Bucket name")
	fs.StringVar(&bucketHost, "bucket-host", bucketHost, "Hostname which accessible from the internet")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return errors.New("config file does not found")
	}

	conf, err := config.ReadConfig(confPath)
	if err != nil {
		return err
	}
	conf.Bucket = bucket
	conf.StorageHost = storageHost
	conf.BucketHost = bucketHost

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return errors.New("bot token is required")
	}

	b, err := bot.NewBot(name, token, storageHost, bucket, bucketHost, conf)
	if err != nil {
		return err
	}

	if err := b.Run(); err != nil {
		return err
	}

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for sig := range signalCh {
			switch sig {
			case syscall.SIGTERM, os.Interrupt:
				shutdown(b)
			}
		}
	}()

	b.Wait()
	return nil
}

func shutdown(p *bot.Bot) {
	shutdownSem <- struct{}{}
	defer func() {
		<-shutdownSem
	}()

	if err := p.Shutdown(); err != nil {
		log.Print(err)
	}
}

func main() {
	if err := process(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
