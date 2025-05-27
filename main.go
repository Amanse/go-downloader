package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/amanse/go-downloader/downloader"
	pluginImpl "github.com/amanse/go-downloader/plugin"
	"github.com/hashicorp/go-plugin"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"downloadPlugin": &pluginImpl.DownloadPlugin{Impl: &pluginImpl.PluginRPC{}},
}

func main() {
	file := flag.String("file", "", "File containing links")
	pluginPath := flag.String("plugin", "", "Path to plugin")

	flag.Parse()
	fmt.Print(*file)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(*pluginPath),
	})

	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client created")

	// Request the plugin
	raw, err := rpcClient.Dispense("downloadPlugin")
	if err != nil {
		log.Fatal(err)
	}

	pl := raw.(pluginImpl.Plugin)
	dat, err := os.ReadFile(*file)
	if err != nil {
		panic("Error reading file")
	}
	urls := strings.Split(string(dat), "\n")
	fmt.Println(urls)
	ll := pl.ProcessUrls(urls)
	d := downloader.New(ll)
	d.StartDownload()
}
