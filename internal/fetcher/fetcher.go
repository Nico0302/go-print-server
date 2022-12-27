package fetcher

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ForwardedCookieHeader = "X-Forwarded-Cookie"
)

type Fetcher struct {
	client *http.Client
}

func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{},
	}
}

func addCookies(req *http.Request, cookies string) {
	list := strings.Split(strings.Replace(cookies, " ", "", -1), ";")
	for _, c := range list {
		name, value, _ := strings.Cut(c, "=")
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
}

func (f *Fetcher) DownloadFile(url string, cookies string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("file request error: %w", err)
	}
	if cookies != "" {
		addCookies(req, cookies)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("file download error: %w", err)
	}
	return resp, nil
}

func GetFileSize(resp *http.Response) int {
	if resp.ContentLength != 0 {
		return int(resp.ContentLength)
	}

	buf := &bytes.Buffer{}
	nRead, err := io.Copy(buf, resp.Body)
	if err != nil {
		return 0
	}

	return int(nRead)
}
