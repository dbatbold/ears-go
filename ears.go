package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"os"
)

func main() {

	earsJson, err := ioutil.ReadFile("ears.json")
	if err != nil {
		panic(err)
	}
	var list []Monitor
	if err := json.Unmarshal(earsJson, &list); err != nil {
		panic(err)
	}

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	for _, monitor := range list {
		req, err := http.NewRequest("HEAD", monitor.Url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("accept", "*/*")
		req.Header.Set("User-Agent", "curl/8.7.1")
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "HEAD request failed '%s': HTTP %s\n", monitor.Url, res.StatusCode)
			continue
		}
		// fmt.Println(res.StatusCode)
		// for k, h := range res.Header {
		// 	fmt.Println(k, h)
		// }
		if len(monitor.Location) > 0 {
			location := res.Header.Get("Location")
			if location != monitor.Location {
				fmt.Println(monitor.Name + ":")
				fmt.Println("\t" + monitor.Location)
				fmt.Println("\t" + location)
			}
		}
	}
	
}

type Monitor struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Location string `json:"location"`
}
