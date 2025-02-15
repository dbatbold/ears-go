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
	notified map[string]bool
)

func main() {
	var day int
	for {
		now := time.Now()
		if day != now.Day() {
			notified = map[string]bool{}
			day = now.Day()
		}
		run()
		time.Sleep(time.Hour)
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
			now := time.Now()
			if res != nil {
				fmt.Fprintf(os.Stderr, "%s HEAD request failed '%s': HTTP %s\n", now, monitor.Url, res.StatusCode)
			} else {
				fmt.Fprintf(os.Stderr, "%s HEAD request failed '%s'\n", now, monitor.Url)
			}
			continue
		}
		if len(monitor.Etag) > 0 {
			if res.StatusCode != http.StatusNotModified {
				monitor.print(res.Header.Get("etag"))
			} else {
				continue
			}
		}
		if monitor.Status == 0 {
			monitor.Status = 200
		}
		if len(monitor.Location) > 0 {
			location := res.Header.Get("Location")
			if http.StatusNotModified != res.StatusCode && location != monitor.Location {
				monitor.print(location)
			}
		}
		if monitor.Status != res.StatusCode {
			fmt.Println(time.Now().Format(time.RFC3339), monitor.Name)
			fmt.Println("\t", monitor.Url)
			fmt.Println("\t", "HTTP Status Code", res.StatusCode)
			continue
		}
		if len(monitor.LastModified) > 0 {
			modified := res.Header.Get("last-modified")
			if modified != monitor.LastModified {
				monitor.print(modified)
			}
		}
		if len(monitor.ContentLength) > 0 {
			length := res.Header.Get("content-length")
			if length != monitor.ContentLength {
				monitor.print(length)
			}
		}
	}

}

type Monitor struct {
	Name          string `json:"name"`
	Url           string `json:"url"`
	Status        int    `json:"status"`
	Visit         string `json:"visit"`
	Location      string `json:"location"`
	Etag          string `json:"etag"`
	LastModified  string `json:"last_modified"`
	ContentLength string `json:"content_length"`
}

func (m *Monitor) print(diff string) {
	if notified[m.Url] {
		return
	}
	now := time.Now()
	fmt.Println(now.Format(time.RFC3339), m.Name)
	fmt.Println("\t", m.Url)
	if len(m.Visit) > 0 {
		fmt.Println("\t", m.Visit)
	}
	if len(m.Etag) > 0 {
		fmt.Println("\t", m.Etag)
	}
	if len(m.LastModified) > 0 {
		fmt.Println("\t", m.LastModified)
	}
	if len(m.Location) > 0 {
		fmt.Println("\t", m.Location)
	}
	if len(m.ContentLength) > 0 {
		fmt.Println("\t", m.ContentLength)
	}
	fmt.Println("\a\t", diff)
	notified[m.Url] = true
}
