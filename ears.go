package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	notified = map[string]bool{}
)

func main() {
	date := time.Now()
	for {
		run()
		time.Sleep(time.Hour)
		if date.Day() != time.Now().Day() {
			notified = map[string]bool{}
			date = time.Now()
		}
	}
}

func parseConfig() (list []Monitor) {
	earsJson, err := ioutil.ReadFile("ears.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(earsJson, &list); err != nil {
		panic(err)
	}
	return
}

func run() {
	list := parseConfig()

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// prevent redirects
		return http.ErrUseLastResponse
	}

	for _, monitor := range list {
		req, err := http.NewRequest("HEAD", monitor.Url, nil)
		if err != nil {
			panic(err)
		}
		if len(monitor.Etag) > 0 {
			req.Header.Set("If-None-Match", monitor.Etag)
		}
		// req.Header.Set("User-Agent", "curl/8.7.1")
		res, err := client.Do(req)
		if err != nil {
			if res != nil {
				fmt.Fprintf(os.Stderr, "HEAD request failed '%s': HTTP %s\n", monitor.Url, res.StatusCode)
			} else {
				fmt.Fprintf(os.Stderr, "HEAD request failed '%s'\n", monitor.Url)
			}
			continue
		}
		// fmt.Println(res.StatusCode)
		// for k, h := range res.Header {
		// 	fmt.Println(k, h)
		// }
		if len(monitor.Location) > 0 {
			location := res.Header.Get("Location")
			if location != monitor.Location {
				monitor.print(location)
			}
		}
		if len(monitor.Etag) > 0 && res.StatusCode != http.StatusNotModified {
			monitor.print(res.Header.Get("etag"))
		}
		if len(monitor.LastModified) > 0 {
			modified := res.Header.Get("last-modified")
			if modified != monitor.LastModified {
				monitor.print(modified)
			}
		}
		if len(monitor.Redirect) > 0 {
			location := res.Header.Get("location")
			if location != monitor.Redirect {
				monitor.print(location)
			}
		}
	}

}

type Monitor struct {
	Name         string `json:"name"`
	Url          string `json:"url"`
	Location     string `json:"location"`
	Etag         string `json:"etag"`
	LastModified string `json:"last_modified"`
	Redirect     string `json:"redirect"`
}

func (m *Monitor) print(diff string) {
	if notified[m.Url] {
		return
	}
	now := time.Now()
	fmt.Println(now.Format(time.RFC3339), m.Name, m.Url)
	if len(m.Redirect) > 0 {
		fmt.Println("\t" + m.Redirect)
		fmt.Println("\t" + diff)
		notified[m.Url] = true
	}
	if len(m.LastModified) > 0 {
		fmt.Println("\t" + m.LastModified)
		fmt.Println("\t" + diff)
		notified[m.Url] = true
	}
	if len(m.Location) > 0 {
		fmt.Println("\t" + m.Location)
		fmt.Println("\t" + diff)
		notified[m.Url] = true
	}
}
