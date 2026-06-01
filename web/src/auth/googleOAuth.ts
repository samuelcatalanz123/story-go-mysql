// "Iniciar sesión con Google" usando OAuth 2.0 + PKCE.
//
// IMPORTANTE: necesita un Client ID de Google. Configúralo en web/.env como
// VITE_GOOGLE_CLIENT_ID=... (ver .env.example). Sin él, el botón avisa que
// falta configuración. El redirect_uri debe coincidir con el registrado en
// Google Cloud y con GOOGLE_REDIRECT_URI del backend.

const GOOGLE_AUTH_URL = "https://accounts.google.com/o/oauth2/v2/auth";
const VERIFIER_KEY = "pkce_code_verifier";

function base64url(bytes: Uint8Array): string {
  const str = btoa(String.fromCharCode(...bytes));
  return str.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

// randomVerifier crea el "code_verifier" de PKCE: una cadena aleatoria.
function randomVerifier(): string {
  const bytes = new Uint8Array(32);
  crypto.getRandomValues(bytes);
  return base64url(bytes);
}

// challenge = base64url(SHA-256(verifier)). Es lo que se manda a Google; el
// verifier original se guarda y se manda al backend al final (eso es PKCE).
async function challenge(verifier: string): Promise<string> {
  const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(verifier));
  return base64url(new Uint8Array(digest));
}

export function googleClientId(): string {
  const env = import.meta.env as Record<string, string | undefined>;
  return env.VITE_GOOGLE_CLIENT_ID ?? "";
}

export function takeCodeVerifier(): string | null {
  const v = sessionStorage.getItem(VERIFIER_KEY);
  sessionStorage.removeItem(VERIFIER_KEY);
  return v;
}

// startGoogleLogin genera PKCE, guarda el verifier y redirige a Google.
export async function startGoogleLogin(): Promise<void> {
  const clientId = googleClientId();
  if (!clientId) {
    alert("Falta configurar VITE_GOOGLE_CLIENT_ID (ver .env.example).");
    return;
  }
  const verifier = randomVerifier();
  sessionStorage.setItem(VERIFIER_KEY, verifier);

  const params = new URLSearchParams({
    client_id: clientId,
    redirect_uri: `${window.location.origin}/auth/oauth/callback`,
    response_type: "code",
    scope: "openid email profile",
    code_challenge: await challenge(verifier),
    code_challenge_method: "S256",
  });
  window.location.href = `${GOOGLE_AUTH_URL}?${params.toString()}`;
}
