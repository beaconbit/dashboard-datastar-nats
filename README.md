# Reef Dashboard

Single-page dashboard with Go backend and Datastar frontend.

## Features

- Responsive layout: TV (4 quarters grid) and mobile (vertical stack)
- Dynamic counters with click interaction
- Automatic screen detection (TV vs mobile)

## Running

```bash
go run main.go
```

Open http://localhost:3001 in a browser.

## TV Layout

- Non-scrollable, fits entire screen
- Four colored quarters (top-left, top-right, bottom-left, bottom-right)
- Click each quarter to increment its counter

## Mobile Layout

- Vertical scrollable layout
- Quarters stack vertically
- Same interactive counters

## Technology

- Backend: Go net/http with HTML templates
- Frontend: Datastar (loaded from CDN) for reactive signals
- CSS: Custom responsive styles with flexbox/grid

## Project Structure

- `main.go` – Go server
- `templates/dashboard.html` – HTML template with embedded CSS/JS
- `go.mod` – Go module definition# dashboard-datastar-nats
