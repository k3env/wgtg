package api

import (
	"fmt"
	"github.com/k3env/wgtg/errors"
	"github.com/k3env/wgtg/wg"
	"net"
	"strings"
)

func (api *MikrotikAPI) loadInterfaces() (err error) {
	var publicIP = ""
	if api.mikrotik.publicAddr != "" {
		if api.logger != nil {
			api.logger.Info().Str("address", api.mikrotik.publicAddr).Msg("Using endpoint address from config")
		}
		publicIP = api.mikrotik.publicAddr
	} else {
		pub, err := api.getPublicIP()
		if err != nil {
			return err
		}
		publicIP = pub.String()
	}

	if publicIP == "" {
		parts := strings.Split(api.mikrotik.endpoint, ":")
		publicIP = parts[0]
		if api.logger != nil {
			api.logger.Info().Msg("No public ip detected, using MikroTIK api IP address")
		}
	}

	res, err := api.client.RunArgs([]string{"/interface/wireguard/print"})
	if err != nil {
		return err
	}
	if len(res.Re) == 0 {
		return errors.NoWGInterfaces
	}
	for _, cfg := range res.Re {
		var wgif wg.WGInterface
		err = wgif.Parse(*cfg)
		if err != nil {
			if api.logger != nil {
				api.logger.Error().Str("interface", cfg.Map["name"]).Err(err).Msg("Error while parsing interface config")
			}
			continue
		}
		ip, network, err := api.getInterfaceIp(cfg.Map["name"])
		if err != nil {
			if api.logger != nil {
				api.logger.Error().Str("interface", cfg.Map["name"]).Err(err).Msg("Error while parsing interface config")
			}
			continue
		}
		err = wgif.SetNetworks(ip, network, publicIP)
		if err != nil {
			if api.logger != nil {
				api.logger.Error().Str("interface", cfg.Map["name"]).Err(err).Msg("Error while parsing interface config")
			}
			continue
		}
		api.Interfaces[wgif.Interface] = &wgif
	}
	return nil
}

func (api *MikrotikAPI) loadPeers() (err error) {
	res, err := api.client.Run("/interface/wireguard/peers/print")
	if err != nil {
		return err
	}
	for _, peer := range res.Re {
		data := peer.Map
		ifName := data["interface"]

		if api.Interfaces[ifName].Interface == ifName {
			parts := strings.Split(data["comment"], "|")
			user := ""
			pk := ""
			var ip net.IP
			if len(parts) < 3 {
				if api.logger != nil {
					api.logger.Warn().Msg("Error while parsing peer comment, skipping")
				}
				continue
			}
			if len(parts) >= 3 {
				pk = parts[2]
			}
			if len(parts) >= 2 {
				ip = net.ParseIP(parts[1])
			}
			if len(parts) >= 1 {
				user = parts[0]
			}
			shared := data["preshared-key"]
			key := data["public-key"]
			ipList := strings.Split(data["allowed-address"], ",")
			wgPeer := &wg.WGPeer{
				Interface:  api.Interfaces[ifName],
				Name:       user,
				IP:         ip,
				PrivateKey: pk,
				PublicKey:  key,
				SharedKey:  shared,
				AllowedIPs: ipList,
			}
			api.Interfaces[ifName].ImportPeer(*wgPeer)
		}
	}
	return nil
}

func (api *MikrotikAPI) addNewPeer(wgInterface *wg.WGInterface, name string, usePsk bool) (*wg.WGPeer, error) {
	peer, _ := wgInterface.AddPeer(name, usePsk)
	cmd := fmt.Sprintf("/interface/wireguard/peers/add =public-key=%s =interface=%s =allowed-address=%s =comment=%s|%s|%s", peer.PublicKey, wgInterface.Interface, strings.Join(peer.AllowedIPs, ","), peer.Name, peer.IP.String(), peer.PrivateKey)
	args := strings.Split(cmd, " ")
	_, err := api.client.RunArgs(args)
	if err != nil {
		return nil, err
	}
	return peer, nil
}

func (api *MikrotikAPI) getPublicIP() (net.IP, error) {
	res, err := api.client.Run("/ip/cloud/print")
	if err != nil {
		return nil, err
	}
	if len(res.Re) == 0 {
		return nil, errors.IPNotFoundError
	}
	addr := res.Re[0].Map["public-address"]
	ip := net.ParseIP(addr)
	return ip, nil
}

func (api *MikrotikAPI) getInterfaceIp(ifname string) (net.IP, *net.IPNet, error) {
	args := []string{"/ip/address/print", fmt.Sprintf("?=interface=%s", ifname)}
	res, err := api.client.RunArgs(args)
	if err != nil {
		return nil, nil, err
	}
	if len(res.Re) == 0 {
		return nil, nil, errors.IPNotFoundError
	}
	addr := res.Re[0].Map["address"]
	return net.ParseCIDR(addr)
}
