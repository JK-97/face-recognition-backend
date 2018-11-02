package model

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io/ioutil"
	"strings"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

// TODO: save file in database to survive from app container removal
func saveImg(img *image.Image, suffix string) string {
	var imageBuf bytes.Buffer
	jpeg.Encode(&imageBuf, *img, nil)

	f, err := ioutil.TempFile("img", suffix+"-*.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(imageBuf.Bytes())
	return f.Name()
}

func base64Str2Img(str string) (*image.Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(str))
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return &img, nil
}
