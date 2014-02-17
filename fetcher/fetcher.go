// Package fetcher defines interface of fetchers and implements a
// default fetcher.
package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"code.google.com/p/goprotobuf/proto"
)

type Fetcher interface {
	Fetch(rawurl string) (fetchInfo *FetchInfo, err error)
}

type DefaultFetcher struct {
}

// Fetch fetches the page from URL.
func (f DefaultFetcher) Fetch(rawurl string) (fetchInfo *FetchInfo, err error) {
	info := &FetchInfo{}
	var key *url.URL
	key, err = url.Parse(rawurl)
	if err != nil {
		return
	}
	info.Url = proto.String(key.String())
	resp, err := http.Get(rawurl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		info.StatusCode = proto.Int32(int32(resp.StatusCode))
		fetchInfo = info
		return
	}
	info.Contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fetchInfo = info
	return
}
