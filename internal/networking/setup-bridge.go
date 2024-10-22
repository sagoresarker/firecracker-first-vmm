package networking

import (
	"fmt"
	"os/exec"

	"github.com/sagoresarker/firecracker-first-vmm/types"
	"github.com/sagoresarker/firecracker-first-vmm/utils"
)

// createBridge creates a bridge with the given name and IP address
func createBridge(bridgeName string, ipAddress string) error {

	cmd := exec.Command("sudo", "ip", "link", "add", "name", bridgeName, "type", "bridge")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create bridge: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "addr", "add", ipAddress, "dev", bridgeName)
	if err := cmd.Run(); err != nil {
		// If assigning IP address fails, we need to delete the bridge
		cmd := exec.Command("sudo", "ip", "link", "delete", bridgeName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to delete bridge after IP assignment failure: %v", err)
		}
		return fmt.Errorf("failed to assign IP address to bridge: %v", err)
	}

	fmt.Printf("Bridge %s created and assigned IP Address %s\n", bridgeName, ipAddress)

	cmd = exec.Command("sudo", "ip", "link", "set", "dev", bridgeName, "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to up the bridge: %v", err)
	}

	cmd = exec.Command("sudo", "iptables", "-t", "nat", "-A", "POSTROUTING", "-o", bridgeName, "-j", "MASQUERADE")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to setup the NAT Rule to the bridge: %v", err)
	}

	// Enable IP forwarding
	cmd = exec.Command("sudo", "sysctl", "-w", "net.ipv4.ip_forward=1")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable IP forwarding: %v", err)
	}

	// Add a NAT rule for the host's network interface
	cmd = exec.Command("sudo", "iptables", "--table", "nat", "--append", "POSTROUTING", "--out-interface", "enp3s0", "-j", "MASQUERADE")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add NAT rule for host's network interface: %v", err)
	}

	return nil
}

// SetupBridgeNetwork sets up the bridge network
func SetupBridgeNetwork() (types.BridgeDetails, error) {
	fmt.Println("Setting up bridge")

	bridgeDetails, err := utils.GenerateValue()
	if err != nil {
		return types.BridgeDetails{}, err
	}

	fmt.Println("Bridge Name:", bridgeDetails.BridgeName)
	fmt.Println("User ID:", bridgeDetails.UserID)
	fmt.Println("Bridge IP Address:", bridgeDetails.BridgeIPAddress)
	fmt.Println("Bridge Gateway IP:", bridgeDetails.BridgeGatewayIP)

	if err := createBridge(bridgeDetails.BridgeName, bridgeDetails.BridgeIPAddress); err != nil {
		fmt.Println("Error creating bridge:", err)
		return types.BridgeDetails{}, err
	}

	return bridgeDetails, nil
}
