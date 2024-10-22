package main

import (
	"fmt"

	"github.com/sagoresarker/firecracker-first-vmm/internal/database"
	"github.com/sagoresarker/firecracker-first-vmm/internal/networking"
	vm "github.com/sagoresarker/firecracker-first-vmm/internal/runner"
)

func main() {
	fmt.Println("Hello Poridhians!")
	database.InitMongoDB()

	networkDetails, err := networking.SetUpNetwork()
	if err != nil {
		fmt.Println("Error setting up network:", err)
		return
	}

	fmt.Println("User ID for the vm1 is -----------:", networkDetails.UserID)

	// bridgeDetails, err := database.GetBridgeDetails()
	// if err != nil {
	// 	fmt.Println("Error getting bridge details:", err)
	// 	return
	// }

	// fmt.Println("Bridge Details:", bridgeDetails)

	tapName1, tapName2 := networking.GetTapNames()

	vm.LaunchFirstInstance(networkDetails.UserID, networkDetails.BridgeName, tapName1, tapName2)
}
