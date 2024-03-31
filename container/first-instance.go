package container

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
	"github.com/sagoresarker/firecracker-first-vmm/database"
	"github.com/sagoresarker/firecracker-first-vmm/networking"
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

	socket_path := user_id + "/tmp/firecracker1.sock"

	database.SaveVMsDetails(user_id, bridge_name, tapName1, tapName2, vm1_eth0_ip, vm2_eth0_ip, mac_address1, mac_address2, bridge_ip_without_mask.String(), bridge_gateway_ip)

	launchVM(tapName1, vm1_eth0_ip, mac_address1, bridge_ip_without_mask.String(), bridge_gateway_ip, socket_path)

}

func launchVM(tapName, vmIP, mac_address, bridgeIP, bridgeGatewayIP, socketPath string) {

	fmt.Println("Launching VM with tap:", tapName)

	vm_eth0_ip_ipv4 := net.ParseIP(vmIP)
	if vm_eth0_ip_ipv4 == nil {
		fmt.Println("Error parsing VM IP address")
		return
	}

	bridge_gateway_ip_ipv4 := net.ParseIP(bridgeGatewayIP)
	fmt.Printf("Bridge Gateway IP: %s and Type %s\n", bridge_gateway_ip_ipv4, reflect.TypeOf(bridge_gateway_ip_ipv4).String())
	if bridge_gateway_ip_ipv4 == nil {
		fmt.Println("Error parsing bridge gateway IP address")
		return
	}

	cfg := firecracker.Config{
		SocketPath:      socketPath,
		LogFifo:         socketPath + ".log",
		MetricsFifo:     socketPath + "-metrics",
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
				PathOnHost:   firecracker.String("files/build/rootfs.ext4"),
			},
		},
		NetworkInterfaces: []firecracker.NetworkInterface{
			{
				StaticConfiguration: &firecracker.StaticNetworkConfiguration{
					MacAddress:  mac_address,
					HostDevName: tapName,
					IPConfiguration: &firecracker.IPConfiguration{
						IPAddr: net.IPNet{
							IP:   vm_eth0_ip_ipv4,
							Mask: net.CIDRMask(24, 32),
						},
						Gateway:     net.ParseIP(bridgeIP),
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
