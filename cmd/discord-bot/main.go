package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/f110/discord-bot/pkg/handler"
	"github.com/spf13/pflag"
)

var shutdownSem = make(chan struct{}, 1)

func process(args []string) error {
	name := ""
	storageHost := ""
	bucket := ""
	bucketHost := ""
	fs := pflag.NewFlagSet("discord-bot", pflag.ContinueOnError)
	fs.StringVarP(&name, "name", "n", name, "Bot name")
	fs.StringVar(&storageHost, "storage", storageHost, "Storage endpoint")
	fs.StringVar(&bucket, "bucket", bucket, "Bucket name")
	fs.StringVar(&bucketHost, "bucket-host", bucketHost, "Hostname which accessible from the internet")
	if err := fs.Parse(args); err != nil {
		return err
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return errors.New("bot token is required")
	}

	b, err := handler.New(name, token, storageHost, bucket, bucketHost)
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

func shutdown(p *handler.Bot) {
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
