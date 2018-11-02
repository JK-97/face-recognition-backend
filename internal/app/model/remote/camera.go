package remote

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
)

// Capture an image from camera server
func Capture() (string, error) {
	cfg := config.Config()
	cameraAddr := cfg.GetString("camera-addr")
	resp, err := http.Get(cameraAddr)
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
