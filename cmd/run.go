package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/davidcollom/dnsupdater/internal/api"
	"github.com/davidcollom/dnsupdater/internal/dns"
	"github.com/davidcollom/dnsupdater/internal/provider"
)

func init() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the IP discovery and DNS updater",
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := api.Load(cfgPath)
			if err != nil {
				log.Fatalf("failed to load config: %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			ipList, err := provider.ResolveAll(ctx, cfg.AddressProviders)
			if err != nil {
				log.Fatalf("IP resolution failed: %v", err)
			}

			u, hostname, err := dns.FromConfig(cfg.DNSService)
			if err != nil {
				log.Fatalf("DNS provider setup failed: %v", err)
			}

			if err := u.Update(ctx, hostname, ipList); err != nil {
				log.Fatalf("DNS update failed: %v", err)
			}

			log.Println("DNS update successful")
		},
	}

	runCmd.Flags().StringP("config", "c", "config.yaml", "Path to config file")
	rootCmd.AddCommand(runCmd)
}
