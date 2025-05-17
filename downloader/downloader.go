package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	Urls       []string `json:urls`
	TotalSize  int      `json:totalSize`
	reqObjs    []*http.Request
	httpClient http.Client
}

func New(urls []string) Downloader {
	d := Downloader{Urls: urls}
	done := make(chan bool)
	go showLoading("Getting total size", done)
	d.populateReqObjects()
	d.populateTotalSize()
	done <- true
	return d
}

func showLoading(text string, done chan bool) {
	chars := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r" + text + " done.   \n") // Clear the loading and show "done"
			return
		default:
			fmt.Printf("\r%s %s", text, chars[i%len(chars)])
			time.Sleep(100 * time.Millisecond) // Adjust speed as needed
			i++
		}
	}
}

func (d *Downloader) jobStart(idx int, wg *sync.WaitGroup, bar *progressbar.ProgressBar) error {
	defer wg.Done()
	destination := "file-" + strconv.Itoa(idx) + ".zip"
	resp, err := d.httpClient.Do(d.reqObjs[idx])
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s, status code: %d", d.reqObjs[idx].URL, resp.StatusCode)
	}

	out, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destination, err)
	}
	defer out.Close()

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read response body for %s: %w", d.reqObjs[idx].URL, err)
		}
		if n == 0 {
			break
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", destination, err)
		}

		bar.Add(n)
	}

	fmt.Printf("Downloaded %s to %s\n", idx, destination)
	return nil
}

func (d *Downloader) StartDownload() {
	var wg sync.WaitGroup
	bar := progressbar.DefaultBytes(int64(d.TotalSize), "Downloading All Files")

	for idx := range d.reqObjs {
		wg.Add(1)
		go func() {
			err := d.jobStart(idx, &wg, bar)
			if err != nil {
				fmt.Println("Download error:", err)
			}
		}()
	}

	wg.Wait()
	fmt.Println("All jobs completed.")
}

func (d *Downloader) populateReqObjects() {
	d.httpClient = http.Client{}
	for _, url := range d.Urls {
		url = d.parseFF(url)
		u := strings.TrimSuffix(url, "\r")
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			log.Println("Error getting http req")
		}
		d.reqObjs = append(d.reqObjs, req)
	}
}

func (d *Downloader) populateTotalSize() {
	size := 0

	for _, reqObj := range d.reqObjs {
		res, err := d.httpClient.Do(reqObj)
		if err != nil {
			log.Printf("Cannot do http req")
		}
		size += int(res.ContentLength)
	}

	d.TotalSize = size
}

func (d *Downloader) parseFF(url string) string {
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

	if match != nil && len(match) > 1 {
		// Extract the URL (group 1)
		result = match[1]
		// fmt.Println("Extracted URL:", result)
	} else {
		fmt.Println("No URL found in the string.")
	}
	return result
}
