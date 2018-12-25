package remote

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
    "bytes"
    "encoding/json"
    "time"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/device"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// Capture an image from camera server
func Capture(deviceName string) (string, error) {
	cfg := config.Config()
	cameraAddr := cfg.GetString("camera-addr")
	requestURL := fmt.Sprintf("%s?device=%s", cameraAddr, deviceName)
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
        var maxRetry = 600; // 10mins
        for maxRetry != 0 {
            time.Sleep(1000 * time.Millisecond)
	        resp, err = http.Get(requestURL)
            if resp.StatusCode == http.StatusOK {
                break
            } else if resp.StatusCode == http.StatusNotFound {
                AddDevices()
            }
            maxRetry = maxRetry - 1
        }
        if maxRetry == 0 {
            return "", fmt.Errorf("check camera status")
        }
    }

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded, nil
}

// AddDevices add post all device in db to Capture services
func AddDevices() {
	cfg := config.Config()
	cameraAddr := cfg.GetString("camera-addr")
	requestURL := fmt.Sprintf("%s", cameraAddr)

	cameras, err := device.GetCameras()
    if err != nil {
        log.Info("Device Post Failed: ", err)
        return
    }

    for _, c := range cameras {
        pc := &schema.CaptureCamera{
            Device: c.DeviceName,
            Rtmp: c.Rtmp,
        }
        jsonValue, _ := json.Marshal(pc)
	    resp, err := http.Post(requestURL,"application/json",  bytes.NewBuffer(jsonValue))
        log.Info("Device Post: ", pc, " result: ", resp, " error: ", err)
    }
}
