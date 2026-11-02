# Advocate Diary — Low Level Design

This is the v1 design behind [hld.md](hld.md). Names below are the ones to implement.

Stack: a React app, a Go HTTP API, PostgreSQL, S3, and email. The Go process serves `/api` and, in production, the built React files. The reminder loop is a goroutine with a one-minute ticker in that same process, guarded so two ticks cannot claim the same row.

React calls the Go API with the session cookie. It never talks to PostgreSQL or S3 directly. The Go API is the only holder of the database URL and the S3 credentials.

A user row is a lawyer. There is no chamber table and no role column.

## 1. Modules

| Module | Writes | Reads |
| --- | --- | --- |
| `session` | users on sign-up / password change | users |
| `clients` | clients | clients, cases |
| `cases` | cases | cases, clients, hearings |
| `documents` | documents, S3 | documents |
| `hearings` | hearings, reminders, `cases.next_hearing_at` | hearings, cases |
| `reminders` | reminder status, notifications | reminders, hearings, cases, clients, users |
| `diary` | nothing | hearings, cases, clients |

Route handlers stay thin: read the session, call one module, map errors to HTTP. Modules receive `userId` as an argument. They do not read the cookie themselves.

## 2. Schema

All primary keys are UUID. All timestamps are `timestamptz`.

`owner_id` on clients and cases is the lawyer. Child rows (documents, hearings) do not repeat it; access is checked by joining to the case.

### users

| Column | Notes |
| --- | --- |
| id | The lawyer. |
| name | |
| email | Unique. |
| password_hash | |
| created_at | |

### clients

| Column | Notes |
| --- | --- |
| id | |
| owner_id | FK users. The lawyer who keeps this client. |
| full_name | Required. |
| phone | Optional. |
| email | Optional. |
| address | Optional. |
| notes | Optional. |
| created_at, updated_at | |

Index: `(owner_id, full_name)`.

### cases

| Column | Notes |
| --- | --- |
| id | |
| owner_id | FK users. Same lawyer as the client. |
| client_id | FK clients. `clients.owner_id` must equal `cases.owner_id`. |
| title | Short name the diary can show. |
| case_number | Court number. Empty until allotted. |
| court_name | |
| case_type | `civil`, `criminal`, `writ`, `family`, `consumer`, `other`. |
| status | `active`, `stayed`, `disposed`. |
| opposite_party | |
| filing_date | Date, nullable. |
| next_hearing_at | Cache. Null when nothing is scheduled in the future. |
| notes | |
| created_at, updated_at | |

Indexes: `(owner_id, status)`, `(client_id)`, `(owner_id, next_hearing_at)`.

### documents

| Column | Notes |
| --- | --- |
| id | |
| case_id | FK cases. |
| original_name | Name from the browser. |
| content_type | |
| byte_size | |
| storage_key | S3 object key. See upload flow. |
| label | `petition`, `order`, `evidence`, `vakalatnama`, `other`. |
| created_at | |

Index: `(case_id, created_at desc)`.

### hearings

| Column | Notes |
| --- | --- |
| id | |
| case_id | FK cases. |
| scheduled_at | |
| court_name | Copied from the case, editable if the matter is transferred. |
| purpose | `mention`, `arguments`, `evidence`, `orders`, `compliance`, `other`. |
| status | `scheduled`, `completed`, `adjourned`, `cancelled`. |
| outcome | Set when completed or adjourned. |
| previous_hearing_id | FK hearings, null on the first date. Set on the new row when adjourning. |
| created_at, updated_at | |

Indexes: `(case_id, scheduled_at)`, `(scheduled_at)` for the diary.

### reminders

One lawyer per hearing, so a hearing has at most two reminder rows.

| Column | Notes |
| --- | --- |
| id | |
| hearing_id | FK hearings. |
| kind | `eve` or `soon`. |
| remind_at | |
| status | `scheduled`, `sending`, `sent`, `failed`, `cancelled`. |
| attempts | Starts at 0. |
| last_error | |
| sent_at | |
| claimed_at | Set when status becomes `sending`, so a crash can return the row to `scheduled` after 10 minutes. |

Unique: `(hearing_id, kind)`.

Index for the worker: `(status, remind_at)`.

### notifications

| Column | Notes |
| --- | --- |
| id | |
| user_id | The lawyer who owns the hearing. |
| reminder_id | |
| title | |
| body | |
| read_at | Null until opened. |
| created_at | |

## 3. State

### Case

