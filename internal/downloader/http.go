package downloader

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/pkg/version"
)

const (
	DefaultDialTimeout           = 30 * time.Second
	DefaultTLSHandshakeTimeout   = 10 * time.Second
	DefaultResponseHeaderTimeout = 30 * time.Second
	DefaultIdleConnTimeout       = 90 * time.Second
)

func defaultHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:   DefaultDialTimeout,
		KeepAlive: 30 * time.Second,
	}).DialContext
	transport.TLSHandshakeTimeout = DefaultTLSHandshakeTimeout
	transport.ResponseHeaderTimeout = DefaultResponseHeaderTimeout
	transport.IdleConnTimeout = DefaultIdleConnTimeout

	return &http.Client{
		Transport: transport,
	}
}

func newRequest(rawURL string) (*http.Request, error) {
	return newRequestWithContext(context.Background(), rawURL)
}

func newRequestWithContext(ctx context.Context, rawURL string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", config.AppName+"/"+version.Version)
	return request, nil
}
