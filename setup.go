package guard

import (
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

var pluginName = "guard"

func init() {
	plugin.Register(pluginName, setup)
}

func setup(caddy *caddy.Controller) error {
	caddy.Next()

	if !strings.EqualFold(pluginName, caddy.Val()) {
		return plugin.Error(pluginName, caddy.ArgErr())
	}

	config, err := CreateConfig(caddy.Dispenser)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	dnsserver.GetConfig(caddy).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return guard{
			Next:   next,
			Config: config,
		}
	})

	return nil
}
