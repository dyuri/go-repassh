package main

import (
	"crypto/sha256"
	"log"
	"net"
	"os"

	"golang.org/x/term"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	h := sha256.New()
	h.Write(key.Marshal())
	log.Printf("Fingerprint: %x", h.Sum(nil))

	log.Printf("Checking host key for %s", hostname)
	khcallback, err := knownhosts.New(os.ExpandEnv("$HOME/.ssh/known_hosts"))
	if err != nil {
		log.Fatal("Failed to load known_hosts file:", err)
	}
	return khcallback(hostname, remote, key)
}

func main() {
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

	client, err := ssh.Dial("tcp", "chaos:22", config)
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
	log.Printf("Terminal size: %d x %d", w, h)

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
