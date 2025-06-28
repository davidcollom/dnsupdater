package api

import "errors"

type CloudflareConfig struct {
	APIToken string `mapstructure:"api_token"`
	ZoneID   string `mapstructure:"zone"`
	Hostname string `mapstructure:"hostname"`
	TTL      int    `mapstructure:"ttl"`
}

type RNDCConfig struct {
	Zone       string `mapstructure:"zone"`
	Hostname   string `mapstructure:"hostname"`
	Nameserver string `mapstructure:"nameserver"`
	Port       int    `mapstructure:"port"`
	KeyName    string `mapstructure:"tsig_key"`
	KeySecret  string `mapstructure:"tsig_secret"`
	Algorithm  string `mapstructure:"algorithm"`
	TTL        int    `mapstructure:"ttl"`
}

type DNSService struct {
	Cloudflare *CloudflareConfig `mapstructure:"cloudflare"`
	RNDC       *RNDCConfig       `mapstructure:"rndc"`
	DynDNS     *DynDNSConfig     `mapstructure:"dyndns"`
}

type DynDNSConfig struct {
	Hostname string `mapstructure:"hostname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	TTL      int    `mapstructure:"ttl"`
	Endpoint string `mapstructure:"endpoint"`
}

func (cfg *Config) DetermineEnabledDNSProvider() (string, any, error) {
	if cfg.DNSService.Cloudflare != nil {
		return "cloudflare", cfg.DNSService.Cloudflare, nil
	} else if cfg.DNSService.RNDC != nil {
		return "rndc", cfg.DNSService.RNDC, nil
	} else if cfg.DNSService.DynDNS != nil {
		return "dyndns", cfg.DNSService.DynDNS, nil
	}
	return "", nil, errors.New("no supported DNS provider configured")
}
