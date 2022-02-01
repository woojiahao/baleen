package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

func buildUrl(base, endpoint string, queryParams map[string]string) (*url.URL, error) {
	url, err := url.Parse(base)
	url.Path = path.Join(url.Path, endpoint)
	log.Printf("Url: %s\n", url.String())
	q := url.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	url.RawQuery = q.Encode()
	return url, err
}

func Get(base, endpoint string, queryParams map[string]string) []byte {
	url, err := buildUrl(base, endpoint, queryParams)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// TODO: Check if query failed
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}
