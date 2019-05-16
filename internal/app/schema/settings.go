package schema

const SETTING_CHECKIN_SCHEDULE = "checkin_schedule"

type CheckinTimestamp struct {
    Hour        int64 `json:"hour" bson:"hour"`
    Minute      int64 `json:"minute" bson:"minute"`
    Second      int64 `json:"second" bson:"second"`
}

func (t *CheckinTimestamp) TranslateToSec() int64 {
    return t.Hour * 3600 + t.Minute * 60 + t.Second
}

// Settings is a program db settings
type Settings struct {
	ID         string   `json:"id" bson:"_id,omitempty"`
	Name       string   `json:"name" bson:"name"`
    Starttime  CheckinTimestamp `json:"starttime" bson:"starttime"`
    Endtime    CheckinTimestamp `json:"endtime" bson:"endtime"`
    Interval   int64    `json:"interval" bson:"interval"`
    Duration   int      `json:"duration" bson:"duration"`
}

// SettingsReq is a settings in http req/resp
type SettingsReq struct {
    Starttime  CheckinTimestamp `json:"starttime" bson:"starttime"`
    Endtime    CheckinTimestamp `json:"endtime" bson:"endtime"`
    Interval   int64    `json:"interval" bson:"interval"`
    Duration   int      `json:"duration" bson:"duration"`
}

// SettingsToReq translate db to http
func SettingsToReq(d *Settings) *SettingsReq {
    return &SettingsReq{
        Starttime: d.Starttime,
        Endtime: d.Endtime,
        Interval: d.Interval,
        Duration: d.Duration,
    }
}

// ReqToSettings translate http to db
func ReqToSettings(d *SettingsReq) Settings {
    return Settings{
        Name:       SETTING_CHECKIN_SCHEDULE,
        Starttime:  d.Starttime,
        Endtime:    d.Endtime,
        Interval: d.Interval,
        Duration: d.Duration,
    }
}
