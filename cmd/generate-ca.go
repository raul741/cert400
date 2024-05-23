/*
Copyright © 2024 raul <raul@bulgariu.xyz>
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generatecaCmd = &cobra.Command{
	Use:   "generate-ca",
	Short: "Generate CA certificates for the webserver to use",
	Long:  `Generate CA certificates for the webserver to use`,
	Run: func(cmd *cobra.Command, args []string) {
		generateCA()
	},
}

func init() {
	rootCmd.AddCommand(generatecaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
