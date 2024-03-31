package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func SaveVMsDetails(userID, bridgeName string, tapName1 string, tapName2 string, vm1_eth0_ip, vm2_eth0_ip, mac_address1, mac_address2, Bridge_ipAddress string) error {
	if mongoClient == nil {
		log.Fatal("MongoDB client not initialized.")
		return nil
	}

	collection := mongoClient.Database("firecrackerdb").Collection("vm-info")
	document := bson.D{
		{Key: "userID", Value: userID},
		{Key: "bridgeName", Value: bridgeName},
		{Key: "tapName1", Value: tapName1},
		{Key: "tapName2", Value: tapName2},
		{Key: "vm1_eth0_ip", Value: vm1_eth0_ip},
		{Key: "vm2_eth0_ip", Value: vm2_eth0_ip},
		{Key: "mac_address1", Value: mac_address1},
		{Key: "mac_address2", Value: mac_address2},
		{Key: "Bridge_ipAddress", Value: Bridge_ipAddress},
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
