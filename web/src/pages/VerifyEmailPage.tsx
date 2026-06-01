import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { verifyEmail } from "../api/auth";
import { EmptyState } from "../ui/EmptyState";
import { Button } from "../ui/Button";

// VerifyEmailPage es a donde lleva el enlace del correo de verificación. Lee el
// token del query string y lo confirma contra el backend.
export function VerifyEmailPage() {
  const [params] = useSearchParams();
  const [status, setStatus] = useState<"loading" | "ok" | "error">("loading");
  const [message, setMessage] = useState("");

  useEffect(() => {
    const token = params.get("token");
    if (!token) {
      setStatus("error");
      setMessage("Falta el token. Abre el enlace que te llegó por correo.");
      return;
    }
    let active = true;
    verifyEmail(token)
      .then(() => active && setStatus("ok"))
      .catch((e) => {
        if (!active) return;
        setStatus("error");
        setMessage(e instanceof Error ? e.message : "No se pudo verificar el correo");
      });
    return () => {
      active = false;
    };
  }, [params]);

  if (status === "loading") return <EmptyState title="Verificando tu correo…" message="Un momento." />;
  if (status === "ok") {
    return (
      <EmptyState
        title="¡Correo verificado! ✅"
        message="Tu cuenta ya está confirmada."
        action={
          <Link to="/characters">
            <Button>Ir a la app</Button>
          </Link>
        }
      />
    );
  }
  return (
    <EmptyState
      title="No se pudo verificar"
      message={message}
      action={
        <Link to="/login">
          <Button variant="secondary">Volver</Button>
        </Link>
      }
    />
  );
}
