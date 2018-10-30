package model

import (
	"encoding/base64"
	"encoding/json"
	_ "image/jpeg" // jpeg is imported for its initialization side-effect,
	// which allows image.Decode to understand JPEG formatted images.
	"io"
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
	Output    string
}

func process(msg redis.Message) error {
	var data fenceMessage
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return err
	}

	go func() {
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data.Image))
		// img, format, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: save file in database to survive from app container removal
		f, err := ioutil.TempFile("img", "event-*.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		io.Copy(f, reader)

		go func() {
			base := filepath.Base(f.Name())
			detail := map[string]string{
				"image_url": "/img/" + base,
			}

			err := pushEvent("test", "device1", nil, detail)
			if err != nil {
				log.Warning(err)
			}
		}()
	}()

	return nil
}
