import { FormEvent, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { ApiError, CaseDetail, api, formatWhen, fromLocalInput } from "./api";
import { FieldHint, useFieldHints } from "./fields";

const purposes = ["mention", "arguments", "evidence", "orders", "compliance", "other"];
const labels = ["petition", "order", "evidence", "vakalatnama", "other"];
const caseTypes = ["civil", "criminal", "writ", "family", "consumer", "other"];
const statuses = ["active", "stayed", "disposed"];

export function CasePage() {
  const { id = "" } = useParams();
  const [item, setItem] = useState<CaseDetail | null>(null);
  const [error, setError] = useState("");
  const [caseError, setCaseError] = useState("");
  const [uploadError, setUploadError] = useState("");
  const [hearingError, setHearingError] = useState("");
  const [note, setNote] = useState("");
  const caseFields = useFieldHints();
  const uploadFields = useFieldHints();
  const hearingFields = useFieldHints();
  async function load() {
    setItem(await api.case(id));
  }
  useEffect(() => {
    load().catch((err) => setError(err.message));
  }, [id]);
  if (!item && error) return <p className="error">{error}</p>;
  if (!item) return <p className="muted">Loading case…</p>;

  async function save(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    if (caseFields.reject(form, [
      { name: "title", label: "Title", required: true },
      { name: "court_name", label: "Court", required: true },
    ])) return;
    const data = new FormData(form);
    setCaseError("");
    try {
      await api.updateCase(id, {
        title: String(data.get("title") || ""),
        court_name: String(data.get("court_name") || ""),
        case_type: String(data.get("case_type") || ""),
        status: String(data.get("status") || ""),
        case_number: String(data.get("case_number") || ""),
        opposite_party: String(data.get("opposite_party") || ""),
        filing_date: String(data.get("filing_date") || ""),
        notes: String(data.get("notes") || ""),
      });
      await load();
    } catch (err) {
      setCaseError(err instanceof ApiError ? err.message : "Could not save the case");
    }
  }
  async function upload(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);
    if (uploadFields.reject(form, [{ name: "file", label: "File", file: true }])) return;
    const file = data.get("file");
    if (!(file instanceof File)) return;
    setUploadError("");
    try {
      await api.upload(id, file, String(data.get("label") || "other"));
      form.reset();
      await load();
    } catch (err) {
      setUploadError(err instanceof ApiError ? err.message : "Upload failed");
    }
  }
  async function addHearing(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    if (hearingFields.reject(form, [
      { name: "scheduled_at", label: "Date and time", required: true },
      { name: "court_name", label: "Court", required: true },
    ])) return;
    const data = new FormData(form);
    setHearingError("");
    setNote("");
    try {
      const created = await api.addHearing(id, {
        scheduled_at: fromLocalInput(String(data.get("scheduled_at") || "")),
        court_name: String(data.get("court_name") || ""),
        purpose: String(data.get("purpose") || ""),
      });
      setNote(created.reminders_planned === 0
        ? "Hearing saved. No reminder will fire; it is inside the lead time or already past."
        : `Hearing saved. ${created.reminders_planned} reminder${created.reminders_planned === 1 ? "" : "s"} planned.`);
      form.reset();
      await load();
    } catch (err) {
      setHearingError(err instanceof ApiError ? err.message : "Could not save the hearing");
    }
  }
  return (
    <section>
      <header className="page-head">
        <p className="crumb"><Link to={`/clients/${item.client_id}`}>{item.client_name}</Link></p>
        <h1 data-testid="case-heading">{item.title}</h1>
        <p>{item.court_name} · {item.case_type} · {item.status}</p>
      </header>
      {error ? <p className="error">{error}</p> : null}
      {note ? <p className="note">{note}</p> : null}
      <form className="card stack" key={item.id} noValidate onSubmit={save} onInput={caseFields.onInput}>
        <h2>Case details</h2>
        <label>
          Title
          <input name="title" defaultValue={item.title} />
          <FieldHint text={caseFields.hints.title} />
        </label>
        <div className="split">
          <label>
            Court
            <input name="court_name" defaultValue={item.court_name} />
            <FieldHint text={caseFields.hints.court_name} />
          </label>
          <label>Case number<input name="case_number" defaultValue={item.case_number} /></label>
        </div>
        <div className="split">
          <label>Type<select name="case_type" defaultValue={item.case_type}>{caseTypes.map((type) => <option key={type}>{type}</option>)}</select></label>
          <label>Status<select name="status" defaultValue={item.status}>{statuses.map((status) => <option key={status}>{status}</option>)}</select></label>
        </div>
        <div className="split">
          <label>Opposite party<input name="opposite_party" defaultValue={item.opposite_party} /></label>
          <label>Filing date<input name="filing_date" type="date" defaultValue={item.filing_date || ""} /></label>
        </div>
        <label>Notes<textarea name="notes" rows={3} defaultValue={item.notes} /></label>
        {caseError ? <p className="error">{caseError}</p> : null}
        <button type="submit" data-testid="save-case">Save case</button>
      </form>

      <h2>Papers</h2>
      {item.documents.length === 0 ? <p className="muted">No papers on this case.</p> : null}
      <ul className="rows">
        {item.documents.map((doc) => (
          <li key={doc.id} className="doc" data-testid="paper">
            <div>
              <strong data-testid="paper-name">{doc.original_name}</strong>
              <span data-testid="paper-label">{doc.label} · {Math.ceil(doc.byte_size / 1024)} KB</span>
            </div>
            <div className="actions">
              <button className="secondary" onClick={async () => {
                const result = await api.download(doc.id);
                window.open(result.url, "_blank", "noopener");
              }}>Open</button>
              <button className="text" onClick={async () => {
                await api.deleteDocument(doc.id);
                await load();
              }}>Delete</button>
            </div>
          </li>
        ))}
      </ul>
      <form className="card inline" noValidate onSubmit={upload} onInput={uploadFields.onInput} onChange={uploadFields.onInput}>
        <label>
          File
          <input data-testid="upload-file" name="file" type="file" accept=".pdf,.doc,.docx,.jpg,.jpeg,.png,.txt" />
          <FieldHint text={uploadFields.hints.file} />
        </label>
        <label>Label<select data-testid="upload-label" name="label" defaultValue="other">{labels.map((label) => <option key={label}>{label}</option>)}</select></label>
        {uploadError ? <p className="error">{uploadError}</p> : null}
        <button type="submit" data-testid="upload">Upload</button>
      </form>

      <h2>Hearings</h2>
      {item.hearings.length === 0 ? <p className="muted">No hearing fixed.</p> : null}
      <ul className="rows">
        {item.hearings.map((hearing) => (
          <HearingRow key={hearing.id} hearing={hearing} onChange={load} onError={setError} />
        ))}
      </ul>
      {item.status === "disposed" ? <p className="muted">This case is disposed, so a new hearing cannot be added.</p> : (
        <form className="card stack" noValidate onSubmit={addHearing} onInput={hearingFields.onInput}>
          <h2 data-testid="fix-hearing-heading">Fix a hearing</h2>
          <div className="split">
            <label>
              When, India time
              <input data-testid="hearing-when" name="scheduled_at" type="datetime-local" />
              <FieldHint text={hearingFields.hints.scheduled_at} />
            </label>
            <label>
              Court
              <input data-testid="hearing-court" name="court_name" defaultValue={item.court_name} />
              <FieldHint text={hearingFields.hints.court_name} />
            </label>
          </div>
          <label>Purpose<select data-testid="hearing-purpose" name="purpose" defaultValue="arguments">{purposes.map((purpose) => <option key={purpose}>{purpose}</option>)}</select></label>
          {hearingError ? <p className="error">{hearingError}</p> : null}
          <button type="submit" data-testid="save-hearing">Save hearing</button>
        </form>
      )}
    </section>
  );
}

