package ssehub

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID      string
	Payload string
}

type Client struct {
	id         string
	room       string
	ctx        context.Context
	cancel     context.CancelFunc
	sendCh     chan Event
	connected  time.Time
	lastActive time.Time
}

func (c *Client) String() string {
	return fmt.Sprintf("client[%s] room=%s", c.id, c.room)
}

func (c *Client) ContextDone() <-chan struct{} {
	return c.ctx.Done()
}

type Hub struct {
	rooms map[string]map[string]*Client
	mu    sync.RWMutex
}
