import { useState } from "react";
import type { FormEvent } from "react";
import { Link } from "react-router-dom";
import { forgotPassword } from "../api/auth";
import { Field } from "../ui/Field";
import { Button } from "../ui/Button";
import styles from "../components/AuthForm.module.css";

export function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      await forgotPassword(email.trim());
      setDone(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error desconocido");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.card}>
        <h1 className={styles.title}>Restablecer contraseña</h1>
        {done ? (
          // Mensaje vago a propósito: no revelamos si el email existe o no.
          <p>
            Si existe una cuenta con ese email, te enviamos un enlace para restablecer la
            contraseña. Revisa tu correo.
          </p>
        ) : (
          <form onSubmit={handleSubmit}>
            {error && (
              <p role="alert" className={styles.alert}>
                {error}
              </p>
            )}
            <Field label="Email" htmlFor="fp-email">
              <input
                id="fp-email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </Field>
            <Button type="submit" disabled={submitting}>
              Enviar enlace
            </Button>
          </form>
        )}
        <p className={styles.footer}>
          <Link to="/login">Volver a iniciar sesión</Link>
        </p>
      </div>
    </div>
  );
}
