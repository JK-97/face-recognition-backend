package model

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Listen subscribes to messages from redis, and handle the message
// Try to reconnect if failed.
func Listen(addr, topic string) {
	go func() {
		for {
			err := listenUntilError(addr, topic)

			log.Printf("%s", err)
			time.Sleep(time.Second * 1)
		}
	}()
}

// listen subscribes to messages from redis, and handle the message.
// Returns error if connection failed.
func listenUntilError(addr, topic string) error {
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return err
	}

	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(topic)

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			process(v)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return err
		}
	}
}
