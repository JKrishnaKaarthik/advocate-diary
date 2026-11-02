# Advocate Diary

A small diary each lawyer keeps for their own practice. It holds that lawyer's client list, case details, uploaded papers, hearing dates, and reminders before court.

Lawyers do not see each other's records. There is no clerk and no shared chamber in this version.

## Run

PostgreSQL and S3 run locally through Docker. The Go API is the server. The React app proxies `/api` to it.

```bash
docker compose up -d
cp .env.example .env
go run ./cmd/server
```

In another terminal:

```bash
cd web
npm install
npm run dev
```

Open http://localhost:5173. Each sign-up is a separate lawyer. With `SMTP_HOST` empty, court reminders are printed in the API log and still land in the inbox.

The design is in [docs/hld.md](docs/hld.md) and [docs/lld.md](docs/lld.md).

| Document | What it answers |
| --- | --- |
| [docs/hld.md](docs/hld.md) | Who uses it, what the system is made of, and how the main flows move through it |
| [docs/lld.md](docs/lld.md) | Tables, states, APIs, the reminder worker, and the upload path |

## v1 boundary

Many lawyers, one login each, data separated by the signed-in lawyer. Clients do not sign in. The app does not file into a court system, send WhatsApp, or bill fees.

## Stack

React for the screens, a Go API, PostgreSQL, and S3 for documents. The reminder loop runs inside the Go process.
