package remote

import (
	"time"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// TestRecord calls remote record with 20 captures
func TestRecord() {
	images := make([]string, 0, 10)
	for i := 0; i < 20; i++ {
		image, err := Capture()
		if err != nil {
			log.Error(err)
			return
		}
		images = append(images, image)
		time.Sleep(time.Millisecond * 500)
	}
	err := Record("1", images)
	if err != nil {
		log.Error(err)
		return
	}
}

// TestDetect calls capture then detect
func TestDetect() {
	testImage, _ := Capture()
	rcgs, err := Detect(testImage)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(rcgs)
}
