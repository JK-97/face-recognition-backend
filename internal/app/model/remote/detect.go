package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// Detect calls remote ai service to recognize faces in image
func Detect(img string) ([]schema.Recognition, error) {
	v := schema.DetectFaceReq{Image: img}

	jsonValue, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	cfg := config.Config()
	url := cfg.GetString("face-ai-addr") + "/detect"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("face ai detect service response error: %v, %s", resp.StatusCode, b)
	}

	var result schema.DetectFaceResp
	json.Unmarshal(b, &result)

	return result.Data.Recognitions, nil
}

// CheckDetectAI pings remote ai to check for availbility
func CheckDetectAI() bool {
	cfg := config.Config()
	url := cfg.GetString("face-ai-addr") + "/detect"
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
