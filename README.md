# Advocate Diary

A small diary each lawyer keeps for their own practice. It holds that lawyer's client list, case details, uploaded papers, hearing dates, and reminders before court.

Lawyers do not see each other's records. There is no clerk and no shared chamber in this version.

## Run

PostgreSQL, S3, the Go API, and the built website run through Docker.

```bash
docker compose up -d --build
```

Open http://localhost:8080. Each sign-up is a separate lawyer. With `SMTP_HOST` empty, court reminders are printed in the API log and still land in the inbox.

To run the API and Vite app on the host instead:

```bash
docker compose up -d postgres minio
cp .env.example .env
go run ./cmd/server
```

```bash
cd web
npm install
npm run dev
```

Then open http://localhost:5173. The React app proxies `/api` to the Go server.

## Pull request checks

Every pull request to `main` runs two GitHub Actions jobs: `go-test` against Postgres, and `ui-regression` (headless Selenium against the Docker stack on http://localhost:8080).

The design is in [docs/hld.md](docs/hld.md) and [docs/lld.md](docs/lld.md).

| Document | What it answers |
| --- | --- |
| [docs/hld.md](docs/hld.md) | Who uses it, what the system is made of, and how the main flows move through it |
| [docs/lld.md](docs/lld.md) | Tables, states, APIs, the reminder worker, and the upload path |

## v1 boundary

Many lawyers, one login each, data separated by the signed-in lawyer. Clients do not sign in. The app does not file into a court system, send WhatsApp, or bill fees.

## Stack

React for the screens, a Go API, PostgreSQL, and S3 for documents. The reminder loop runs inside the Go process.
