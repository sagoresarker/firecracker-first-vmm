package types

// type logger struct {
// 	LogPath       string `json:"log_path"`
// 	Level         string `json:"level"`
// 	ShowLevel     bool   `json:"show_level"`
// 	ShowLogOrigin bool   `json:"show_log_origin"`
// }

// type machineConfig struct {
// 	MemSizeMiB uint `json:"mem_size_mib"`
// 	VCPUCount  uint `json:"vcpu_count"`
// }

// type bootSource struct {
// 	KernelImagePath string `json:"kernel_image_path"`
// 	BootArgs        string `json:"boot_args"`
// }

// type drive struct {
// 	DriveID      string `json:"drive_id"`
// 	PathOnHost   string `json:"path_on_host"`
// 	IsRootDevice bool   `json:"is_root_device"`
// 	IsReadOnly   bool   `json:"is_read_only"`
// }

// type networkInterface struct {
// 	IfaceID     string `json:"iface_id"`
// 	GuestMAC    string `json:"guest_mac"`
// 	HostDevName string `json:"host_dev_name"`
// }

// type action struct {
// 	ActionType string `json:"action_type"`
// }

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
