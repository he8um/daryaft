package inspect

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/downloader"
	"github.com/he8um/daryaft/pkg/version"
)

type Options struct {
	URL    string
	Client *http.Client
}

func URL(ctx context.Context, options Options) (Result, error) {
	rawURL := strings.TrimSpace(options.URL)
	if err := ValidateURL(rawURL); err != nil {
		return Result{}, err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	client := options.Client
	if client == nil {
		client = downloader.DefaultHTTPClient()
	}

	head, err := request(ctx, client, http.MethodHead, rawURL, false)
	if err != nil {
		return Result{}, err
	}
	defer func() {
		_ = head.Body.Close()
	}()

	if head.StatusCode == http.StatusMethodNotAllowed {
		return inspectWithRangeFallback(ctx, client, rawURL)
	}

	result := resultFromResponse(rawURL, head)
	if !metadataSufficient(result) {
		fallback, err := inspectWithRangeFallback(ctx, client, rawURL)
		if err == nil {
			return fallback, nil
		}
	}

	return result, nil
}

func inspectWithRangeFallback(ctx context.Context, client *http.Client, rawURL string) (Result, error) {
	response, err := request(ctx, client, http.MethodGet, rawURL, true)
	if err != nil {
		return Result{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	return resultFromResponse(rawURL, response), nil
}

func request(ctx context.Context, client *http.Client, method, rawURL string, rangeProbe bool) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create inspect request: %w", err)
	}
	request.Header.Set("User-Agent", config.AppName+"/"+version.Version)
	if rangeProbe {
		request.Header.Set("Range", "bytes=0-0")
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("inspect %q: %w", rawURL, err)
	}
	return response, nil
}

func resultFromResponse(rawURL string, response *http.Response) Result {
	finalURL := rawURL
	if response.Request != nil && response.Request.URL != nil {
		finalURL = response.Request.URL.String()
	}

	contentLength, contentLengthKnown := contentLengthFromResponse(response)
	resumeSupported, resumeSupportKnown := resumeSupportFromResponse(response)

	return Result{
		URL:                rawURL,
		FinalURL:           finalURL,
		Status:             response.Status,
		StatusCode:         response.StatusCode,
		Filename:           downloader.FilenameFromResponse(finalURL, response.Header, ""),
		ContentLength:      contentLength,
		ContentLengthKnown: contentLengthKnown,
		ContentType:        response.Header.Get("Content-Type"),
		AcceptRanges:       response.Header.Get("Accept-Ranges"),
		ResumeSupported:    resumeSupported,
		ResumeSupportKnown: resumeSupportKnown,
		ETag:               response.Header.Get("ETag"),
		LastModified:       response.Header.Get("Last-Modified"),
	}
}

func contentLengthFromResponse(response *http.Response) (int64, bool) {
	if response.StatusCode == http.StatusPartialContent {
		if total, ok := contentRangeTotal(response.Header.Get("Content-Range")); ok {
			return total, true
		}
	}
	if value := strings.TrimSpace(response.Header.Get("Content-Length")); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err == nil && parsed >= 0 {
			return parsed, true
		}
	}
	return 0, false
}

func resumeSupportFromResponse(response *http.Response) (bool, bool) {
	if response.StatusCode == http.StatusPartialContent {
		return true, true
	}

	acceptRanges := strings.TrimSpace(strings.ToLower(response.Header.Get("Accept-Ranges")))
	switch acceptRanges {
	case "bytes":
		return true, true
	case "":
		return false, false
	default:
		return false, true
	}
}

func contentRangeTotal(value string) (int64, bool) {
	_, total, ok := strings.Cut(value, "/")
	if !ok || strings.TrimSpace(total) == "" || strings.TrimSpace(total) == "*" {
		return 0, false
	}
	parsed, err := strconv.ParseInt(strings.TrimSpace(total), 10, 64)
	if err != nil || parsed < 0 {
		return 0, false
	}
	return parsed, true
}

func metadataSufficient(result Result) bool {
	return result.ContentLengthKnown &&
		result.ContentType != "" &&
		result.ResumeSupportKnown &&
		result.ETag != "" &&
		result.LastModified != ""
}

func ValidateURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("scheme must be http or https")
	}
	if parsed.Host == "" {
		return fmt.Errorf("host is required")
	}
	return nil
}
