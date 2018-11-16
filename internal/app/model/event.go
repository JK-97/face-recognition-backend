package model

import (
	"github.com/google/uuid"
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
)

// Event is a dashboard event:
// "event_id": uuid,
// "title": "事件标题",
// "labels": {"key": "value"},
// "created_time": int(UTC),
// "device": "设备名称",
// "related_app": "关联app",
// "detail": {"image": "base64image", "image_url": "video": "", "audio": "", "audio_url": ""}
type Event struct {
	EventID     string            `json:"event_id"`
	Title       string            `json:"title"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreatedTime int64             `json:"created_time"`
	Device      string            `json:"device"`
	RelatedApp  string            `json:"related_app"`
	Detail      map[string]string `json:"detail,omitempty"`
}

// pushEvent pushes a dashboard event to redis MQ.
func pushEvent(title, device string, labels, detail map[string]string) error {
	cfg := config.Config()
	relatedApp := cfg.GetString("app-name")

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	createdTime := util.NowMilli()

	e := Event{id.String(), title, labels, createdTime, device, relatedApp, detail}

	jsonfied, err := json.Marshal(e)
	if err != nil {
		return err
	}

	addr, topic := cfg.GetString("event-out-addr"), cfg.GetString("event-out-chan")
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return err
	}
	conn.Do("PUBLISH", topic, jsonfied)

	return nil
}
