import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { oauthGoogle } from "../api/auth";
import { takeCodeVerifier } from "../auth/googleOAuth";
import { useAuth } from "../auth/AuthContext";
import { EmptyState } from "../ui/EmptyState";
import { Button } from "../ui/Button";

// OAuthCallbackPage es a donde Google redirige tras el consentimiento, con un
// ?code=... en la URL. Aquí lo cambiamos (junto al code_verifier de PKCE) por
// nuestros propios tokens y guardamos la sesión.
export function OAuthCallbackPage() {
  const [params] = useSearchParams();
  const { setSession } = useAuth();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    const code = params.get("code");
    const verifier = takeCodeVerifier();

    if (!code || !verifier) {
      setError("Faltan datos del callback de Google (code o verifier).");
      return;
    }

    oauthGoogle(code, verifier)
      .then((res) => {
        if (!active) return;
        setSession(res);
        navigate("/characters", { replace: true });
      })
      .catch((e: unknown) => {
        if (active) setError(e instanceof Error ? e.message : "No se pudo iniciar sesión con Google");
      });

    return () => {
      active = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (error) {
    return (
      <EmptyState
        title="Error al iniciar sesión con Google"
        message={error}
        action={<Button onClick={() => navigate("/login")}>Volver a iniciar sesión</Button>}
      />
    );
  }
  return <EmptyState title="Conectando con Google…" message="Un momento por favor." />;
}
