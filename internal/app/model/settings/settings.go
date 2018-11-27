package settings

import (
	"context"
    "time"

	"github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/options"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/timer"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// RegisterSettingsTimer just init timer
var RegisterSettingsTimer = timer.RegisterHandler(AutoCheckinTimer, true)

func collection() *mongo.Collection {
	return model.DB.Collection("settings")
}

// GetSettings get settings
func GetSettings() (*schema.SettingsReq, int64, error) {
    // TODO region as default
	doc := collection().FindOne(context.Background(), map[string]string{"name": schema.SETTING_CHECKIN_SCHEDULE})
	result := &schema.Settings{}
	err := doc.Decode(&result)
	if err != nil {
		return nil, 0, err
	}
    return schema.SettingsToReq(result), result.LastChecktime, nil
}

// UpdateSettings update settings in db: create if not exists
func UpdateSettings(h *schema.SettingsReq) error {
    d := schema.ReqToSettings(h)
    opt := &options.UpdateOptions{}
    opt.SetUpsert(true)
	_, err := collection().UpdateOne(context.Background(), map[string]string{"name": schema.SETTING_CHECKIN_SCHEDULE}, d, opt)

    timer.UpdateTimer()
	return err
}

// UpdateSettingsWithLast update last checking time
func UpdateSettingsWithLast(h *schema.SettingsReq, last int64) error {
    d := schema.ReqToSettings(h)
    d.LastChecktime = last
	_, err := collection().UpdateOne(context.Background(), map[string]string{"name": schema.SETTING_CHECKIN_SCHEDULE}, d)
	return err
}

// AutoCheckinTimer auto checking
func AutoCheckinTimer() (int64, error) {
    s, last, err := GetSettings()
    if err != nil {
        return 0, err
    }
    
    now := time.Now()
    nowSeconds := int64(now.Hour() * 3600 + now.Minute() * 60 + now.Second())
    start := s.Starttime.TranslateToSec()
    end := s.Endtime.TranslateToSec()
    if last == 0 {
        last = start
    }

    var nextTime int64
    if nowSeconds < end && nowSeconds >= start {
        pending := last + s.Interval
        if nowSeconds > pending && pending < end {
            last = int64((nowSeconds - start) / s.Interval) * s.Interval

            // TODO maybe need try more time
	        err = checkin.DefaultCheckiner.Start()
            waitTimer := time.NewTimer(30 * time.Second)
            go func() {
                <- waitTimer.C
                t, err := checkin.DefaultCheckiner.Stop()
                log.Info("AutoCheckTimer: %ld: ", t, err)
            } ()

        }

        if nowSeconds + s.Interval > 86400 {
            last = start
            nextTime = start + 86400 - nowSeconds
        } else {
            nextTime = nowSeconds - last + s.Interval
        }
    } else if nowSeconds < start {
        nextTime = start + nowSeconds
        last = start
    } else {
        nextTime = start + 86400 - nowSeconds
        last = start
    }

    err = UpdateSettingsWithLast(s, last)
    return nextTime, err
}
