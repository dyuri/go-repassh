package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var ignoreHostKey bool
var port string

var rootCmd = &cobra.Command{
	Use:   "repassh <host>",
	Short: "repassh is a simple SSH client that uses the SSH agent",
	Long:  `repassh is a simple SSH client that uses the SSH agent ...`,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&ignoreHostKey, "ignore-host-key", false, "Ignore host key verification")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "22", "Port to connect to")
}

func Execute() {
	log.SetLevel(log.InfoLevel) // TODO
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
