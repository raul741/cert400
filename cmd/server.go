/*
Copyright © 2024 raul
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringP("cert", "c", "$HOME/.config/cert400/server.crt", "CA certificate")
	serverCmd.Flags().StringP("key", "k", "$HOME/.config/cert400/server.key", "CA private key")
}
func setServerParams(cmd *cobra.Command) error {
	// Get "-f" flag parameters and configure them if not empty
	// parameterFile, err := cmd.Flags().GetString("file")
	// if err != nil {
	// 	return err
	// }
	// if parameterFile != "" {
	// 	playIndividualFile = true
	// 	filePath = parameterFile
	// }
	//
	// // Get "-d" flag parameters and configure them if not empty
	// parameterDir, err := cmd.Flags().GetString("directory")
	// if err != nil {
	// 	return err
	// }
	// if parameterDir != "" {
	// 	recurseFolders = true
	// 	folderPath = parameterDir
	// }
	//
	// debugParameter, err := cmd.Flags().GetBool("debug")
	// if err != nil {
	// 	return err
	// }
	// if debugParameter == true {
	// 	debuggingMode = true
	// }
	return nil
}
