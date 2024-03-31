package networking

import (
	"fmt"
	"math/rand"
	"net"
	"os"

	"github.com/coreos/go-iptables/iptables"
	"github.com/sagoresarker/firecracker-first-vmm/users"
	"github.com/vishvananda/netlink"
)

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

func generateValue() (bridgeName string, userID string, bridge_ip_address string, gateway_ip string) {
	fmt.Println("Generate a value for bridge-name, user-id and ip-address")

	startRange := "10.0.0.0"
	endRange := "10.255.255.255"

	userID = users.GetUserID()
	bridgeName = userID + "-" + "br"

	bridge_ip_address, gateway_ip, err := generateBridgeIPAddress(startRange, endRange)

	if err != nil {
		fmt.Println("Error Generating IP adress:", err)
		return
	}

	bridge_ip_address = bridge_ip_address + "/24"

	return bridgeName, userID, bridge_ip_address, gateway_ip

}

func createBridge(bridgeName string, ipAddress string) error {
	// Create a new bridge
	bridge := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: bridgeName,
		},
	}
	if err := netlink.LinkAdd(bridge); err != nil {
		return fmt.Errorf("failed to create bridge: %v", err)
	}

	// Assign IP address to the bridge
	addr, err := netlink.ParseAddr(ipAddress)
	if err != nil {
		// Clean up the bridge if IP assignment fails
		if delErr := netlink.LinkDel(bridge); delErr != nil {
			return fmt.Errorf("failed to delete bridge after IP assignment failure: %v", err)
		}
		return fmt.Errorf("failed to parse IP address: %v", err)
	}
	if err := netlink.AddrAdd(bridge, addr); err != nil {
		// Clean up the bridge if IP assignment fails
		if delErr := netlink.LinkDel(bridge); delErr != nil {
			return fmt.Errorf("failed to delete bridge after IP assignment failure: %v", err)
		}
		return fmt.Errorf("failed to assign IP address to bridge: %v", err)
	}
	fmt.Printf("Bridge %s created and assigned IP Address %s\n", bridgeName, ipAddress)

	// Bring up the bridge
	if err := netlink.LinkSetUp(bridge); err != nil {
		return fmt.Errorf("failed to up the bridge: %v", err)
	}

	// Setup NAT rule for the bridge
	ipt, err := iptables.New()
	if err != nil {
		return fmt.Errorf("failed to initialize iptables: %v", err)
	}
	if err := ipt.AppendUnique("nat", "POSTROUTING", "-o", bridgeName, "-j", "MASQUERADE"); err != nil {
		return fmt.Errorf("failed to setup the NAT Rule to the bridge: %v", err)
	}

	// Enable IP forwarding
	if err := enableIPForwarding(); err != nil {
		return fmt.Errorf("failed to enable IP forwarding: %v", err)
	}

	hostInterface := "enp3s0"

	// Add a NAT rule for the host's network interface
	hostLink, err := netlink.LinkByName(hostInterface)
	if err != nil {
		return fmt.Errorf("failed to get host interface: %v", err)
	}
	if err := ipt.AppendUnique("nat", "POSTROUTING", "-o", hostLink.Attrs().Name, "-j", "MASQUERADE"); err != nil {
		return fmt.Errorf("failed to add NAT rule for host's network interface: %v", err)
	}

	return nil
}

func enableIPForwarding() error {
	file, err := os.OpenFile("/proc/sys/net/ipv4/ip_forward", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString("1"); err != nil {
		return err
	}

	return nil
}

func SetupBridgeNetwork() (bridge string, userID string, bridge_ip_address string, bridge_gateway_ip string, err error) {
	fmt.Println("Setting up bridge")

	bridgeName, userID, bridge_ip_address, bridge_gateway_ip := generateValue()

	fmt.Println("Bridge Name:", bridgeName)
	fmt.Println("User ID:", userID)
	fmt.Println("bridge_ip_address:", bridge_ip_address)
	fmt.Println("bridge_gateway_ip:", bridge_gateway_ip)

	if err = createBridge(bridgeName, bridge_ip_address); err != nil {
		fmt.Println("Error creating bridge:", err)
		return
	}

	return bridgeName, userID, bridge_ip_address, bridge_gateway_ip, nil
}
