package util

import "net"

const (
	localhost = "127.0.0.1"
)

// ExternalIP 获取 interface 的 IPv4 地址
func ExternalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return localhost
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			// 该 interface 是关闭的
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return localhost
		}
		for _, addr := range addrs {
			ip := getIPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}
	return localhost
}

// getIPFromAddr 获取 addr 中的 IPv4 地址
func getIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	return ip
}
