package ssehub

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	heartbeatInterval = 1 * time.Second
)

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

func (c *Client) Wait() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Client) writerLoop(w http.ResponseWriter, flusher http.Flusher) {
	// send initial comment to establish connection
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-c.ctx.Done():
			return
		case ev := <-c.sendCh:
			c.lastActive = time.Now()
			writeEvent(w, ev)
			flusher.Flush()
		}
	}
}

func writeEvent(w http.ResponseWriter, ev Event) {
	var err error
	if ev.Event != "" {
		_, err = fmt.Fprintf(w, "event: %s\n", ev.Event)
	}

	// split data by newline to follow SSE spec
	for _, line := range splitLines(ev.Payload) {
		_, err = fmt.Fprintf(w, "data: %s\n", line)
	}
	_, err = fmt.Fprint(w, "\n")

	if err != nil {
		log.Printf("[SSE] Error writing data: %s", err)
	}
}

func splitLines(s string) []string {
	// naive split; could use strings.Split if no special behavior needed
	var lines []string
	current := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			ev := Event{
				Event:   "heartbeat",
				Payload: fmt.Sprintf("heartbeat %d", time.Now().Unix()),
			}

			select {
			case c.sendCh <- ev:
			default:
				// if even heartbeat can't be enqueued, treat as slow/unresponsive and cancel
				log.Printf("[SSE] client unresponsive, closing %s", c)
				c.cancel()
				return
			}
		}
	}
}
