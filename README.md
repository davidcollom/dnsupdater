# `dnsupdater`

A dynamic DNS updater written in Go. It fetches your current public IP address from one or more
providers and updates a DNS record using the DNS service of your choice.

The updater is designed for connections whose IP address changes (home broadband, mobile, VPN) and
works with multiple DNS backends behind a single pluggable interface.

## Features

- Multiple IP address providers, queried together and de-duplicated
- Pluggable DNS update backends: **Cloudflare**, **RNDC**, and **DynDNS**
- YAML configuration via [Viper](https://github.com/spf13/viper)
- [Cobra](https://github.com/spf13/cobra) CLI with a `--config` flag

## Requirements

- Go 1.25 or newer
- A DNS provider account (Cloudflare API token, RNDC key, or DynDNS credentials)

## Installation

```bash
go build -o dnsupdater .
```

Or run directly:

```bash
go run .
```

## Configuration

Create a `config.yaml` in the working directory (or point at your own file with `--config`):

```yaml
address_providers:
  - type: ipify
  - type: icanhazip
  - type: custom_json
    opts:
      url: https://api.my-ip.io/v2/ip.json
      json_field: ip
  - type: custom
    opts:
      url: https://ifconfig.me/ip

dns_service:
  cloudflare:
    api_token: your-cloudflare-api-token
    zone: example.com
    hostname: home.example.com
    ttl: 300
```

### Address providers

| Type | Description |
| ---- | ----------- |
| `ipify` | Uses the [ipify](https://www.ipify.org/) API. |
| `icanhazip` | Uses the [icanhazip](https://icanhazip.com/) API. |
| `custom` | Fetches a plain-text URL and uses the response as the IP address. |
| `custom_json` | Fetches a JSON URL and extracts the address using the `json_field` option. |

### DNS services

Only one DNS service should be configured at a time.

**Cloudflare**

```yaml
dns_service:
  cloudflare:
    api_token: your-cloudflare-api-token
    zone: example.com
    hostname: home.example.com
    ttl: 300
```

**RNDC**

```yaml
dns_service:
  rndc:
    zone: example.com
    hostname: home.example.com
    nameserver: 127.0.0.1
    port: 53
    tsig_key: your-key
    tsig_secret: your-base64-secret
    algorithm: hmac-sha256
    ttl: 300
```

**DynDNS**

```yaml
dns_service:
  dyndns:
    hostname: home.example.com
    username: your-dyndns-username
    password: your-dyndns-password
    ttl: 300
    endpoint: https://members.dyndns.org/nic/update
```

## Usage

```bash
# Use the default config.yaml in the current directory
go run .

# Use a specific config file
go run . --config /etc/dnsupdater/config.yaml
go run . -c /etc/dnsupdater/config.yaml
```

The `--config` / `-c` flag is registered as a persistent flag on the root command, so it applies to
every subcommand.

## Layout

| Path | Description |
| ---- | ----------- |
| `main.go` | Entry point. |
| `cmd/` | Cobra command definition and config flag. |
| `internal/api/` | Configuration types and loading. |
| `internal/provider/` | Public IP address providers. |
| `internal/dns/` | DNS update providers and the `DNSUpdater` interface. |

Adding a new DNS backend means implementing the `DNSUpdater` interface and registering it in
`internal/dns/dns.go`:

```go
type DNSUpdater interface {
    Update(ctx context.Context, hostname string, ips []net.IP) error
}
```

## License

See the repository for licensing details.
