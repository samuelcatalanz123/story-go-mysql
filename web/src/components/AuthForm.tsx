import { useState } from "react";
import type { FormEvent, ReactNode } from "react";
import { Field } from "../ui/Field";
import { Button } from "../ui/Button";
import styles from "./AuthForm.module.css";

type Props = {
  title: string;
  submitLabel: string;
  onSubmit: (email: string, password: string) => void;
  submitting?: boolean;
  error?: string | null;
  footer: ReactNode;
};

export function AuthForm({ title, submitLabel, onSubmit, submitting, error, footer }: Props) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [localError, setLocalError] = useState<string | null>(null);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (email.trim() === "" || password === "") {
      setLocalError("Email y contraseña son obligatorios");
      return;
    }
    setLocalError(null);
    onSubmit(email.trim(), password);
  }

  const shownError = localError ?? error;

  return (
    <div className={styles.wrapper}>
      <form className={styles.card} onSubmit={handleSubmit}>
        <h1 className={styles.title}>{title}</h1>
        {shownError && (
          <p role="alert" className={styles.alert}>
            {shownError}
          </p>
        )}
        <Field label="Email" htmlFor="auth-email">
          <input
            id="auth-email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </Field>
        <Field label="Contraseña" htmlFor="auth-password">
          <input
            id="auth-password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </Field>
        <Button type="submit" disabled={submitting}>
          {submitLabel}
        </Button>
        <p className={styles.footer}>{footer}</p>
      </form>
    </div>
  );
}
