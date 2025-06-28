package provider

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func init() {
	RegisterProvider("ipify", func(map[string]interface{}) Provider {
		return &IpifyProvider{}
	})
}

type IpifyProvider struct{}

func (p *IpifyProvider) GetIP(ctx context.Context) (net.IP, error) {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.RetryWaitMax = 1 * time.Second
	httpClient.Backoff = retryablehttp.LinearJitterBackoff

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org?format=json", nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return net.ParseIP(data.IP), nil
}
