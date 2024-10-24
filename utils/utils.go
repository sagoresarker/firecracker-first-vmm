package utils

import (
	"math/rand"
	"time"
)

func createUserID() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 2)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	return "user" + string(result)
}

func getUserID() string {
	return createUserID()
}
