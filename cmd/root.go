package cmd

import (
	"github.com/dyuri/go-repassh/ssh"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "repassh <host>",
	Short: "repassh is a simple SSH client that uses the SSH agent",
	Long:  `repassh is a simple SSH client that uses the SSH agent ...`,
	Args:  cobra.MinimumNArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			ssh.Connect(args[0])
		} else {
			log.Fatal("No host specified")
		}
	},
}

func init() {
	// TODO add flags
};

func Execute() {
	log.SetLevel(log.InfoLevel) // TODO
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
