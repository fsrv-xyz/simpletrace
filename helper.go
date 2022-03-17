package simpletrace

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func randomID(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println(err)
		return ""
	}
	return hex.EncodeToString(bytes)
}
