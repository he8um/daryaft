package downloader

import (
	"net/http"
	"time"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/pkg/version"
)

const DefaultHTTPTimeout = 30 * time.Second

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: DefaultHTTPTimeout,
	}
}

func newRequest(rawURL string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", config.AppName+"/"+version.Version)
	return request, nil
}
