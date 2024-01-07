package wg

import "net"

type WGConfig struct {
	PrivateKey string
	IP         string
	ServerKey  string
	SharedKey  string
	AllowedIPs []string
	Endpoint   string
}

type WGPeer struct {
	Interface  *WGInterface
	Name       string
	IP         net.IP
	PrivateKey string
	PublicKey  string
	SharedKey  string
	AllowedIPs []string
}

type WGInterface struct {
	PublicIP    string
	ListenPort  int
	LocalSubnet string
	PublicKey   string
	Interface   string
	AllowedIPs  []string
	Network     *net.IPNet
	RouterIP    net.IP // Used as base point for client ip allocation
	Peers       []*WGPeer
	nextIp      net.IP
}
