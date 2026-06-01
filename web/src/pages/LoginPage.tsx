import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { AuthForm } from "../components/AuthForm";
import { GoogleButton } from "../components/GoogleButton";
import { login } from "../api/auth";
import { useAuth } from "../auth/AuthContext";
import { useToast } from "../ui/Toast";

export function LoginPage() {
  const { setSession } = useAuth();
  const navigate = useNavigate();
  const toast = useToast();
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(email: string, password: string) {
    setSubmitting(true);
    setError(null);
    try {
      const res = await login({ email, password });
      setSession(res);
      toast.success("Sesión iniciada");
      navigate("/characters");
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Error desconocido");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <AuthForm
      title="Iniciar sesión"
      submitLabel="Entrar"
      onSubmit={handleSubmit}
      submitting={submitting}
      error={error}
      footer={
        <>
          <div style={{ marginBottom: "var(--space-4)" }}>
            <GoogleButton />
          </div>
          <div style={{ marginBottom: "var(--space-2)" }}>
            <Link to="/forgot-password">¿Olvidaste tu contraseña?</Link>
          </div>
          ¿No tienes cuenta? <Link to="/register">Regístrate</Link>
        </>
      }
    />
  );
}
