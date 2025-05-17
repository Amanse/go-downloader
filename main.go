package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/amanse/go-downloader/downloader"
)

func main() {
	file := flag.String("file", "", "File containing links")

	flag.Parse()
	fmt.Print(*file)
	dat, err := os.ReadFile(*file)
	if err != nil {
		panic("Error reading file")
	}
	urls := strings.Split(string(dat), "\n")
	d := downloader.New(urls)
	d.StartDownload()
}
