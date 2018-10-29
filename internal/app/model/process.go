package model

import (
	"github.com/gomodule/redigo/redis"
)

func process(msg redis.Message) {
	// fmt.Printf("%s: message: %s\n", msg.Channel, msg.Data)
}
