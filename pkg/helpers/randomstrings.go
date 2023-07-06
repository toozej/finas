package helpers

import (
	"crypto/rand"
	"math/big"

	log "github.com/sirupsen/logrus"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	max := big.NewInt(int64(len(charset)))
	result := make([]byte, length)

	for i := range result {
		seededRand, err := rand.Int(rand.Reader, max)
		if err != nil {
			log.Error("Unable to generate random string for DockerCmd container name")
		}
		result[i] = charset[seededRand.Int64()]
	}

	return string(result)
}
