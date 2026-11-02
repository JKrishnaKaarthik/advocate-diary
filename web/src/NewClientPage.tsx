import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ApiError, api } from "./api";
import { FieldHint, useFieldHints } from "./fields";

const clientRules = [
  { name: "full_name", label: "Name", required: true },
  { name: "email", label: "Email", email: true },
];

export function NewClientPage() {
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);
  const fields = useFieldHints();

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    if (fields.reject(form, clientRules)) return;
    const data = new FormData(form);
    setError("");
    setPending(true);
    try {
      const client = await api.createClient({
        full_name: String(data.get("full_name") || "").trim(),
        phone: String(data.get("phone") || "").trim(),
        email: String(data.get("email") || "").trim(),
        address: String(data.get("address") || "").trim(),
        notes: String(data.get("notes") || "").trim(),
      });
      navigate(`/clients/${client.id}`);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Could not save the client");
      setPending(false);
    }
  }

  return (
    <section>
      <header className="page-head">
        <p className="crumb"><Link to="/clients">Clients</Link></p>
        <h1 data-testid="new-client-heading">Add a client</h1>
        <p>Name is required. Phone, email, and address can wait.</p>
      </header>
      <form className="card stack" noValidate onSubmit={onSubmit} onInput={fields.onInput}>
        <label>
          Name
          <input name="full_name" autoFocus />
          <FieldHint text={fields.hints.full_name} />
        </label>
        <div className="split">
          <label>Phone<input name="phone" /></label>
          <label>
            Email
            <input name="email" type="email" />
            <FieldHint text={fields.hints.email} />
          </label>
        </div>
        <label>Address<input name="address" /></label>
        <label>Notes<textarea name="notes" rows={3} /></label>
        {error ? <p className="error">{error}</p> : null}
        <div className="actions">
          <button type="submit" data-testid="save-client" disabled={pending}>{pending ? "Saving…" : "Save client"}</button>
          <Link className="button secondary" to="/clients" data-testid="cancel-client">Cancel</Link>
        </div>
      </form>
    </section>
  );
}
