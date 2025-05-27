package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	pluginImpl "github.com/amanse/go-downloader/plugin"
	"github.com/hashicorp/go-plugin"
)

type FFPlugin struct{}

func (ff *FFPlugin) ProcessUrls(urls []string) []string {
	var result []string
	for _, url := range urls {
		log.Println("Processing URL:", url)
		result = append(result, parseFF(url))
	}
	return result
}

func parseFF(url string) string {
	client := http.Client{}
	var result string
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bs := string(body)
	re := regexp.MustCompile(`(?i)window\.open\("([^"]+)"\)`)

	// Find the first match
	match := re.FindStringSubmatch(bs)

	if len(match) > 1 {
		// Extract the URL (group 1)
		result = match[1]
		// fmt.Println("Extracted URL:", result)
	} else {
		fmt.Println("No URL found in the string.")
	}
	return result
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	ffPlugin := &FFPlugin{}

	var pluginMap = map[string]plugin.Plugin{
		"downloadPlugin": &pluginImpl.DownloadPlugin{Impl: ffPlugin},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
