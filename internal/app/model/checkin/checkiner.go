package checkin

import (
	"fmt"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// DefaultCheckiner is the singleton of Checkiner
var DefaultCheckiner = NewCheckiner()

type stopRespType struct {
	time int64
	err  error
}

// Checkiner periodically run checkin.
type Checkiner struct {
	status  schema.CheckinStatus
	startCh chan chan error
	stopCh  chan chan stopRespType
}

// Start starts periodical checkin
func (c *Checkiner) Start() error {
	respCh := make(chan error)
	c.startCh <- respCh
	err := <-respCh
	return err
}

// Stop stops periodical checkin
func (c *Checkiner) Stop() (int64, error) {
	respCh := make(chan stopRespType)
	c.stopCh <- respCh
	resp := <-respCh
	return resp.time, resp.err
}

// Status returns its status
func (c *Checkiner) Status() schema.CheckinStatus {
	return c.status
}

// NewCheckiner creates a new Checkiner
func NewCheckiner() *Checkiner {
	c := Checkiner{}

	c.status = schema.STOPPED
	c.startCh = make(chan chan error)
	c.stopCh = make(chan chan stopRespType)

	go c.serve()
	return &c
}

func (c *Checkiner) serve() {
	for {
		startTime := c.waitStart()
		c.status = schema.CHECKING
		c.detecting(startTime)
		c.status = schema.STOPPED
	}
}

func (c *Checkiner) waitStart() int64 {
	for {
		select {
		case startResp := <-c.startCh:
			startTime := util.NowMilli()
			err := remote.CheckDetectAI()
			startResp <- err
			if err == nil {
				return startTime
			}
		case stopResp := <-c.stopCh:
			stopResp <- stopRespType{0, fmt.Errorf("checkin is not started")}
		}
	}
}

func (c *Checkiner) detecting(startTime int64) {
	for {
		select {
		case startResp := <-c.startCh:
			startResp <- fmt.Errorf("checkin already started")
		case stopResp := <-c.stopCh:
			saveCheckin(seal{
				startTime: startTime,
				endTime:   util.NowMilli(),
			})
			stopResp <- stopRespType{startTime, nil}
			return
		default:
			err := checkin()
			if err != nil {
				log.Error(err)
			}
		}
	}
}
