CREATE TABLE users (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE clients (
    id uuid PRIMARY KEY,
    owner_id uuid NOT NULL REFERENCES users (id),
    full_name text NOT NULL,
    phone text NOT NULL DEFAULT '',
    email text NOT NULL DEFAULT '',
    address text NOT NULL DEFAULT '',
    notes text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX clients_owner_name ON clients (owner_id, full_name);

CREATE TABLE cases (
    id uuid PRIMARY KEY,
    owner_id uuid NOT NULL REFERENCES users (id),
    client_id uuid NOT NULL REFERENCES clients (id),
    title text NOT NULL,
    case_number text NOT NULL DEFAULT '',
    court_name text NOT NULL,
    case_type text NOT NULL CHECK (case_type IN ('civil', 'criminal', 'writ', 'family', 'consumer', 'other')),
    status text NOT NULL CHECK (status IN ('active', 'stayed', 'disposed')),
    opposite_party text NOT NULL DEFAULT '',
    filing_date date,
    next_hearing_at timestamptz,
    notes text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX cases_owner_status ON cases (owner_id, status);
CREATE INDEX cases_client ON cases (client_id);
CREATE INDEX cases_owner_next ON cases (owner_id, next_hearing_at);

CREATE TABLE documents (
    id uuid PRIMARY KEY,
    case_id uuid NOT NULL REFERENCES cases (id),
    original_name text NOT NULL,
    content_type text NOT NULL,
    byte_size bigint NOT NULL,
    storage_key text NOT NULL,
    label text NOT NULL CHECK (label IN ('petition', 'order', 'evidence', 'vakalatnama', 'other')),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX documents_case_created ON documents (case_id, created_at DESC);

CREATE TABLE hearings (
    id uuid PRIMARY KEY,
    case_id uuid NOT NULL REFERENCES cases (id),
    scheduled_at timestamptz NOT NULL,
    court_name text NOT NULL,
    purpose text NOT NULL CHECK (purpose IN ('mention', 'arguments', 'evidence', 'orders', 'compliance', 'other')),
    status text NOT NULL CHECK (status IN ('scheduled', 'completed', 'adjourned', 'cancelled')),
    outcome text NOT NULL DEFAULT '',
    previous_hearing_id uuid REFERENCES hearings (id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX hearings_case_at ON hearings (case_id, scheduled_at);
CREATE INDEX hearings_at ON hearings (scheduled_at);

CREATE TABLE reminders (
    id uuid PRIMARY KEY,
    hearing_id uuid NOT NULL REFERENCES hearings (id),
    kind text NOT NULL CHECK (kind IN ('eve', 'soon')),
    remind_at timestamptz NOT NULL,
    status text NOT NULL CHECK (status IN ('scheduled', 'sending', 'sent', 'failed', 'cancelled')),
    attempts int NOT NULL DEFAULT 0,
    last_error text NOT NULL DEFAULT '',
    sent_at timestamptz,
    claimed_at timestamptz,
    UNIQUE (hearing_id, kind)
);

CREATE INDEX reminders_due ON reminders (status, remind_at);

CREATE TABLE notifications (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users (id),
    reminder_id uuid NOT NULL REFERENCES reminders (id),
    title text NOT NULL,
    body text NOT NULL,
    read_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX notifications_user ON notifications (user_id, created_at DESC);
