import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { DiaryEntry, api, clock, dayKey, diaryWindow, monthDay, weekdayName } from "./api";

export function DiaryPage() {
  const [items, setItems] = useState<DiaryEntry[] | null>(null);
  const [error, setError] = useState("");
  useEffect(() => {
    const range = diaryWindow();
    api.diary(range.from, range.to).then(setItems).catch((err) => setError(err.message));
  }, []);
  const groups = new Map<string, DiaryEntry[]>();
  for (const item of items || []) {
    const key = dayKey(item.scheduled_at);
    groups.set(key, [...(groups.get(key) || []), item]);
  }
  return (
    <section>
      <header className="page-head">
        <h1 data-testid="diary-heading">Court diary</h1>
        <p>Sittings for the next seven days. Time, matter, and court.</p>
      </header>
      {error ? <p className="error">{error}</p> : null}
      {items && items.length === 0 ? (
        <div className="empty">
          <p>Nothing is listed for the next seven days.</p>
          <Link to="/clients">Open a client and fix a hearing</Link>
        </div>
      ) : null}
      {[...groups.entries()].map(([day, rows]) => (
        <section key={day} className="timetable">
          <header className="day-head">
            <h2>{weekdayName(day)}</h2>
            <span>{monthDay(day)}</span>
          </header>
          <ul>
            {rows.map((row) => (
              <li key={row.id}>
                <Link className="slot" data-testid="diary-slot" to={`/cases/${row.case_id}`}>
                  <time>{clock(row.scheduled_at)}</time>
                  <strong data-testid="diary-client">{row.client_name}</strong>
                  <span data-testid="diary-matter">{row.title}{row.case_number ? ` · ${row.case_number}` : ""}</span>
                  <span data-testid="diary-court">{row.court_name} · {row.purpose}</span>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      ))}
    </section>
  );
}
