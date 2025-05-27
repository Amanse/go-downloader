## Go-downloader
A very basic multi threaded downloader with plugin support to download from anywhere

### Usage

Building the plugins

```bash
go build -o normal ./plugins/normal/normal_impl.go
go build -o ffplugin ./plugins/ffplugin/ffPlugin_impl.go
```

Using plugins
```
go-downloader --file links.txt --plugin ./ffplugin
```

### Plugins info

#### FFplugin
This is the plugin for f*ckingfast[dot]com

You do not need URLs you get after going to site and clicking the button, it will auto-resolve that, you can directly links that you get from _you know where_.
It lets golang handle the multi threading part by launching the downloads in go funcs all together.

#### Normal Plugin
This assumes the URLs provided to be direct download urls and directly downloads them
