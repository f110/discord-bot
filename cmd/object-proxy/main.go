package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/f110/discord-bot/pkg/storage"
	"github.com/spf13/pflag"
)

var shutdownSem = make(chan struct{}, 1)

func objectProxy(args []string) error {
	bind := ""
	storageHost := ""
	bucket := ""
	fs := pflag.NewFlagSet("object-proxy", pflag.ContinueOnError)
	fs.StringVarP(&bind, "bind", "l", bind, "Listen interface and port")
	fs.StringVar(&storageHost, "storage", storageHost, "Storage endpoint")
	fs.StringVar(&bucket, "bucket", bucket, "Bucket name")
	if err := fs.Parse(args); err != nil {
		return err
	}

	st := storage.New(storageHost, bucket)
	p := storage.NewProxy(bind, st)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for sig := range signalCh {
			switch sig {
			case syscall.SIGTERM, os.Interrupt:
				shutdown(p)
			}
		}
	}()

	if err := p.Start(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func shutdown(p *storage.Proxy) {
	shutdownSem <- struct{}{}
	defer func() {
		<-shutdownSem
	}()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	if err := p.Shutdown(ctx); err != nil {
		log.Print(err)
	}
}

func main() {
	if err := objectProxy(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
}
