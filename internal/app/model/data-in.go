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
			err := listenUntilErr(addr, topic)

			log.Printf("%s", err)
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
			err := process(v)
			if err != nil {
				log.Printf("%s", err)
			}
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return err
		}
	}
}
