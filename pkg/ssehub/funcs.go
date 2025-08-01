package ssehub

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	heartbeatInterval = 25 * time.Second
)

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
	if ev.ID != "" {
		fmt.Fprintf(w, "id: %s\n", ev.ID)
	}

	// split data by newline to follow SSE spec
	for _, line := range splitLines(ev.Payload) {
		fmt.Fprintf(w, "data: %s\n", line)
	}
	fmt.Fprint(w, "\n")
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
