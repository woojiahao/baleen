package api

import (
	"bytes"
	"encoding/json"
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

func Post(base, endpoint string, headers map[string]string, body map[string]string) []byte {
	url, err := buildUrl(base, endpoint, map[string]string{})
	if err != nil {
		panic(err)
	}

	encodedBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(encodedBody))
	if err != nil {
		panic(err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return respBody
}
