package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/davidcollom/dnsupdater/internal/api"
)

type CloudflareUpdater struct {
	APIToken string
	ZoneID   string
	Record   string
	TTL      int
}

func (c *CloudflareUpdater) Update(ctx context.Context, hostname string, ips []net.IP) error {
	// Placeholder for Cloudflare API logic
	fmt.Printf("[Cloudflare] Updating %s with IPs: %v\n", hostname, ips)
	return nil
}
func GetCloudflareFactory() func(any) DNSUpdater {
	return func(cfg any) DNSUpdater {
		return &CloudflareUpdater{
			APIToken: cfg.(*api.CloudflareConfig).APIToken,
			ZoneID:   cfg.(*api.CloudflareConfig).ZoneID,
			Record:   cfg.(*api.CloudflareConfig).Hostname,
			TTL:      cfg.(*api.CloudflareConfig).TTL,
		}
	}
}
