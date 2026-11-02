import { FormEvent, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { ApiError, ClientDetail, api, formatWhen } from "./api";
import { FieldHint, useFieldHints } from "./fields";

const caseTypes = ["civil", "criminal", "writ", "family", "consumer", "other"];

export function ClientPage() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const [client, setClient] = useState<ClientDetail | null>(null);
  const [error, setError] = useState("");
  const [caseError, setCaseError] = useState("");
  const clientFields = useFieldHints();
  const caseFields = useFieldHints();
  useEffect(() => {
    api.client(id).then(setClient).catch((err) => setError(err.message));
  }, [id]);
  if (error && !client) return <p className="error">{error}</p>;
  if (!client) return <p className="muted">Loading client…</p>;
  async function save(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    if (clientFields.reject(form, [
      { name: "full_name", label: "Name", required: true },
      { name: "email", label: "Email", email: true },
    ])) return;
    const data = new FormData(form);
    setError("");
    try {
      await api.updateClient(id, {
        full_name: String(data.get("full_name") || "").trim(),
        phone: String(data.get("phone") || "").trim(),
        email: String(data.get("email") || "").trim(),
        address: String(data.get("address") || "").trim(),
        notes: String(data.get("notes") || "").trim(),
      });
      setClient(await api.client(id));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Could not save");
    }
  }
  async function openCase(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    if (caseFields.reject(form, [
      { name: "title", label: "Title", required: true },
      { name: "court_name", label: "Court", required: true },
    ])) return;
    const data = new FormData(form);
    setCaseError("");
    try {
      const created = await api.createCase({
        client_id: id,
        title: String(data.get("title") || ""),
        court_name: String(data.get("court_name") || ""),
        case_type: String(data.get("case_type") || ""),
        case_number: String(data.get("case_number") || ""),
        opposite_party: String(data.get("opposite_party") || ""),
      });
      navigate(`/cases/${created.id}`);
    } catch (err) {
      setCaseError(err instanceof ApiError ? err.message : "Could not open the case");
    }
  }
  return (
    <section>
      <header className="page-head">
        <h1 data-testid="client-heading">{client.full_name}</h1>
        <p>Client record and the matters tied to them.</p>
      </header>
      <form className="card stack" noValidate onSubmit={save} onInput={clientFields.onInput}>
        <label>
          Name
          <input name="full_name" defaultValue={client.full_name} />
          <FieldHint text={clientFields.hints.full_name} />
        </label>
        <div className="split">
          <label>Phone<input name="phone" defaultValue={client.phone} /></label>
          <label>
            Email
            <input name="email" type="email" defaultValue={client.email} />
            <FieldHint text={clientFields.hints.email} />
          </label>
        </div>
        <label>Address<input name="address" defaultValue={client.address} /></label>
        <label>Notes<textarea name="notes" rows={3} defaultValue={client.notes} /></label>
        {error ? <p className="error">{error}</p> : null}
        <button type="submit">Save changes</button>
      </form>
      <h2>Cases</h2>
      {client.cases.length === 0 ? <p className="muted">No case yet.</p> : null}
      <ul className="rows">
        {client.cases.map((item) => (
          <li key={item.id}>
            <Link to={`/cases/${item.id}`}>
              <strong>{item.title}</strong>
              <span>{item.court_name} · {item.status}{item.case_number ? ` · ${item.case_number}` : ""}</span>
              <span>{item.next_hearing_at ? `Next hearing ${formatWhen(item.next_hearing_at)}` : "No hearing fixed"}</span>
            </Link>
          </li>
        ))}
      </ul>
      <form className="card stack" noValidate onSubmit={openCase} onInput={caseFields.onInput}>
        <h2 data-testid="open-case-heading">Open a case</h2>
        <label>
          Title
          <input name="title" />
          <FieldHint text={caseFields.hints.title} />
        </label>
        <div className="split">
          <label>
            Court
            <input name="court_name" />
            <FieldHint text={caseFields.hints.court_name} />
          </label>
          <label>
            Type
            <select name="case_type" defaultValue="civil">
              {caseTypes.map((type) => <option key={type} value={type}>{type}</option>)}
            </select>
          </label>
        </div>
        <div className="split">
          <label>Case number<input name="case_number" /></label>
          <label>Opposite party<input name="opposite_party" /></label>
        </div>
        {caseError ? <p className="error">{caseError}</p> : null}
        <button type="submit" data-testid="open-case">Open case</button>
      </form>
    </section>
  );
}
