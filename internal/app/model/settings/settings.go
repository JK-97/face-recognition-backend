package settings

import (
	"context"
    "time"
    "errors"

	"github.com/google/uuid"
	"github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/options"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/timer"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

var registerSettingsTimer = timer.RegisterHandler(AutoCheckinTimer)

func collection() *mongo.Collection {
	return model.DB.Collection("settings")
}

// GetSettings get settings
func GetSettings() (*schema.SettingsReq, error) {
    db := collection()
    if db == nil {
        return nil, errors.New("DB uninit")
    }
	doc := db.FindOne(context.Background(), map[string]string{"name": schema.SETTING_CHECKIN_SCHEDULE})
	result := &schema.Settings{}
	err := doc.Decode(&result)
	if err != nil {
		return nil, err
	}
    return schema.SettingsToReq(result), nil
}

// UpdateSettings update settings in db: create if not exists
func UpdateSettings(h *schema.SettingsReq) error {
    d := schema.ReqToSettings(h)
    _, err := GetSettings()
    if err == mongo.ErrNoDocuments {
        uuid, _ := uuid.NewUUID()
        d.ID = uuid.String()
    } else if err != nil {
        return err
    }
    opt := &options.UpdateOptions{}
    opt.SetUpsert(true)
	_, err = collection().UpdateOne(context.Background(),
        map[string]string{"name": schema.SETTING_CHECKIN_SCHEDULE},
        map[string]schema.Settings{"$set": d},
        opt)

    timer.UpdateTimer(registerSettingsTimer, nextCheckinTime())
	return err
}

func nextCheckinTime() int64 {
    now := time.Now()
    nowSeconds := int64(now.Hour() * 3600 + now.Minute() * 60 + now.Second())

    s, err := GetSettings()
    if err != nil {
        log.Warn("nextCheckingTime: ", err)
        return 0
    }

    start := s.Starttime.TranslateToSec()
    end := s.Endtime.TranslateToSec()

    var nextTime int64
    if nowSeconds <= end && nowSeconds >= start {
        dt := nowSeconds - start
        nextTime = s.Interval - (dt - (dt / s.Interval) * s.Interval)
    } else if nowSeconds < end {
        nextTime = start - nowSeconds
    } else {
        nextTime = 86400 - nowSeconds + start
    }

    if nextTime == 0 {
        nextTime = 1
    }

    return nextTime
}

// AutoCheckinTimer auto checking
func AutoCheckinTimer(init bool) (int64, error) {

    if init {
        return nextCheckinTime(), nil
    }

    // TODO maybe need try more time
    id, err := checkin.DefaultCheckiner.Start()
    if err == nil {
        waitTimer := time.NewTimer(30 * time.Second)
        go func() {
            <- waitTimer.C
            t, err := checkin.DefaultCheckiner.Stop(id)
            log.Info("AutoCheckTimer: %d: ", t, err)
        } ()
    }

    return nextCheckinTime(), err
}
