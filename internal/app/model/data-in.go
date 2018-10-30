package model

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"gitlab.jiangxingai.com/luyor/tf-pose-backend/config"
	"gitlab.jiangxingai.com/luyor/tf-pose-backend/log"
)

// Listen subscribes to messages from redis, and handle the message
// Try to reconnect if failed.
func Listen() {
	cfg := config.Config()
	addr, topic := cfg.GetString("data-in-addr"), cfg.GetString("data-in-chan")

	go func() {
		for {
			err := listenUntilErr(addr, topic)

			log.Errorf("Redis connection failed: %s", err)
			time.Sleep(time.Second * 1)
		}
	}()
}

// listen subscribes to messages from redis, and handle the message.
// Returns error if connection failed.
func listenUntilErr(addr, topic string) error {
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return err
	}

	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(topic)

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			go func() {
				err := process(v)
				if err != nil {
					log.Errorf("Data processing failed: %s", err)
				}
			}()
		case redis.Subscription:
			log.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return err
		}
	}
}
