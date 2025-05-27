package main

import (
	pluginImpl "github.com/amanse/go-downloader/plugin"
	"github.com/hashicorp/go-plugin"
)

type NormalPlugin struct{}

func (ff *NormalPlugin) ProcessUrls(urls []string) []string {
	return urls
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	ffPlugin := &NormalPlugin{}

	var pluginMap = map[string]plugin.Plugin{
		"downloadPlugin": &pluginImpl.DownloadPlugin{Impl: ffPlugin},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
