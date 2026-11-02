import { FormEvent, useState } from "react";

export type Rule = {
  name: string;
  label: string;
  required?: boolean;
  email?: boolean;
  min?: number;
  file?: boolean;
};

export function checkFields(form: HTMLFormElement, rules: Rule[]): Record<string, string> {
  const data = new FormData(form);
  const out: Record<string, string> = {};
  for (const rule of rules) {
    if (rule.file) {
      const file = data.get(rule.name);
      if (!(file instanceof File) || file.size === 0) out[rule.name] = "Choose a file to upload";
      continue;
    }
    const value = String(data.get(rule.name) ?? "").trim();
    if (rule.required && value === "") {
      out[rule.name] = `${rule.label} is required`;
      continue;
    }
    if (value && rule.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
      out[rule.name] = "Enter a valid email";
      continue;
    }
    if (rule.min && value.length > 0 && value.length < rule.min) {
      out[rule.name] = `${rule.label} must be at least ${rule.min} characters`;
    }
  }
  return out;
}

export function FieldHint({ text }: { text?: string }) {
  if (!text) return null;
  return <p className="info">{text}</p>;
}

export function useFieldHints() {
  const [hints, setHints] = useState<Record<string, string>>({});
  function reject(form: HTMLFormElement, rules: Rule[]): boolean {
    const next = checkFields(form, rules);
    setHints(next);
    return Object.keys(next).length > 0;
  }
  function onInput(event: FormEvent<HTMLFormElement>) {
    const name = (event.target as { name?: string }).name;
    if (!name) return;
    setHints((current) => {
      if (!current[name]) return current;
      const next = { ...current };
      delete next[name];
      return next;
    });
  }
  function clear() {
    setHints({});
  }
  return { hints, reject, onInput, clear };
}
