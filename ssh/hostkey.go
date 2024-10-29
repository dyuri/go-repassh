package ssh

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

func visualHostKeyFingerprint(key ssh.PublicKey) string {
	sb := strings.Builder{}

	h := sha256.New()
	h.Write(key.Marshal())
	sum := h.Sum(nil)

	for i := 0; i < len(sum); i++ {
		if i % 8 == 0 && i != 0 {
			sb.WriteString("\n")
		}
		c := sum[i]
		bg := 16 + (c % 216)
		var fg int
		if (bg - 16) % 36 < 12 {
			fg = 255
		} else {
			fg = 232
		}
		sb.WriteString(fmt.Sprintf("\033[38;5;%dm\033[48;5;%dm%02x\033[0m", fg, bg, c))
	}
	sb.WriteString("\n")

	return sb.String()
}

