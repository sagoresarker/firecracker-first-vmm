package networking

import (
	"fmt"
)

var (
	tapName1, tapName2       string
	bridge_ip_address_global = ""
	bridge_gateway_ip_global = ""
)

func SetUpNetwork() (string, string, string, string, string, string) {

	fmt.Println("Setup full networking")
	bridgeName, userID, bridge_ip_address, bridge_gateway_ip, err := SetupBridgeNetwork()
	bridge_ip_address_global = bridge_ip_address
	bridge_gateway_ip_global = bridge_gateway_ip

	if err != nil {
		fmt.Println("(From SetUp network) - Error setting up bridge network:", err)
	}
	fmt.Println("(From SetUp network) - Bridge IP Address:", bridge_ip_address)
	tapName1, tapName2, err = SetupTapNetwork(bridgeName)

	if err != nil {
		fmt.Println("(From SetUp network) - Error setting up tap network:", err)
	}

	fmt.Println("Bridge Name:", bridgeName)
	fmt.Println("Tap Names:", tapName1, tapName2)

	return userID, bridgeName, tapName1, tapName2, bridge_ip_address, bridge_gateway_ip
}

func GetTapNames() (string, string) {

	fmt.Println("The Tap Names are (from Get Method):", tapName1, tapName2)

	return tapName1, tapName2
}

func GetBridgeIPAddress() (string, string) {
	fmt.Println("The Bridge IP Address is (from Get Method):", bridge_ip_address_global)
	return bridge_ip_address_global, bridge_gateway_ip_global
}