`active` and `stayed` may still receive hearings. `disposed` may not. Moving back from `disposed` to `active` is allowed and does not invent hearings.

### Hearing

```text
scheduled → completed
scheduled → adjourned     (and a new scheduled hearing is inserted)
scheduled → cancelled
```

`completed`, `adjourned`, and `cancelled` are terminal. A wrong outcome is corrected by editing `outcome` text, not by reopening status.

### Reminder

```text
scheduled → sending → sent
sending   → scheduled   (attempt failed, attempts < 3)
sending   → failed      (attempts reaches 3)
scheduled → cancelled   (hearing no longer scheduled)
```

`sending` is the claim. A second tick skips rows that are not `scheduled`.

## 4. Flows

### 4.1 Sign up

1. Reject a duplicate email with 409.
2. Insert the user with a slow password hash.
3. Set the session cookie. The new diary is empty.

### 4.2 Create a case

1. Load the client by id and session `userId`. Missing client → 404.
2. Insert the case with `owner_id = userId`, `status = active`, and `next_hearing_at = null`.
3. Return the case.

### 4.3 Upload

1. Require a case whose `owner_id` is the session user.
2. Reject if the file is over 20 MB.
3. Reject unless the extension is `pdf`. A stated content type must be `application/pdf`. An empty type or `application/octet-stream` is accepted, because browsers often send that for a valid file. |
4. Build the S3 object key as `lawyers/{userId}/cases/{caseId}/{documentId}/{safeName}`.
5. `PutObject` into the documents bucket.
6. Insert the `documents` row with that key.
7. If the insert fails, `DeleteObject`, then return 500.

Download loads the document through the case and checks `owner_id`, then returns a presigned GET URL that expires in five minutes. The object key itself is never returned, and the bucket is not public.

### 4.4 Add a hearing

In one transaction:

1. Lock the case row (`SELECT … FOR UPDATE`) where `owner_id` is the session user. Reject if missing, or if status is `disposed`.
2. Insert the hearing with `status = scheduled`.
3. Run the reminder planner for that hearing.
4. Recompute `cases.next_hearing_at` as the minimum `scheduled_at` of hearings on this case with `status = scheduled` and `scheduled_at > now()`.

If `scheduled_at` is in the past, the hearing is still saved (backfill) and the planner writes nothing. The response includes `remindersPlanned: 0` so the screen can say that no notice will fire.

### 4.5 Adjourn

In one transaction:

1. Load the hearing through its case. Reject unless the case owner is the session user and the hearing status is `scheduled`.
2. Require `nextAt` later than the current `scheduled_at`.
3. Set the old hearing to `adjourned` and store `outcome`.
4. Cancel reminders for the old hearing where status is `scheduled`.
5. Insert the new hearing: `status = scheduled`, `scheduled_at = nextAt`, `previous_hearing_id` = old id, purpose copied unless overridden.
6. Plan reminders for the new hearing.
7. Recompute `next_hearing_at`.

### 4.6 Complete or cancel

Same transaction shape as adjourn, without a successor hearing. Status becomes `completed` or `cancelled`. Pending reminders are cancelled. `next_hearing_at` is recomputed.

### 4.7 Reminder planner

Called only from the hearing transaction. It writes reminders for the case owner and for nobody else.

```text
cancel every reminder of this hearing where status = scheduled

if hearing.status is not scheduled: return
if hearing.scheduled_at <= now: return

eve  = 18:00 Asia/Kolkata on the calendar day before scheduled_at
soon = scheduled_at minus 2 hours

for each (kind, remind_at) in { eve, soon }:
  if remind_at > now:
    insert reminder status = scheduled
```

Unique `(hearing_id, kind)` makes a retry of the planner safe after the cancel.

### 4.8 Reminder tick

```text
every 60 seconds:
  select id from reminders
    where status = scheduled and remind_at <= now()
    order by remind_at
    limit 50
    for update skip locked

  for each id:
    set status = sending, attempts = attempts + 1
    commit that claim

    load hearing, case, client, and the case owner
    insert notification for that owner
    send email to that owner's address

    on success: status = sent, sent_at = now()
    on failure and attempts < 3: status = scheduled, last_error set
    on failure and attempts >= 3: status = failed, last_error set
```

The claim commits before the email send so two overlapping ticks do not both send. If the process dies after claim and before the final status write, the row stays `sending`. A sweep at the start of each tick returns rows stuck in `sending` for more than 10 minutes to `scheduled` without incrementing attempts again. That sweep is the recovery path for a crash, and it is why a duplicate email is possible and acceptable.

