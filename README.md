# Real-time Room Relay Service
-------

1. Creating named rooms (UUID + link).
2. Clients to subscribe to a room via Server-Sent Events (SSE).
3. Any HTTP request sent to the room's relay endpoint to be broadcast as an event to all connected clients.

Primary use case: external systems or users post events into a room and subscribers receive them in real time.

## Features

* `POST /rooms` : create a room by name; returns `room_id` and link.
* `GET /rooms/{room_id}/events` : SSE stream for clients to receive events.
* `POST /rooms/{room_id}/relay` : send arbitrary HTTP payload into the room; broadcast to subscribers.
* In-memory room and client management.
* JSON event schema with metadata.

## Architecture

* **Hub**: in-memory registry of rooms.
* **Room**: holds active SSE clients and broadcasts events.
* **Client**: SSE subscriber with buffered send channel and disconnect handling.
* **Event**: normalized representation of incoming HTTP request including method, headers, body, source IP, timestamp, and UUID.

## References

- https://html.spec.whatwg.org/multipage/server-sent-events.html
