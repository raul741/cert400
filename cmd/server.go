/*
Copyright © 2024 raul <raul@bulgariu.xyz>
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start cert400 web server",
	Long:  `Start cert400 web server`,
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// TODO: Allow choosing CA certificates by parameter
	// serverCmd.Flags().StringP("cert", "c", "$HOME/.config/cert400/server.crt", "CA certificate")
	// serverCmd.Flags().StringP("key", "k", "$HOME/.config/cert400/server.key", "CA private key")
}
