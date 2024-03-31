package networking

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	macAddress1 string
	macAddress2 string
	mutex       sync.Mutex
)

func generateMACAddress() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 6)
	for i := range bytes {
		bytes[i] = byte(rand.Intn(256))
	}
	bytes[0] &= 0xfc // Set the first two bits to 0 to ensure the MAC address is unique
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])
}

func GetMACAddress() (string, string) {
	mutex.Lock()
	defer mutex.Unlock()
	if macAddress1 == "" || macAddress2 == "" {
		macAddress1 = generateMACAddress()
		macAddress2 = generateMACAddress()
	}
	fmt.Println("MAC Address 1:", macAddress1)
	fmt.Println("MAC Address 2:", macAddress2)

	return macAddress1, macAddress2
}
