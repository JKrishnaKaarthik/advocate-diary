# Advocate Diary

A small practice diary for one advocate chamber. It keeps the client list, case details, uploaded papers, hearing dates, and reminders before court.

This repository currently holds the design, not the application.

| Document | What it answers |
| --- | --- |
| [docs/hld.md](docs/hld.md) | Who uses it, what the system is made of, and how the main flows move through it |
| [docs/lld.md](docs/lld.md) | Tables, states, APIs, the reminder worker, and the upload path |

## v1 boundary

One chamber, one login shared by the advocate and an optional clerk. Clients do not sign in. The app does not file into a court system, send WhatsApp, or bill fees.

## Suggested build

One web app, PostgreSQL, file storage, and a reminder loop in the same deployment. Details and the reason for staying a monolith are in the high-level design.

Build order once implementation starts:

1. Sign-in and chamber scope
2. Clients and cases
3. Diary home
4. Hearings and reminders
5. Document upload
