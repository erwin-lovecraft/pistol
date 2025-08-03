# SSE Rooms Viewer (Chi + HTMX)

> A lightweight real-time room event viewer using Server-Sent Events. Built with Go (Chi router) and a minimal
> HTMX-enhanced web UI to display incoming events in a sidebar with detail inspection. Ideal for debugging webhooks,
> internal APIs, or any event stream per room.

## Features

* SSE hub supporting multiple rooms and clients.
* `/rooms/{roomID}/events`: SSE endpoint streaming JSON events.
* `/rooms/{roomID}/views`: Web UI that connects to the SSE feed, shows live messages, and displays details in a
  human-friendly layout.
* Styled to feel clean and modern (inspired by Apple Human Interface Guidelines).
* Demo push endpoint to inject sample events: `POST /rooms/{roomID}/push`.

## Quickstart

### 1. Clone / prepare

```sh
git clone https://github.com/erwin-lovecraft/pistol
cd pistol
```

### 2. Directory layout

```
./
  cmd/
    /serverd
      /main.go                  # server implementation
  internal/
    adapters/
    config/
    core/
    web/
      view.html             # HTML template for the UI (required)
```

Create `templates/view.html` with the provided template (renders the sidebar + detail view and subscribes to SSE).

### 3. Build & run

```sh
make run
```

### 4. Open the UI

Visit in browser:

```
http://localhost:8080/rooms/demo/views
```

### 5. Inject a demo event

```sh
curl -X POST http://localhost:8080/rooms/demo/push
```

The event appears live in the sidebar. Click on it to see headers (as a compact key\:value table) and body.

## API Endpoints

* `GET /rooms/{roomID}/events` - SSE stream for room events. Clients subscribe here.
* `GET /rooms/{roomID}/views` - UI page for a room.
* `ANY /rooms/{roomID}/relay` - To send event into the room.

## Template Customization

You can modify `internal/web/template.html` to adapt styling, add filters, or replace the detail panel logic. The server
injects `{{.RoomID}}` into the view context.

## Extension Ideas

* Persist events to a backing store and load history on page load.
* HTMX partials: load message detail via separate endpoint instead of client-side rendering.
* Authentication per room / access control.
* Dark mode toggle with `prefers-color-scheme` media query.
* Real upstream producers forwarding real webhook payloads into rooms.

## Development Tips

* Enable live reload by using a file watcher to restart Go binary or re-parse templates on change.
* Replace the demo push handler with an adapter that accepts external webhook payloads and broadcasts them as
  `SSEMessage`.

## License

MIT
