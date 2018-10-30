package model

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"gitlab.jiangxingai.com/luyor/tf-pose-backend/config"
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
	EventID     []byte            `json:"event_id"`
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

	id, err := uuid.NewV1()
	if err != nil {
		return err
	}
	createdTime := time.Now().Unix()

	e := Event{id.Bytes(), title, labels, createdTime, device, relatedApp, detail}

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
