package util

import (
	"github.com/k3env/wgtg/errors"
	"net"
)

func NextIP(ip net.IP, inc uint) (net.IP, error) {
	i := ip.To4()
	if i == nil {
		return nil, errors.InvalidIPError
	}
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	if v3 == 255 {
		return NextIP(net.IPv4(v0, v1, v2, v3), 2)
	}
	if v3 == 0 {
		return NextIP(net.IPv4(v0, v1, v2, v3), 1)
	}
	return net.IPv4(v0, v1, v2, v3), nil
}
