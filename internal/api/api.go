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
	return requestWithHeadersAndBody("POST", base, endpoint, headers, body)
}

func Patch(base, endpoint string, headers map[string]string, body map[string]string) []byte {
	return requestWithHeadersAndBody("PATCH", base, endpoint, headers, body)
}

func requestWithHeadersAndBody(method, base, endpoint string, headers, body map[string]string) []byte {
	url, err := buildUrl(base, endpoint, map[string]string{})
	if err != nil {
		panic(err)
	}

	// TODO: Talk about having to parse any body as a raw messsage to avoid escaping the " in a nested JSON - in UpdateDatabaseProperties
	rawBody := make(map[string]*json.RawMessage)
	for k, v := range body {
		rawProperty := json.RawMessage(v)
		rawBody[k] = &(rawProperty)
	}
	encodedBody, _ := json.Marshal(rawBody)
	req, err := http.NewRequest(method, url.String(), bytes.NewBuffer(encodedBody))
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
