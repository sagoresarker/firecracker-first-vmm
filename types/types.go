package types

type NetworkDetails struct {
	UserID          string
	BridgeName      string
	TapName1        string
	TapName2        string
	BridgeIPAddress string
	BridgeGatewayIP string
}

type BridgeDetails struct {
	BridgeName      string
	UserID          string
	BridgeIPAddress string
	BridgeGatewayIP string
}

type VMMDetails struct {
	UserID          string
	BridgeName      string
	TapName1        string
	TapName2        string
	VM1Eth0IP       string
	VM2Eth0IP       string
	MacAddress1     string
	MacAddress2     string
	BridgeIPAddress string
	BridgeGatewayIP string
}

type LauncherDetails struct {
	TapName         string
	VMIP            string
	MacAddress      string
	BridgeIP        string
	BridgeGatewayIP string
	SocketPath      string
}
