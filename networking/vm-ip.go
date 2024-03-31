package networking

import (
	"fmt"
	"net"
)

func GetVMIPs(bridgeIP string) (string, string, error) {
	// Parse the bridge IP address
	ip := net.ParseIP(bridgeIP)
	if ip == nil {
		return "", "", fmt.Errorf("invalid bridge IP address")
	}

	// Convert the IP address to IPv4
	ip = ip.To4()

	// Ensure the IP address is in the correct range for a /24 subnet
	if ip[3] != 7 {
		return "", "", fmt.Errorf("bridge IP address is not in the correct range for a /24 subnet")
	}

	// Get the network address and subnet mask
	network := ip.Mask(net.CIDRMask(24, 32))

	// Get the first two IP addresses in the subnet excluding the network address and broadcast address
	ip1 := net.IP(make([]byte, 4))
	ip2 := net.IP(make([]byte, 4))
	copy(ip1, network)
	copy(ip2, network)

	// Increment the last octet for the second and third IP address
	ip1[3] += 2
	ip2[3] += 3

	return ip1.String(), ip2.String(), nil
}
