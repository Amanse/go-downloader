package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Plugin interface {
	ProcessUrls(urls []string) []string
}

type PluginRPC struct {
	client *rpc.Client
}

func (p *PluginRPC) ProcessUrls(urls []string) []string {
	var resp []string
	err := p.client.Call("Plugin.ProcessUrls", &urls, &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type PluginRPCServer struct {
	// This is the real implementation
	Impl Plugin
}

func (s *PluginRPCServer) ProcessUrls(args *[]string, resp *[]string) error {
	*resp = s.Impl.ProcessUrls(*args)
	return nil
}

type DownloadPlugin struct {
	// Impl Injection
	Impl Plugin
}

func (p *DownloadPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PluginRPCServer{Impl: p.Impl}, nil
}

func (DownloadPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PluginRPC{client: c}, nil
}
