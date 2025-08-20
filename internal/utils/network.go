package utils

import (
	"net"
	"time"
)

func CheckGitHubAccess() bool {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", "github.com:443", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
