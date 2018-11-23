package remote

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
)

// Capture an image from camera server
func Capture(deviceName string) (string, error) {
	cfg := config.Config()
	cameraAddr := cfg.GetString("camera-addr")
	requestURL := fmt.Sprintf("%s?device=%s", cameraAddr, deviceName)
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded, nil
}
