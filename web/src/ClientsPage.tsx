import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Client, api } from "./api";

export function ClientsPage() {
  const [q, setQ] = useState("");
  const [items, setItems] = useState<Client[]>([]);
  const [error, setError] = useState("");
  async function load(query = q) {
    setItems(await api.clients(query));
  }
  useEffect(() => {
    load("").catch((err) => setError(err.message));
  }, []);
  return (
    <section>
      <header className="page-head with-action">
        <div>
          <h1 data-testid="clients-heading">Clients</h1>
          <p>People whose matters you are handling.</p>
        </div>
        <Link className="button" to="/clients/new" data-testid="add-client">Add client</Link>
      </header>
      <form className="search" onSubmit={(event) => { event.preventDefault(); load().catch((err) => setError(err.message)); }}>
        <input value={q} onChange={(event) => setQ(event.target.value)} placeholder="Search by name" aria-label="Search clients" />
        <button type="submit" className="secondary">Search</button>
      </form>
      {error ? <p className="error">{error}</p> : null}
      {items.length === 0 ? <p className="muted">No clients yet.</p> : null}
      <ul className="rows">
        {items.map((client) => (
          <li key={client.id}>
            <Link to={`/clients/${client.id}`}>
              <strong>{client.full_name}</strong>
              <span>{[client.phone, client.email].filter(Boolean).join(" · ") || "No contact saved"}</span>
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}
