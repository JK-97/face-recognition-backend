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

// Record calls remote ai service to record a person's face
func Record(personID string, images []string) error {
	v := schema.RecordReq{
		ID:     personID,
		Images: images,
	}
	jsonValue, err := json.Marshal(v)
	if err != nil {
		return err
	}

	cfg := config.Config()
	url := cfg.GetString("face-ai-addr") + "/record"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("face ai detect service response error: %v, %s", resp.StatusCode, b)
	}
	return nil
}
