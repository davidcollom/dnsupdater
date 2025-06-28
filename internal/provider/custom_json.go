package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func init() {
	RegisterProvider("custom_json", func(opts map[string]interface{}) Provider {
		url, _ := opts["url"].(string)
		field, _ := opts["json_field"].(string)
		return &CustomJSONHTTPProvider{URL: url, JSONField: field}
	})
}

type CustomJSONHTTPProvider struct {
	URL       string
	JSONField string // dot notation (e.g. "data.ip")
}

func (p *CustomJSONHTTPProvider) GetIP(ctx context.Context) (net.IP, error) {
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

	if p.JSONField == "" {
		return net.ParseIP(strings.TrimSpace(string(body))), nil
	}

	// Parse JSON and follow dot notation path
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("[%s]failed to unmarshal JSON: %w", p.URL, err)
	}

	fields := strings.Split(p.JSONField, ".")
	var current interface{} = raw
	for _, key := range fields {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("[%s] invalid JSON path", p.URL)
		}
		current = m[key]
	}

	strVal, ok := current.(string)
	if !ok {
		return nil, fmt.Errorf("[%s] target JSON field is not a string", p.URL)
	}

	return net.ParseIP(strings.TrimSpace(strVal)), nil
}
