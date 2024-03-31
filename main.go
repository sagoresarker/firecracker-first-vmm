package main

import (
	"fmt"

	vm "github.com/sagoresarker/firecracker-first-vmm/container"
	"github.com/sagoresarker/firecracker-first-vmm/database"
	"github.com/sagoresarker/firecracker-first-vmm/networking"
)

func main() {
	fmt.Println("Hello Poridhians!")
	database.InitMongoDB()

	user_id, bridge_name, _, _, _, _ := networking.SetUpNetwork()

	fmt.Println("User ID for the vm1 is -----------:", user_id)

	// bridgeDetails, err := database.GetBridgeDetails()
	// if err != nil {
	// 	fmt.Println("Error getting bridge details:", err)
	// 	return
	// }

	// fmt.Println("Bridge Details:", bridgeDetails)

	tapName1, tapName2 := networking.GetTapNames()

	vm.LaunchFirstInstance(user_id, bridge_name, tapName1, tapName2)
}
