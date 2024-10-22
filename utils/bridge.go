package utils

import (
	"fmt"
	"net"

	"math/rand"

	"github.com/sagoresarker/firecracker-first-vmm/types"
)

// it return bridge ip address and gateway ip address
func generateBridgeIPAddress(startRange, endRange string) (string, string, error) {
	// Parse start and end IP addresses
	startIP := net.ParseIP(startRange).To4()
	endIP := net.ParseIP(endRange).To4()

	if startIP == nil || endIP == nil {
		return "", "", fmt.Errorf("invalid IP address range")
	}

	// Convert IP addresses to integers
	start := int(startIP[0])<<24 | int(startIP[1])<<16 | int(startIP[2])<<8 | int(startIP[3])
	end := int(endIP[0])<<24 | int(endIP[1])<<16 | int(endIP[2])<<8 | int(endIP[3])

	// Generate a random IP address within the range
	randomIP := make(net.IP, 4)
	ipInt := rand.Intn((end-start)/256) + start // Divide by 256 to ensure the last octet is always 0

	randomIP[0] = byte(ipInt >> 24 & 0xFF)
	randomIP[1] = byte(ipInt >> 16 & 0xFF)
	randomIP[2] = byte(ipInt >> 8 & 0xFF)
	randomIP[3] = 7 // Set the last octet to 1

	bridgeIp := randomIP.String()

	randomIP[3] = 1

	gateway_ip := randomIP.String()

	return bridgeIp, gateway_ip, nil
}

func GenerateValue() (types.BridgeDetails, error) {
	fmt.Println("Generate a value for bridge-name, user-id and ip-address")

	startRange := "10.0.0.0"
	endRange := "10.255.255.255"

	userID := getUserID()
	bridgeName := userID + "-" + "br"

	bridge_ip_address, gateway_ip, err := generateBridgeIPAddress(startRange, endRange)

	if err != nil {
		fmt.Println("Error Generating IP adress:", err)
		return types.BridgeDetails{}, err
	}

	bridge_ip_address = bridge_ip_address + "/24"

	return types.BridgeDetails{
		BridgeName:      bridgeName,
		UserID:          userID,
		BridgeIPAddress: bridge_ip_address,
		BridgeGatewayIP: gateway_ip,
	}, nil

}
