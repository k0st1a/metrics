// Package routing for find preferred ip to destination ip.
package routing

import (
	"fmt"
	"net"

	"github.com/google/gopacket/routing"
)

func Route(dst string) (string, error) {
	router, err := routing.New()
	if err != nil {
		return "", fmt.Errorf("router create error:%w", err)
	}

	host, _, err := net.SplitHostPort(dst)
	if err != nil {
		return "", fmt.Errorf("spliet dst error:%w", err)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return "", fmt.Errorf("bad dst address")
	}

	_, _, src, err := router.Route(ip)
	if err != nil {
		return "", fmt.Errorf("route error:%w", err)
	}

	return src.String(), nil
}
