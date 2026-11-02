import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { DiaryEntry, api, monthWindow, shortWhen, weekBuckets } from "./api";

function Contact({ label, children }: { label: string; children: string }) {
  return (
    <div>
      <dt>{label}</dt>
      <dd>{children || "Not on file"}</dd>
    </div>
  );
}

export function InboxPage() {
  const [items, setItems] = useState<DiaryEntry[] | null>(null);
  const [error, setError] = useState("");
  useEffect(() => {
    const range = monthWindow();
    api.diary(range.from, range.to).then(setItems).catch((err) => setError(err.message));
  }, []);
  const weeks = weekBuckets(items || []);
  return (
    <section>
      <header className="page-head">
        <h1 data-testid="inbox-heading">Client reminders</h1>
        <p>Hearings in the coming month. Name, case type, time, and how to reach the client.</p>
      </header>
      {error ? <p className="error">{error}</p> : null}
      {items && items.length === 0 ? <p className="muted">No hearings in the next month.</p> : null}
      {weeks.map((week) => (
        <section key={week.title} className="week">
          <h2>{week.title}</h2>
          <div className="board">
            {week.rows.map((row) => (
              <article key={row.id} className="hearing-card" data-testid="inbox-hearing">
                <header className="card-top">
                  <time>{shortWhen(row.scheduled_at)}</time>
                  <span className="pill">{row.case_type}</span>
                </header>
                <h3 data-testid="inbox-client">{row.client_name}</h3>
                <p className="matter" data-testid="inbox-matter">{row.title}{row.case_number ? ` · ${row.case_number}` : ""}</p>
                <dl className="contact-list">
                  <Contact label="Phone">{row.phone}</Contact>
                  <Contact label="Email">{row.email}</Contact>
                  <Contact label="Address">{row.address}</Contact>
                </dl>
                <Link to={`/cases/${row.case_id}`}>Open case</Link>
              </article>
            ))}
          </div>
        </section>
      ))}
    </section>
  );
}
