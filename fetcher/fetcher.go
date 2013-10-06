// Package fetcher defines interface of fetchers and implements a
// default fetcher.
package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Fetcher interface {
	Fetch(rawurl string) (fetchInfo *FetchInfo, err error)
}

type DefaultFetcher struct {
}

type FetchInfo struct {
	Key        *url.URL
	Contents   []byte
	StatusCode int
}

// Fetch fetches the page from URL.
func (f DefaultFetcher) Fetch(rawurl string) (fetchInfo *FetchInfo, err error) {
	info := &FetchInfo{}
	info.Key, err = url.Parse(rawurl)
	if err != nil {
		return
	}
	resp, err := http.Get(rawurl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// TODO: Handle the case.
		return
	}
	info.Contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fetchInfo = info
	return
}
