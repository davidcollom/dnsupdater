package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/davidcollom/dnsupdater/internal/api"
	"github.com/miekg/dns"
)

type RNDCDNSUpdater struct {
	Config *api.RNDCConfig
}

func (r *RNDCDNSUpdater) Update(ctx context.Context, hostname string, ips []net.IP) error {
	fqdn := dns.Fqdn(hostname)
	m := new(dns.Msg)
	m.SetUpdate(dns.Fqdn(r.Config.Zone))

	for _, ip := range ips {
		record := &dns.A{
			Hdr: dns.RR_Header{
				Name:   fqdn,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.Config.TTL),
			},
			A: ip,
		}
		m.Remove([]dns.RR{record})
		m.Insert([]dns.RR{record})
	}

	tsigKey := dns.Fqdn(r.Config.KeyName)
	keyAlg := r.Config.Algorithm
	secret := r.Config.KeySecret

	m.SetTsig(tsigKey, keyAlg, 300, time.Now().Unix())

	dnsClient := new(dns.Client)
	dnsClient.TsigSecret = map[string]string{tsigKey: secret}

	target := net.JoinHostPort(r.Config.Nameserver, fmt.Sprintf("%d", r.Config.Port))

	resp, _, err := dnsClient.Exchange(m, target)
	if err != nil {
		return fmt.Errorf("DNS update failed: %w", err)
	}
	if resp.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("DNS update error: %s", dns.RcodeToString[resp.Rcode])
	}

	return nil
}

func GetRNDCFactory() func(any) DNSUpdater {
	return func(cfg any) DNSUpdater {
		return &RNDCDNSUpdater{
			Config: cfg.(*api.RNDCConfig),
		}
	}
}