Email body is plain text: client name, case title, case number, court, purpose, and the sitting time in Asia/Kolkata. Inbox title is `Hearing · {case title}`.

## 5. HTTP

All routes require a session except sign-up and sign-in. An id the session lawyer does not own returns 404.

| Method | Path | Body / result |
| --- | --- | --- |
| POST | `/api/sign-up` | name, email, password → session cookie |
| POST | `/api/session` | email, password → session cookie |
| DELETE | `/api/session` | clears the cookie |
| GET | `/api/clients?q=` | this lawyer's clients, name search |
| POST | `/api/clients` | full_name, phone, email, address, notes |
| GET | `/api/clients/:id` | client plus its cases |
| PATCH | `/api/clients/:id` | same fields, partial |
| POST | `/api/cases` | client_id, title, court_name, case_type, and optional case fields |
| GET | `/api/cases?status=&clientId=&q=` | list rows include `next_hearing_at` and client name |
| GET | `/api/cases/:id` | case, client, documents, hearings newest first |
| PATCH | `/api/cases/:id` | case fields except `owner_id` and `next_hearing_at` |
| POST | `/api/cases/:id/documents` | multipart file + label |
| GET | `/api/documents/:id` | five-minute S3 presigned GET URL |
| DELETE | `/api/documents/:id` | deletes the S3 object, then the row |
| POST | `/api/cases/:id/hearings` | scheduled_at, court_name, purpose |
| POST | `/api/hearings/:id/complete` | outcome |
| POST | `/api/hearings/:id/cancel` | outcome optional |
| POST | `/api/hearings/:id/adjourn` | next_at, outcome, optional purpose |
| GET | `/api/diary?from=&to=` | this lawyer's scheduled hearings in range |
| GET | `/api/notifications` | unread first |
| POST | `/api/notifications/:id/read` | sets read_at |

Diary `from` and `to` are inclusive instants. The home screen asks for start of today through seven days, in Asia/Kolkata, converted to UTC before the query.

## 6. Diary query

```sql
select h.id, h.scheduled_at, h.court_name, h.purpose,
       c.id as case_id, c.title, c.case_number,
       cl.full_name as client_name
from hearings h
join cases c on c.id = h.case_id
join clients cl on cl.id = c.client_id
where c.owner_id = $1
  and h.status = 'scheduled'
  and h.scheduled_at >= $2
  and h.scheduled_at < $3
order by h.scheduled_at;
```

`$1` is the session user id.

## 7. Validation errors

| Condition | HTTP |
| --- | --- |
| Missing required field, bad enum, file type, file size, adjourn `nextAt` not later | 400 |
| No session | 401 |
| Id not owned by this lawyer | 404 |
| Email already registered | 409 |
| Case is `disposed` and the caller is adding a hearing | 409 |
| Delete client while any case points at them | 409 |

The client delete conflict is the only delete rule in v1. Documents otherwise delete as written above. Clients are edited, not removed, once a case exists.

## 8. Screens

| Screen | Data |
| --- | --- |
| Diary | `GET /api/diary` for today through seven days, plus unread notification count |
| Clients | Search list. Row opens the client and their cases. |
| Case file | Details, document list with upload, hearing list with add / complete / adjourn |
| Inbox | This lawyer's notifications. Opening one marks it read and links to the case. |

The case file is the only place a hearing is created. The diary links to the case; it does not edit the hearing in place.

## 9. Security

- Password hashing with a slow hash (argon2 or bcrypt). Cookie is `HttpOnly`, `Secure`, `SameSite=Lax`.
- `owner_id` is never accepted from the client. It is the session user id.
- S3 object keys are server-generated. `original_name` is displayed only, after stripping path separators.
- The bucket blocks public access. The React app receives only a five-minute presigned GET URL after the Go API has checked `owner_id`.
- Download and delete load the document through the case and check `owner_id`. Delete calls `DeleteObject`, then removes the row.
- The reminder tick is not an HTTP route. It is not reachable from the browser.
- A lawyer cannot list, open, or search another lawyer's clients or cases. The response is 404, not 403, so existence of the other lawyer's ids is not revealed.

## 10. What this design deliberately leaves unbuilt

No clerk role, no shared chamber, no audit table, no full-text search, and no way to transfer a case to another lawyer. A reminder always belongs to one hearing, and replanning always cancels unsent rows before inserting new ones.
