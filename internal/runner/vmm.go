package runner

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/sagoresarker/firecracker-first-vmm/internal/database"
	"github.com/sagoresarker/firecracker-first-vmm/internal/networking"
	"github.com/sagoresarker/firecracker-first-vmm/types"
	"github.com/sirupsen/logrus"
)

func LaunchFirstInstance(user_id, bridge_name, tapName1, tapName2 string) {

	bridge_ip_address, bridge_gateway_ip := networking.GetBridgeIPAddress()
	bridge_ip_without_mask, _, err := net.ParseCIDR(bridge_ip_address)
	if err != nil {
		fmt.Println("Error parsing bridge IP address:", err)
		return
	}

	vm1_eth0_ip, vm2_eth0_ip, err := networking.GetVMIPs(bridge_ip_without_mask.String())
	if err != nil {
		fmt.Println("Error getting VM1 IP:", err)
		return
	}

	mac_address1, mac_address2 := networking.GetMACAddress()

	socket_dir := "Socketfiles/" + user_id + "/tmp/"
	socket_path := socket_dir + "firecracker1.sock"

	// Check if the directory exists
	if _, err := os.Stat(socket_dir); os.IsNotExist(err) {
		// Create the directory with necessary permissions
		err := os.MkdirAll(socket_dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	vmmDetails := types.VMMDetails{
		UserID:          user_id,
		BridgeName:      bridge_name,
		TapName1:        tapName1,
		TapName2:        tapName2,
		VM1Eth0IP:       vm1_eth0_ip,
		VM2Eth0IP:       vm2_eth0_ip,
		MacAddress1:     mac_address1,
		MacAddress2:     mac_address2,
		BridgeIPAddress: bridge_ip_without_mask.String(),
		BridgeGatewayIP: bridge_gateway_ip,
	}

	database.SaveVMsDetails(vmmDetails)

	launcherDetails := types.LauncherDetails{
		TapName:         tapName1,
		VMIP:            vm1_eth0_ip,
		MacAddress:      mac_address1,
		BridgeIP:        bridge_ip_without_mask.String(),
		BridgeGatewayIP: bridge_gateway_ip,
		SocketPath:      socket_path,
	}

	launchVM(launcherDetails)

}

func launchVM(launcherDetails types.LauncherDetails) {

	fmt.Println("Launching VM with tap:", launcherDetails.TapName)

	vm_eth0_ip_ipv4 := net.ParseIP(launcherDetails.VMIP)
	if vm_eth0_ip_ipv4 == nil {
		fmt.Println("Error parsing VM IP address")
		return
	}

	bridge_gateway_ip_ipv4 := net.ParseIP(launcherDetails.BridgeGatewayIP)
	fmt.Printf("Bridge Gateway IP: %s and Type %s\n", bridge_gateway_ip_ipv4, reflect.TypeOf(bridge_gateway_ip_ipv4).String())
	if bridge_gateway_ip_ipv4 == nil {
		fmt.Println("Error parsing bridge gateway IP address")
		return
	}

	cfg := firecracker.Config{
		SocketPath:      launcherDetails.SocketPath,
		LogFifo:         launcherDetails.SocketPath + ".log",
		MetricsFifo:     launcherDetails.SocketPath + "-metrics",
		LogLevel:        "Debug",
		KernelImagePath: "files/vmlinux",
		KernelArgs:      "ro console=ttyS0 reboot=k panic=1 pci=off",
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(2),
			MemSizeMib: firecracker.Int64(512),
			Smt:        firecracker.Bool(false),
		},
		Drives: []models.Drive{
			{
				DriveID:      firecracker.String("1"),
				IsRootDevice: firecracker.Bool(true),
				IsReadOnly:   firecracker.Bool(false),
				PathOnHost:   firecracker.String("files/rootfs.ext4"),
			},
		},
		NetworkInterfaces: []firecracker.NetworkInterface{
			{
				StaticConfiguration: &firecracker.StaticNetworkConfiguration{
					MacAddress:  launcherDetails.MacAddress,
					HostDevName: launcherDetails.TapName,
					IPConfiguration: &firecracker.IPConfiguration{
						IPAddr: net.IPNet{
							IP:   vm_eth0_ip_ipv4,
							Mask: net.CIDRMask(24, 32),
						},
						Gateway:     net.ParseIP(launcherDetails.BridgeIP),
						IfName:      "eth0",
						Nameservers: []string{"8.8.8.8", "8.8.4.4"},
					},
				},
			},
		},
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	entry := logrus.NewEntry(logger)
	ctx := context.Background()

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithLogger(entry))
	if err != nil {
		fmt.Printf("Failed to create VM: %v\n", err)
		return
	}

	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Printf("Received signal: %s\n", sig)
		vmmCancel()
	}()

	if err := m.Start(vmmCtx); err != nil {
		fmt.Printf("Failed to start VM: %v\n", err)
		return
	}

	if err := m.Wait(vmmCtx); err != nil {
		fmt.Printf("VM exited with error: %v\n", err)
	} else {
		fmt.Println("VM exited successfully")
	}
}
