package networking

import (
	"log"

	"github.com/sagoresarker/firecracker-first-vmm/types"
)

var (
	tapName1, tapName2       string
	bridge_ip_address_global = ""
	bridge_gateway_ip_global = ""
)

// SetUpNetwork sets up the network (bridge, tap, ip address, gateway ip address)
func SetUpNetwork() (types.NetworkDetails, error) {
	log.Println("Setup full networking")
	bridgeDetails, err := SetupBridgeNetwork()
	bridge_ip_address_global = bridgeDetails.BridgeIPAddress
	bridge_gateway_ip_global = bridgeDetails.BridgeGatewayIP

	if err != nil {
		log.Println("(From SetUp network) - Error setting up bridge network:", err)
	}
	log.Println("(From SetUp network) - Bridge IP Address:", bridgeDetails.BridgeIPAddress)
	tapName1, tapName2, err = SetupTapNetwork(bridgeDetails.BridgeName)

	if err != nil {
		log.Println("(From SetUp network) - Error setting up tap network:", err)
	}

	log.Println("Bridge Name:", bridgeDetails.BridgeName)
	log.Println("Tap Names:", tapName1, tapName2)

	return types.NetworkDetails{
		UserID:          bridgeDetails.UserID,
		BridgeName:      bridgeDetails.BridgeName,
		TapName1:        tapName1,
		TapName2:        tapName2,
		BridgeIPAddress: bridgeDetails.BridgeIPAddress,
		BridgeGatewayIP: bridgeDetails.BridgeGatewayIP,
	}, nil
}

func GetTapNames() (string, string) {
	log.Println("The Tap Names are (from Get Method):", tapName1, tapName2)
	return tapName1, tapName2
}

func GetBridgeIPAddress() (string, string) {
	log.Println("The Bridge IP Address is (from Get Method):", bridge_ip_address_global)
	return bridge_ip_address_global, bridge_gateway_ip_global
}
