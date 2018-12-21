package timer

import (
	"fmt"
	"io/ioutil"
	"net/http"
    "encoding/json"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/device"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
)

// InitReload run check cameras timer
func InitReload() {
    cfg := config.Config()
    appid := cfg.GetString("appid")
    if appid != "" {
        log.Info("HAS APPID")
        RegisterHandler(autoReloadTimer)
    } else {
        log.Info("NO APPID")
    }
}

func autoReloadTimer(init bool) (int64, error) {
    log.Info("start autoReloadTimer")
    if checkin.DefaultCheckiner.Status() == schema.CHECKING {
        log.Info("checking is running, won't autoReloadTimer")
        return util.NowMilli() + 300000, nil
    } else if init {
        log.Info("reload cameras: true")
        return util.NowMilli() + 5000, reloadCameras()
    } else {
        log.Info("reload cameras: false")
        return util.NowMilli() + 300000, reloadCameras()
    }
}

func reloadCameras() error {
    cfg := config.Config()
    appid := cfg.GetString("appid")
    gatewayAddr := cfg.GetString("gateway-addr")
    requestURL := fmt.Sprintf("%s/internalapi/v1/%s/device/all", gatewayAddr, appid)
    resp, err := http.Get(requestURL)
    if err == nil && resp.StatusCode == http.StatusOK {
	    b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return err
        }

        var p schema.CamerasResp
        err = json.Unmarshal(b, &p)
        if err != nil {
            return err
        }

        device.RemoveCameras()
        for _, c := range p.Data {
            device.AddCamera(&schema.Camera{
                Name: c.Name,
                Rtmp: c.StreamAddr,
                DeviceName: c.ID,
            })
        }
        remote.AddDevices()
    }
    return err
}
