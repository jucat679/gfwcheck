package proxy

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func Check(proxyAddr string) bool {

	p := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyAddr)
	}
	transport := &http.Transport{Proxy: p}
	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	resp, err := client.Get("https://www.google.com")
	if err != nil || http.StatusOK != resp.StatusCode {
		log.Printf("Failed connect to Google using proxy [%s]\n", proxyAddr)
		return false
	} else {
		log.Printf("Connect to Google using proxy [%s] success!\n", proxyAddr)
		return true
	}
}