function HearingRow({ hearing, onChange, onError }: {
  hearing: CaseDetail["hearings"][number];
  onChange: () => Promise<void>;
  onError: (message: string) => void;
}) {
  const [mode, setMode] = useState<"complete" | "adjourn" | "">("");
  const [problem, setProblem] = useState("");
  const fields = useFieldHints();
  function pick(next: "complete" | "adjourn") {
    fields.clear();
    setProblem("");
    setMode(mode === next ? "" : next);
  }
  async function run(action: () => Promise<unknown>) {
    try {
      await action();
      setMode("");
      await onChange();
    } catch (err) {
      const message = err instanceof ApiError ? err.message : "Could not update the hearing";
      setProblem(message);
      onError(message);
    }
  }
  return (
    <li data-testid="hearing">
      <div>
        <strong data-testid="hearing-when-status">{formatWhen(hearing.scheduled_at)} · {hearing.status}</strong>
        <span data-testid="hearing-court-purpose">{hearing.court_name} · {hearing.purpose}</span>
        {hearing.outcome ? <span>{hearing.outcome}</span> : null}
        {problem ? <p className="info">{problem}</p> : null}
      </div>
      {hearing.status === "scheduled" ? (
        <div className="actions">
          <button className="secondary" onClick={() => pick("complete")}>Complete</button>
          <button className="secondary" onClick={() => pick("adjourn")}>Adjourn</button>
          <button className="text" onClick={() => run(() => api.cancelHearing(hearing.id, ""))}>Cancel</button>
        </div>
      ) : null}
      {mode === "complete" ? (
        <form className="inline" noValidate onInput={fields.onInput} onSubmit={(event) => {
          event.preventDefault();
          const form = event.currentTarget;
          if (fields.reject(form, [{ name: "outcome", label: "Outcome", required: true }])) return;
          const outcome = String(new FormData(form).get("outcome") || "").trim();
          run(() => api.completeHearing(hearing.id, outcome));
        }}>
          <input name="outcome" className={fields.hints.outcome ? "invalid" : undefined} placeholder="What happened" />
          <FieldHint text={fields.hints.outcome} />
          <button type="submit">Save outcome</button>
        </form>
      ) : null}
      {mode === "adjourn" ? (
        <form className="inline" noValidate onInput={fields.onInput} onSubmit={(event) => {
          event.preventDefault();
          const form = event.currentTarget;
          if (fields.reject(form, [
            { name: "outcome", label: "Reason", required: true },
            { name: "next_at", label: "Next date and time", required: true },
          ])) return;
          const data = new FormData(form);
          run(() => api.adjournHearing(hearing.id, {
            next_at: fromLocalInput(String(data.get("next_at") || "")),
            outcome: String(data.get("outcome") || "").trim(),
          }));
        }}>
          <input name="outcome" className={fields.hints.outcome ? "invalid" : undefined} placeholder="Why it was adjourned" />
          <input name="next_at" className={fields.hints.next_at ? "invalid" : undefined} type="datetime-local" aria-label="Next date, India time" />
          <FieldHint text={fields.hints.outcome} />
          <FieldHint text={fields.hints.next_at} />
          <button type="submit">Fix next date</button>
        </form>
      ) : null}
    </li>
  );
}
