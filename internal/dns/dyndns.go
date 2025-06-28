package dns

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/davidcollom/dnsupdater/internal/api"
)

type DynDNSUpdater struct {
	Config *api.DynDNSConfig
}

func (d *DynDNSUpdater) Update(ctx context.Context, hostname string, ips []net.IP) error {
	// Most DynDNS services only support updating a single IP (usually the current one)
	if len(ips) == 0 {
		return fmt.Errorf("no IPs provided to update")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.Config.Endpoint, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Set("hostname", hostname)
	q.Set("myip", ips[0].String())
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(d.Config.Username, d.Config.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DynDNS update failed with status: %s", resp.Status)
	}

	return nil
}

func GetDynDNSFactory() func(any) DNSUpdater {
	return func(cfg any) DNSUpdater {
		return &DynDNSUpdater{
			Config: cfg.(*api.DynDNSConfig),
		}
	}
}
