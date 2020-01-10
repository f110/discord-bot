package storage

import (
	"context"
	"io"
	"log"
	"net/http"
)

type Proxy struct {
	server *http.Server

	storage *Storage
}

func NewProxy(bind string, st *Storage) *Proxy {
	s := &http.Server{
		Addr: bind,
	}
	p := &Proxy{server: s, storage: st}
	s.Handler = p

	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	obj, err := p.storage.Get(req.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer obj.Close()

	_, err = io.Copy(w, obj)
	if err != nil {
		log.Print(err)
	}
}

func (p *Proxy) Start() error {
	log.Printf("Start %s", p.server.Addr)
	return p.server.ListenAndServe()
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}
