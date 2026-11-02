export type Lawyer = { id: string; name: string; email: string };

export type Client = {
  id: string;
  full_name: string;
  phone: string;
  email: string;
  address: string;
  notes: string;
};

export type CaseSummary = {
  id: string;
  title: string;
  case_number: string;
  court_name: string;
  case_type: string;
  status: string;
  opposite_party: string;
  filing_date: string | null;
  next_hearing_at: string | null;
  notes: string;
  client_id: string;
  client_name: string;
};

export type Document = {
  id: string;
  original_name: string;
  content_type: string;
  byte_size: number;
  label: string;
  created_at: string;
};

export type Hearing = {
  id: string;
  scheduled_at: string;
  court_name: string;
  purpose: string;
  status: string;
  outcome: string;
  previous_hearing_id: string | null;
};

export type CaseDetail = CaseSummary & { documents: Document[]; hearings: Hearing[] };
export type ClientDetail = Client & { cases: CaseSummary[] };

export type DiaryEntry = {
  id: string;
  scheduled_at: string;
  court_name: string;
  purpose: string;
  case_id: string;
  title: string;
  case_number: string;
  case_type: string;
  client_name: string;
  phone: string;
  email: string;
  address: string;
};

export type Notification = {
  id: string;
  title: string;
  body: string;
  read_at: string | null;
  created_at: string;
  case_id: string;
};

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, { credentials: "include", ...init });
  if (res.status === 204) return undefined as T;
  const text = await res.text();
  const body = text ? JSON.parse(text) : null;
  if (!res.ok) throw new ApiError(res.status, body?.error || "Request failed");
  return body as T;
}

function json(body: unknown): RequestInit {
  return { headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) };
}

export const api = {
  session: () => request<Lawyer>("/api/session"),
  signUp: (body: { name: string; email: string; password: string }) =>
    request<Lawyer>("/api/sign-up", { method: "POST", ...json(body) }),
  signIn: (body: { email: string; password: string }) =>
    request<Lawyer>("/api/session", { method: "POST", ...json(body) }),
  signOut: () => request<void>("/api/session", { method: "DELETE" }),
  clients: (q = "") => request<Client[]>(`/api/clients?q=${encodeURIComponent(q)}`),
  createClient: (body: Partial<Client> & { full_name: string }) =>
    request<Client>("/api/clients", { method: "POST", ...json(body) }),
  client: (id: string) => request<ClientDetail>(`/api/clients/${id}`),
  updateClient: (id: string, body: Partial<Client>) =>
    request<Client>(`/api/clients/${id}`, { method: "PATCH", ...json(body) }),
  createCase: (body: Record<string, string>) => request<CaseSummary>("/api/cases", { method: "POST", ...json(body) }),
  case: (id: string) => request<CaseDetail>(`/api/cases/${id}`),
  updateCase: (id: string, body: Record<string, string>) =>
    request<CaseSummary>(`/api/cases/${id}`, { method: "PATCH", ...json(body) }),
  upload: (caseId: string, file: File, label: string) => {
    const data = new FormData();
    data.set("file", file);
    data.set("label", label);
    return request<Document>(`/api/cases/${caseId}/documents`, { method: "POST", body: data });
  },
  download: (id: string) => request<{ url: string }>(`/api/documents/${id}`),
  deleteDocument: (id: string) => request<void>(`/api/documents/${id}`, { method: "DELETE" }),
  addHearing: (caseId: string, body: { scheduled_at: string; court_name: string; purpose: string }) =>
    request<Hearing & { reminders_planned: number }>(`/api/cases/${caseId}/hearings`, { method: "POST", ...json(body) }),
  completeHearing: (id: string, outcome: string) =>
    request<Hearing>(`/api/hearings/${id}/complete`, { method: "POST", ...json({ outcome }) }),
  cancelHearing: (id: string, outcome: string) =>
    request<Hearing>(`/api/hearings/${id}/cancel`, { method: "POST", ...json({ outcome }) }),
  adjournHearing: (id: string, body: { next_at: string; outcome: string; purpose?: string }) =>
    request<Hearing>(`/api/hearings/${id}/adjourn`, { method: "POST", ...json(body) }),
  diary: (from: string, to: string) =>
    request<DiaryEntry[]>(`/api/diary?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`),
  notifications: () => request<Notification[]>("/api/notifications"),
  markRead: (id: string) => request<void>(`/api/notifications/${id}/read`, { method: "POST" }),
};

