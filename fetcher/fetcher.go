// Package fetcher defines interface of fetchers and implements a
// default fetcher.
package fetcher

import (
	"io/ioutil"
	"net/http"
)

type Fetcher interface {
	Fetch(url string) (body []byte, err error)
}

type DefaultFetcher struct {
}

// Fetch fetches the page from URL.
func (f DefaultFetcher) Fetch(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// TODO: Handle the case.
		return []byte{}, nil
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return
}
