package command

import (
	"github.com/f110/discord-bot/pkg/config"
	"github.com/f110/discord-bot/pkg/handler"
)

var Manager = newManager()

type Plugin interface {
	Enable(conf *config.Config) error
	Subscribe(handler *handler.Handler)
}

type manager struct {
	plugins map[string]Plugin
}

func newManager() *manager {
	return &manager{plugins: make(map[string]Plugin)}
}

func (m *manager) Register(name string, plugin Plugin) {
	m.plugins[name] = plugin
}

func (m *manager) Fetch(name string) (Plugin, bool) {
	p, ok := m.plugins[name]
	if !ok {
		return nil, false
	}

	return p, true
}
