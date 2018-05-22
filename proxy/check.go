package proxy

import (
	"net/http"
	"net/url"
)

func CheckGFW(proxyAddr string) bool {

	p := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyAddr)
	}
	transport := &http.Transport{Proxy: p}
	client := &http.Client{Transport: transport}

	resp, err := client.Get("https://www.google.com")
	if err != nil || http.StatusOK != resp.StatusCode {
		return false
	} else {
		return true
	}
}
