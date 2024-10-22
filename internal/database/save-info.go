package database

import (
	"context"
	"log"
	"time"

	"github.com/sagoresarker/firecracker-first-vmm/types"
	"go.mongodb.org/mongo-driver/bson"
)

// SaveVMsDetails saves the details of the VMs to the MongoDB database
func SaveVMsDetails(vmmDetails types.VMMDetails) error {
	if mongoClient == nil {
		log.Fatal("MongoDB client not initialized.")
		return nil
	}

	collection := mongoClient.Database("firecrackerdb").Collection("vm-info")
	document := bson.D{
		{Key: "userID", Value: vmmDetails.UserID},
		{Key: "bridgeName", Value: vmmDetails.BridgeName},
		{Key: "tapName1", Value: vmmDetails.TapName1},
		{Key: "tapName2", Value: vmmDetails.TapName2},
		{Key: "vm1_eth0_ip", Value: vmmDetails.VM1Eth0IP},
		{Key: "vm2_eth0_ip", Value: vmmDetails.VM2Eth0IP},
		{Key: "mac_address1", Value: vmmDetails.MacAddress1},
		{Key: "mac_address2", Value: vmmDetails.MacAddress2},
		{Key: "Bridge_ipAddress", Value: vmmDetails.BridgeIPAddress},
		{Key: "bridge_gateway_ip", Value: vmmDetails.BridgeGatewayIP},
		{Key: "created_at", Value: time.Now()},
	}

	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Println("Error saving network details to MongoDB:", err)
		return err
	}

	log.Println("Network details saved to MongoDB.")

	return nil
}

func GetVMDetails() ([]bson.M, error) {
	if mongoClient == nil {
		log.Fatal("MongoDB client not initialized.")
		return nil, nil
	}

	collection := mongoClient.Database("firecrackerdb").Collection("vm-info")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error fetching bridge details from MongoDB:", err)
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Println("Error fetching bridge details from MongoDB:", err)
		return nil, err
	}

	return results, nil
}
