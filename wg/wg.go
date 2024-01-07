package wg

import (
	"fmt"
	"github.com/go-routeros/routeros/proto"
	"github.com/k3env/wgtg/errors"
	"github.com/k3env/wgtg/util"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"strconv"
)

func (wgif *WGInterface) Parse(row proto.Sentence) error {
	port, err := strconv.Atoi(row.Map["listen-port"])
	if err != nil {
		return err
	}
	wgif.ListenPort = port
	wgif.PublicKey = row.Map["public-key"]
	wgif.Interface = row.Map["name"]
	wgif.Peers = make([]*WGPeer, 0)
	return nil
}

func (wgif *WGInterface) SetNetworks(localIP net.IP, localNet *net.IPNet, publicIP string) (err error) {
	poolMask, _ := localNet.Mask.Size()
	localSubnet := fmt.Sprintf("%s/%d", localNet.IP.String(), poolMask)
	wgif.AllowedIPs = []string{"0.0.0.0/0", localSubnet}
	wgif.RouterIP = localIP
	wgif.LocalSubnet = localSubnet
	wgif.Network = localNet
	wgif.nextIp, err = util.NextIP(localIP, 1)
	if err != nil {
		return
	}
	wgif.PublicIP = publicIP
	return
}

func (wgif *WGInterface) AddPeer(name string, usePsk bool) (*WGPeer, error) {
	if wgif == nil {
		return nil, errors.NoWGInterfaces
	}
	if !wgif.Network.Contains(wgif.nextIp) {
		return nil, errors.NoAllocableIPsError
	}
	psk := ""
	if usePsk {
		px, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, err
		}
		psk = px.String()
	}
	kx, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	pub := kx.PublicKey().String()
	key := kx.String()
	peer := WGPeer{
		Interface:  wgif,
		Name:       name,
		IP:         wgif.nextIp,
		PrivateKey: key,
		PublicKey:  pub,
		SharedKey:  psk,
		AllowedIPs: wgif.AllowedIPs,
	}
	if err != nil {
		return nil, err
	}
	wgif.ImportPeer(peer)
	return &peer, nil
}

func (wgif *WGInterface) ImportPeer(peer WGPeer) {
	wgif.Peers = append(wgif.Peers, &peer)
	var next net.IP = wgif.nextIp
	for peer.IP.Equal(next) {
		next, _ = util.NextIP(next, 1)
	}
	wgif.nextIp = next
}

func (p *WGPeer) ExportConfig() WGConfig {
	cfg := WGConfig{
		PrivateKey: p.PrivateKey,
		IP:         p.IP.String(),
		ServerKey:  p.Interface.PublicKey,
		SharedKey:  p.SharedKey,
		AllowedIPs: p.AllowedIPs,
		Endpoint:   fmt.Sprintf("%s:%d", p.Interface.PublicIP, p.Interface.ListenPort),
	}
	return cfg
}
