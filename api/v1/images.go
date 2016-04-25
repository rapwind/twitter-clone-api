package v1

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"time"

	"crypto/rand"
	"encoding/binary"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/db"
	"github.com/techcampman/twitter-d-server/env"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/service"
)

func uploadImage(c *gin.Context) {
	// Check image size
	if c.Request.ContentLength > constant.ImageMaxSize {
		errors.Send(c, errors.ImageSizeTooLarge())
		return
	}

	f, h, err := c.Request.FormFile("data")
	if err != nil {
		errors.Send(c, errors.BadParams("image", "an unimportable file"))
		return
	}
	defer f.Close()
	// Check Content-Type
	ct := h.Header["Content-Type"][0]
	if _, err = db.GetFormat(ct); err != nil {
		errors.Send(c, errors.UnsupportedMediaType())
		return
	}

	// Decode image
	img, iFmt, err := image.Decode(f)
	if err != nil {
		errors.Send(c, errors.UnsupportedMediaType())
		return
	}

	// Set decoded Content-Type
	ct, err = db.GetContentType(iFmt)
	if err != nil {
		errors.Send(c, errors.UnsupportedMediaType())
		return
	}
	iBuf := new(bytes.Buffer)
	if err = service.EncodeImage(ct, iBuf, img); err != nil {
		errors.Send(c, errors.UnsupportedMediaType())
		return
	}
	if iBuf.Len() > constant.ImageMaxSize {
		errors.Send(c, errors.ImageSizeTooLarge())
		return
	}

	// Check type
	typ := c.Request.FormValue("type")
	if !(typ == constant.ImageTypeProfile || typ == constant.ImageTypeTweet || typ == constant.ImageTypeHeader) {
		errors.Send(c, errors.BadParams("type", "invalid image type"))
		return
	}
	bp := basePath(c, typ)

	// Upload to S3
	p, err := service.UploadImage(bp, iBuf.Bytes(), ct)
	if err != nil {
		errors.Send(c, fmt.Errorf("failed to upload a image to strage server"))
		return
	}

	resp := map[string]interface{}{
		"imageUrl": env.GetStorageURL().String() + "/" + p,
	}

	c.JSON(http.StatusCreated, resp)
}

func basePath(c *gin.Context, typ string) string {
	hash := randomHash()
	return "images/" + typ + "/" + hash + time.Now().Format("20060102150405")
}

func randomHash() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}
