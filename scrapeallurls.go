package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"os"
)

func downloadFile(client *http.Client, url string, filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File already exists, skipping download: %s\n", filename)
		return false
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Printf("File downloaded: %s\n", filename)
	return true
}

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	targetUrl, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fqdn := targetUrl.Scheme + "://" + targetUrl.Host // Extract FQDN

	regex := regexp.MustCompile(`="(http[^"]*|/[^"]*)`)
	matches := regex.FindAllString(string(body), -1)

	for _, match := range matches {
		link := match[2:] // The captured URL or path

		if strings.HasPrefix(link, "/") {
			link = fqdn + link // Prepend FQDN to paths starting with /
		}

		filename := strings.ReplaceAll(link, "/", "-")
		extension := filepath.Ext(link)
		filename = filename + extension

		downloaded := downloadFile(client, link, filename)
		if downloaded {
			time.Sleep(time.Second)
		}
	}
}
