import { FormEvent, createContext, useContext, useEffect, useState } from "react";
import { Link, NavLink, Navigate, Outlet, Route, Routes, useNavigate } from "react-router-dom";
import { ApiError, Lawyer, api, monthWindow } from "./api";
import { FieldHint, useFieldHints } from "./fields";
import { CasePage } from "./CasePage";
import { ClientPage } from "./ClientPage";
import { ClientsPage } from "./ClientsPage";
import { NewClientPage } from "./NewClientPage";
import { DiaryPage } from "./DiaryPage";
import { InboxPage } from "./InboxPage";

type Session = { lawyer: Lawyer | null; setLawyer: (lawyer: Lawyer | null) => void };

const SessionContext = createContext<Session | null>(null);

export function useSession() {
  const value = useContext(SessionContext);
  if (!value) throw new Error("session missing");
  return value;
}

export function App() {
  const [lawyer, setLawyer] = useState<Lawyer | null | undefined>(undefined);
  useEffect(() => {
    api.session().then(setLawyer).catch(() => setLawyer(null));
  }, []);
  if (lawyer === undefined) return <p className="boot">Opening the diary…</p>;
  return (
    <SessionContext.Provider value={{ lawyer, setLawyer }}>
      <Routes>
        <Route path="/sign-in" element={<Auth mode="in" />} />
        <Route path="/sign-up" element={<Auth mode="up" />} />
        <Route element={<Layout />}>
          <Route path="/" element={<DiaryPage />} />
          <Route path="/clients" element={<ClientsPage />} />
          <Route path="/clients/new" element={<NewClientPage />} />
          <Route path="/clients/:id" element={<ClientPage />} />
          <Route path="/cases/:id" element={<CasePage />} />
          <Route path="/inbox" element={<InboxPage />} />
        </Route>
      </Routes>
    </SessionContext.Provider>
  );
}

function Layout() {
  const { lawyer, setLawyer } = useSession();
  const navigate = useNavigate();
  const [unread, setUnread] = useState(0);
  useEffect(() => {
    if (!lawyer) return;
    const range = monthWindow();
    api.diary(range.from, range.to).then((items) => setUnread(items.length)).catch(() => setUnread(0));
  }, [lawyer]);
  if (!lawyer) return <Navigate to="/sign-in" replace />;
  return (
    <div className="shell">
      <aside className="side">
        <p className="mark">Advocate Diary</p>
        <nav>
          <NavLink to="/" end data-testid="nav-diary">Diary</NavLink>
          <NavLink to="/clients" data-testid="nav-clients">Clients</NavLink>
          <NavLink to="/inbox" data-testid="nav-inbox">Inbox{unread > 0 ? <span className="count">{unread}</span> : null}</NavLink>
        </nav>
        <div className="who">
          <strong>{lawyer.name}</strong>
          <span>{lawyer.email}</span>
          <button
            className="text"
            data-testid="sign-out"
            onClick={async () => {
              await api.signOut();
              setLawyer(null);
              navigate("/sign-in");
            }}
          >
            Sign out
          </button>
        </div>
      </aside>
      <main>
        <Outlet />
      </main>
    </div>
  );
}

function Auth({ mode }: { mode: "in" | "up" }) {
  const { lawyer, setLawyer } = useSession();
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);
  const fields = useFieldHints();
  if (lawyer) return <Navigate to="/" replace />;
  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const rules = mode === "up"
      ? [
          { name: "name", label: "Name", required: true },
          { name: "email", label: "Email", required: true, email: true },
          { name: "password", label: "Password", required: true, min: 8 },
        ]
      : [
          { name: "email", label: "Email", required: true, email: true },
          { name: "password", label: "Password", required: true },
        ];
    if (fields.reject(form, rules)) return;
    const data = new FormData(form);
    setPending(true);
    setError("");
    try {
      const next =
        mode === "up"
          ? await api.signUp({
              name: String(data.get("name") || ""),
              email: String(data.get("email") || ""),
              password: String(data.get("password") || ""),
            })
          : await api.signIn({
              email: String(data.get("email") || ""),
              password: String(data.get("password") || ""),
            });
      setLawyer(next);
      navigate("/");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Could not sign in");
    } finally {
      setPending(false);
    }
  }
  return (
    <div className="auth">
      <form className="card" noValidate onSubmit={onSubmit} onInput={fields.onInput}>
        <p className="mark">Advocate Diary</p>
        <h1>{mode === "up" ? "Open your diary" : "Welcome back"}</h1>
        <p className="lede">
          {mode === "up"
            ? "Your clients, cases, and hearing dates stay on this account alone."
            : "Sign in to the diary you keep for your own matters."}
        </p>
        {mode === "up" ? (
          <label>
            Name
            <input name="name" autoComplete="name" />
            <FieldHint text={fields.hints.name} />
          </label>
        ) : null}
        <label>
          Email
          <input name="email" type="email" autoComplete="email" />
          <FieldHint text={fields.hints.email} />
        </label>
        <label>
          Password
          <input name="password" type="password" autoComplete={mode === "up" ? "new-password" : "current-password"} />
          <FieldHint text={fields.hints.password} />
        </label>
        {error ? <p className="error">{error}</p> : null}
        <button type="submit" data-testid="auth-submit" disabled={pending}>{pending ? "Please wait…" : mode === "up" ? "Create account" : "Sign in"}</button>
        <p className="switch">
          {mode === "up" ? <Link to="/sign-in" data-testid="go-sign-in">I already have an account</Link> : <Link to="/sign-up" data-testid="go-sign-up">Create an account</Link>}
        </p>
      </form>
    </div>
  );
}
