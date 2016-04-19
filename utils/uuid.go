package utils

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// GetNewUUIDv4 ... get new UUIDv4
func GetNewUUIDv4() string {
	return fmt.Sprintf("%x", uuid.NewV4().Bytes()[:])
}
