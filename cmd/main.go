package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ip-dns-updater",
	Short: "Fetch public IP and update DNS record",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "Path to config file")
}
