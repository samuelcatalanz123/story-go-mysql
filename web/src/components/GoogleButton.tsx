import { Button } from "../ui/Button";
import { startGoogleLogin } from "../auth/googleOAuth";

// GoogleButton inicia el flujo "Iniciar sesión con Google" (OAuth + PKCE).
export function GoogleButton() {
  return (
    <Button type="button" variant="secondary" onClick={() => void startGoogleLogin()}>
      Iniciar sesión con Google
    </Button>
  );
}