export function fromLocalInput(value: string): string {
  return `${value}:00+05:30`;
}

function istDate(date: Date): string {
  return new Intl.DateTimeFormat("en-CA", { timeZone: "Asia/Kolkata", year: "numeric", month: "2-digit", day: "2-digit" }).format(date);
}

export function formatWhen(value: string): string {
  return new Intl.DateTimeFormat("en-IN", {
    timeZone: "Asia/Kolkata",
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "numeric",
    minute: "2-digit",
  }).format(new Date(value));
}

export function clock(value: string): string {
  return new Intl.DateTimeFormat("en-IN", {
    timeZone: "Asia/Kolkata",
    hour: "numeric",
    minute: "2-digit",
  }).format(new Date(value));
}

export function shortWhen(value: string): string {
  return new Intl.DateTimeFormat("en-IN", {
    timeZone: "Asia/Kolkata",
    weekday: "short",
    day: "numeric",
    month: "short",
    hour: "numeric",
    minute: "2-digit",
  }).format(new Date(value));
}

function calendarDate(key: string): Date {
  const [y, m, d] = key.split("-").map(Number);
  return new Date(Date.UTC(y, m - 1, d, 6, 30));
}

export function weekdayName(key: string): string {
  return new Intl.DateTimeFormat("en-IN", { timeZone: "Asia/Kolkata", weekday: "long" }).format(calendarDate(key));
}

export function monthDay(key: string): string {
  return new Intl.DateTimeFormat("en-IN", { timeZone: "Asia/Kolkata", day: "numeric", month: "long" }).format(calendarDate(key));
}

function weekStart(key: string): string {
  const [y, m, d] = key.split("-").map(Number);
  const date = new Date(Date.UTC(y, m - 1, d));
  const weekday = date.getUTCDay();
  const mondayOffset = weekday === 0 ? -6 : 1 - weekday;
  date.setUTCDate(date.getUTCDate() + mondayOffset);
  return date.toISOString().slice(0, 10);
}

export function weekBuckets<T extends { scheduled_at: string }>(items: T[]): { title: string; rows: T[] }[] {
  const today = istDate(new Date());
  const thisWeek = weekStart(today);
  const nextWeekDate = new Date(`${thisWeek}T00:00:00Z`);
  nextWeekDate.setUTCDate(nextWeekDate.getUTCDate() + 7);
  const nextWeek = nextWeekDate.toISOString().slice(0, 10);
  const groups = new Map<string, { title: string; rows: T[] }>();
  for (const item of items) {
    const start = weekStart(dayKey(item.scheduled_at));
    const title = start === thisWeek ? "This week" : start === nextWeek ? "Next week" : `Week of ${monthDay(start)}`;
    const bucket = groups.get(start) || { title, rows: [] };
    bucket.rows.push(item);
    groups.set(start, bucket);
  }
  return [...groups.entries()].sort(([a], [b]) => a.localeCompare(b)).map(([, bucket]) => bucket);
}

function rangeFromToday(addMonths: number, addDays: number): { from: string; to: string } {
  const today = istDate(new Date());
  const [y, m, d] = today.split("-").map(Number);
  const end = new Date(Date.UTC(y, m - 1 + addMonths, d + addDays)).toISOString().slice(0, 10);
  return { from: `${today}T00:00:00+05:30`, to: `${end}T00:00:00+05:30` };
}

export function diaryWindow(): { from: string; to: string } {
  return rangeFromToday(0, 7);
}

export function monthWindow(): { from: string; to: string } {
  return rangeFromToday(1, 0);
}

export function dayKey(value: string): string {
  return istDate(new Date(value));
}

export function dayLabel(key: string): string {
  const [y, m, d] = key.split("-").map(Number);
  return new Intl.DateTimeFormat("en-IN", {
    timeZone: "Asia/Kolkata",
    weekday: "long",
    day: "numeric",
    month: "long",
  }).format(new Date(Date.UTC(y, m - 1, d, 6, 30)));
}
