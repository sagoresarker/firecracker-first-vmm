package networking

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func createTap(tapName, bridgeName string) error {
	// Check if a TAP interface with the same name already exists
	existingTap, err := netlink.LinkByName(tapName)
	if err == nil {
		// TAP interface with the same name exists, delete it
		if err := netlink.LinkDel(existingTap); err != nil {
			return fmt.Errorf("failed to delete existing tap: %v", err)
		}
	}

	tapLink := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: tapName,
		},
		Mode: netlink.TUNTAP_MODE_TAP,
	}

	if err := netlink.LinkAdd(tapLink); err != nil {
		return fmt.Errorf("failed to create tap: %v", err)
	}

	// Bring up the TAP interface
	if err := netlink.LinkSetUp(tapLink); err != nil {
		// Clean up the partially created TAP interface
		netlink.LinkDel(tapLink)
		return fmt.Errorf("failed to bring up tap: %v", err)
	}

	// Get the bridge link
	bridgeLink, err := netlink.LinkByName(bridgeName)
	if err != nil {
		// Clean up the created TAP interface
		netlink.LinkDel(tapLink)
		return fmt.Errorf("failed to get bridge link: %v", err)
	}

	// Assign the TAP interface to the bridge
	if err := netlink.LinkSetMaster(tapLink, bridgeLink.(*netlink.Bridge)); err != nil {
		// Clean up the created TAP interface
		netlink.LinkDel(tapLink)
		return fmt.Errorf("failed to assign tap to bridge: %v", err)
	}

	fmt.Printf("Tap %s assigned to Bridge %s\n", tapName, bridgeName)

	netlink.LinkDel(tapLink)

	return nil
}

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
