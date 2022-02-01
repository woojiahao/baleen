package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

func buildUrl(base, endpoint string, queryParams map[string]string) (*url.URL, error) {
	baseUrl, err := url.Parse(base)
	baseUrl.Path = path.Join(baseUrl.Path, endpoint)
	q := baseUrl.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	baseUrl.RawQuery = q.Encode()
	log.Printf("Url: %s\n", baseUrl.String())
	return baseUrl, err
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
