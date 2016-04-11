package db

import "time"

// Cache is a interface for cache stores
type Cache interface {
	// Get is used to get the value of key
	Get(key string) (reply interface{}, err error)
	// Set key to hold the string value and set key to timeout after a given number of seconds
	Set(key string, value interface{}, expire time.Duration) (err error)
	// Delete is used to remove the specified a key
	Delete(key string) (result int, err error)
}
