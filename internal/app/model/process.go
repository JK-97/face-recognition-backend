package model

import (
	"encoding/json"
	_ "image/jpeg" // jpeg is imported for its initialization side-effect,
	// which allows image.Decode to understand JPEG formatted images.

	"github.com/gomodule/redigo/redis"
	"gitlab.jiangxingai.com/luyor/tf-pose-backend/log"
)

// PoseMessage is the detection result of tf pose detector
type PoseMessage struct {
	Timestamp float32
	Image     string
	Output    string
}

func process(msg redis.Message) error {
	var data PoseMessage
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return err
	}

	// reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data.Image))
	// config, format, err := image.DecodeConfig(reader)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go func() {
		err := pushEvent("test", "device1", nil, nil)
		if err != nil {
			log.Warning(err)
		}
	}()

	return nil
}
