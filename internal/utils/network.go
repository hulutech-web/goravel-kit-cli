package utils

import (
	"net"
	"time"
)

func CheckGiteeAccess() bool {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", "gitee.com:443", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
