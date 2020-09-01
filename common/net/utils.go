package net

import (
	"fmt"
	appenv "github.com/vazmin/eagle-eye-kratos/common/env"
	"net"
)

func GetIP() (ips []string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

func GetGRPCAddrs() ([]string, error) {
	ips, err := GetIP()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ips); i++ {
		ips[i] = fmt.Sprintf("grpc://%s:%s", ips[i], appenv.GRPC_PORT)
	}
	return ips, nil
}
