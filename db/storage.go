package db

import "fmt"

// Content-Types for Storage
const (
	Jpeg = "image/jpeg"
	Png  = "image/png"
)

// Formats available Content-Type and format strings map
var Formats = map[string]string{
	Jpeg: "jpeg",
	Png:  "png",
}

// Storage interface for any object storage
type Storage interface {
	Put(path string, data []byte, contentType string) error
	Del(path string) error
}

// GetFormat returns format string by specified Content-Type
func GetFormat(ct string) (f string, err error) {
	f, ok := Formats[ct]
	if !ok {
		err = fmt.Errorf("Unknown Content-Type: %s", ct)
	}
	return
}

// GetContentType returns format string by specified Content-Type
func GetContentType(fm string) (ct string, err error) {
	for c, f := range Formats {
		if f == fm {
			ct = c
			return
		}
	}
	err = fmt.Errorf("Unknown Content-Type: %s", fm)
	return
}
