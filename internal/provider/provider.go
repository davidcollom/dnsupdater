package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/davidcollom/dnsupdater/internal/api"
)

type Provider interface {
	GetIP(ctx context.Context) (net.IP, error)
}

var providerRegistry = map[string]func(map[string]interface{}) Provider{}

func RegisterProvider(name string, factory func(map[string]interface{}) Provider) {
	if _, exists := providerRegistry[name]; exists {
		panic(fmt.Sprintf("Provider already registered: %s", name))
	}
	providerRegistry[name] = factory
}

func GetFactory(name string) func(map[string]interface{}) Provider {
	return providerRegistry[name]
}

func ResolveAll(ctx context.Context, configs []api.ProviderConfig) ([]net.IP, error) {
	var (
		wg      sync.WaitGroup
		mutex   sync.Mutex
		results []net.IP
		seen    = map[string]struct{}{}
	)

	for _, cfg := range configs {
		cfg := cfg // capture for closure
		wg.Add(1)

		go func() {
			defer wg.Done()

			factory := GetFactory(cfg.Type)
			if factory == nil {
				log.Printf("Unknown provider: %s", cfg.Type)
				return
			}

			p := factory(cfg.Opts)

			log.Printf("Resolving IP from provider: %s", cfg.Type)
			ip, err := p.GetIP(ctx)
			if err != nil {
				log.Printf("Error from provider %s: %v", cfg.Type, err)
				return
			}
			log.Printf("Resolved IP from provider %s: %s", cfg.Type, ip.String())

			if ip == nil {
				return
			}

			mutex.Lock()
			defer mutex.Unlock()
			if _, ok := seen[ip.String()]; !ok {
				seen[ip.String()] = struct{}{}
				results = append(results, ip)
			}
		}()
	}

	wg.Wait()

	if len(results) == 0 {
		return nil, errors.New("no valid IPs resolved from any provider")
	}

	return results, nil
}
