package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/techcampman/twitter-d-server/db"
	"github.com/techcampman/twitter-d-server/env"
	"github.com/techcampman/twitter-d-server/logger"
)

// EncodeImage encodes a image
func EncodeImage(ct string, iBuf *bytes.Buffer, img image.Image) (err error) {
	switch ct {
	case db.Jpeg:
		err = jpeg.Encode(iBuf, img, &jpeg.Options{Quality: 100})
	case db.Png:
		err = png.Encode(iBuf, img)
	default:
		err = fmt.Errorf("Bad content-type: %v", ct)
	}

	if err != nil {
		logger.Error(err)
	}
	return
}

// UploadImage uploads a image to a storage
func UploadImage(basePath string, data []byte, ct string, tag ...string) (path string, err error) {

	// upload with extension
	fm, err := db.GetFormat(ct)
	if err != nil {
		logger.Error(err)
		return
	}

	if fm == db.Formats[db.Jpeg] {
		fm = "jpg"
	}

	if len(tag) > 0 {
		path = basePath + "." + fm + tag[0]
	} else {
		path = basePath + "." + fm
	}
	err = env.GetStorage().Put(path, data, ct)
	if err != nil {
		logger.Error(err)
	}
	return
}
