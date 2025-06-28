package dns

import (
	"context"
	"errors"
	"net"

	"github.com/davidcollom/dnsupdater/internal/api"
)

// DNSUpdater defines a pluggable interface for DNS update providers.
type DNSUpdater interface {
	Update(ctx context.Context, hostname string, ips []net.IP) error
}

func FromConfig(service api.DNSService) (DNSUpdater, string, error) {
	switch {
	case service.Cloudflare != nil:
		u := GetCloudflareFactory()(service.Cloudflare)
		return u, service.Cloudflare.Hostname, nil
	case service.RNDC != nil:
		u := GetRNDCFactory()(service.RNDC)
		return u, service.RNDC.Hostname, nil
	case service.DynDNS != nil:
		u := GetDynDNSFactory()(service.DynDNS)
		return u, service.DynDNS.Hostname, nil
	default:
		return nil, "", errors.New("no supported DNS provider configured")
	}
}
