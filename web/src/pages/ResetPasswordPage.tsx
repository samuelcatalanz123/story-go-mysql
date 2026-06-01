import { useState } from "react";
import type { FormEvent } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { resetPassword } from "../api/auth";
import { Field } from "../ui/Field";
import { Button } from "../ui/Button";
import { useToast } from "../ui/Toast";
import styles from "../components/AuthForm.module.css";

export function ResetPasswordPage() {
  const [params] = useSearchParams();
  const token = params.get("token") ?? "";
  const navigate = useNavigate();
  const toast = useToast();
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (password.length < 8) {
      setError("La contraseña debe tener al menos 8 caracteres");
      return;
    }
    setSubmitting(true);
    setError(null);
    try {
      await resetPassword(token, password);
      toast.success("Contraseña actualizada, inicia sesión");
      navigate("/login");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error desconocido");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.card}>
        <h1 className={styles.title}>Nueva contraseña</h1>
        {!token ? (
          <p role="alert" className={styles.alert}>
            Falta el token. Abre el enlace que te llegó por correo.
          </p>
        ) : (
          <form onSubmit={handleSubmit}>
            {error && (
              <p role="alert" className={styles.alert}>
                {error}
              </p>
            )}
            <Field label="Nueva contraseña" htmlFor="rp-password">
              <input
                id="rp-password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </Field>
            <Button type="submit" disabled={submitting}>
              Guardar contraseña
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
