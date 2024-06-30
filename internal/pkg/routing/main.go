// Package routing for find preferred ip to destination ip.
package routing

import (
	"fmt"
	"net"

	"github.com/google/gopacket/routing"
)

type router struct {
	router routing.Router
}

func New() (*router, error) {
	r, err := routing.New()
	if err != nil {
		return nil, fmt.Errorf("router create error:%w", err)
	}

	return &router{
		router: r,
	}, nil
}

func (r *router) Route(dst net.IP) (net.IP, error) {
	_, _, src, err := r.router.Route(dst)
	if err != nil {
		return nil, fmt.Errorf("route error:%w", err)
	}

	return src, nil
}

func ParseHost(a string) (net.IP, error) {
	host, _, err := net.SplitHostPort(a)
	if err != nil {
		return nil, fmt.Errorf("split address error:%w", err)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, fmt.Errorf("bad host address:%v", host)
	}

	return ip, nil
}
