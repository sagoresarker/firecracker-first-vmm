package networking

import (
	"fmt"
	"os/exec"
)

// createTap creates a tap with the given name and assigns it to the bridge
func createTap(tapName string, bridgeName string) error {
	cmd := exec.Command("sudo", "ip", "tuntap", "add", "dev", tapName, "mode", "tap")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tap: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", "dev", tapName, "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring up tap: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", "dev", tapName, "master", bridgeName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to assign tap to bridge: %v", err)
	}

	fmt.Printf("Tap %s assigned to Bridge %s\n", tapName, bridgeName)

	return nil
}

// SetupTapNetwork sets up the tap network with the given bridge name
func SetupTapNetwork(bridgeName string) (string, string, error) {
	fmt.Println("Setting up tap")

	tapName1 := bridgeName + "-tap" + "-1"
	tapName2 := bridgeName + "-tap" + "-2"

	if err := createTap(tapName1, bridgeName); err != nil {
		fmt.Println("Error creating tap for VM1:", err)
		return "", "", err
	}
	if err := createTap(tapName2, bridgeName); err != nil {
		fmt.Println("Error creating tap for VM2:", err)
		return "", "", err
	}
	return tapName1, tapName2, nil
}
