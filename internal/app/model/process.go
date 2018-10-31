package model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gomodule/redigo/redis"
	"gitlab.jiangxingai.com/luyor/tf-fence-backend/log"
)

// fenceMessage is the detection result of tf fence detector
type fenceMessage struct {
	Timestamp float32
	Image     string
	Device    string
	Output    []recognition
}

type recognition struct {
	// Box: ymin, xmin, ymax, xmax
	Box [4]float32
}

func process(msg redis.Message) error {
	var data fenceMessage
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return err
	}

	go func() {
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data.Image))
		img, _, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, rcg := range data.Output {
			box := rcg.Box

			if isOutside(box) {
				height, width := float32(img.Bounds().Max.Y), float32(img.Bounds().Max.X)
				rect := image.Rect(int(box[1]*width), int(box[0]*height), int(box[3]*width), int(box[2]*height))
				personImg := img.(subImager).SubImage(rect)

				fname := saveImg(personImg)
				go func() {
					base := filepath.Base(fname)
					detail := map[string]string{
						"image_url": "/img/" + base,
					}

					err := pushEvent("人员跨越电子围栏", data.Device, nil, detail)
					if err != nil {
						log.Warning(err)
					}
				}()
			}
		}
	}()

	return nil
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

// TODO: save file in database to survive from app container removal
func saveImg(personImg image.Image) string {
	var imageBuf bytes.Buffer
	jpeg.Encode(&imageBuf, personImg, nil)

	f, err := ioutil.TempFile("img", "event-*.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(imageBuf.Bytes())
	return f.Name()
}
