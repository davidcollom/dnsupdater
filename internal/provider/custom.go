package provider

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func init() {
	RegisterProvider("custom", func(opts map[string]interface{}) Provider {
		url, _ := opts["url"].(string)
		return &CustomHTTPProvider{URL: url}
	})
}

type CustomHTTPProvider struct {
	URL string
}

func (p *CustomHTTPProvider) GetIP(ctx context.Context) (net.IP, error) {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.RetryWaitMax = 1 * time.Second
	httpClient.Backoff = retryablehttp.LinearJitterBackoff

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(strings.TrimSpace(string(body))), nil
}
