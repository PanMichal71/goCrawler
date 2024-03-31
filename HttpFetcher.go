package main

import (
	"io"
	"log"
	"net/http"
)

type HTTPFetcher struct {
}

func (f *HTTPFetcher) FetchHTML(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {

		log.Println(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(body), nil
}
