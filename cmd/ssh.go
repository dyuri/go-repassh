package cmd

import (
	"strings"

	"github.com/dyuri/go-repassh/ssh"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh <host>",
	Short: "Connect to <host> via SSH",
	Long:  `Connect to <host> via SSH ...`,
	Run:   func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			ssh.Connect(strings.Join([]string{args[0], port}, ":"), ignoreHostKey)
		} else {
			log.Fatal("No host specified")
		}
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}

