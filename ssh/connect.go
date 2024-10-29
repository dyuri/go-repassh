package ssh

import (
	"fmt"
	"os"
	"net"

	"github.com/charmbracelet/log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/term"
)

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	fingerprint := visualHostKeyFingerprint(key)

	log.Info("Checking host key for", "hostname", hostname)
	fmt.Printf("Host key fingerprint:\n%s", fingerprint)
	
	khcallback, err := knownhosts.New(os.ExpandEnv("$HOME/.ssh/known_hosts"))
	if err != nil {
		log.Error("Failed to load known_hosts file:", err)
	}
	
	if (khcallback != nil) {
		err = khcallback(hostname, remote, key)
		if err != nil {
			log.Error("Host key verification failed:", err)
			line := knownhosts.Line([]string{knownhosts.HashHostname(knownhosts.Normalize(hostname))}, key)
			fmt.Printf("To suppress this error, add the following line to your known_hosts file:\n%s\n", line)
		}
	}

	return err
}

func Connect(host string) {
	agentSocket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", agentSocket)
	if err != nil {
		log.Fatal("Failed to connect to the SSH agent:", err)
	}
	defer conn.Close()

	agentClient := agent.NewClient(conn)

	config := &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	conn.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session:", err)
	}
	defer session.Close()

	w, h, _ := term.GetSize(int(os.Stdin.Fd()))
	log.Info("Terminal size", "w", w, "h", h)

	termtype := os.Getenv("TERM")
	if termtype == "" {
		termtype = "xterm-256color"
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty(termtype, h, w, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	// put terminal into raw mode
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// set input and output
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatal("Failed to start shell:", err)
	}
	/*
	if err := session.Start("bash"); err != nil {
		log.Fatal("Failed to run:", err)
	}
	*/

	err = session.Wait()
	if err != nil {
		log.Fatal("Failed to wait:", err)
	}
}

