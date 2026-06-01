import { createContext, useContext, useEffect, useState } from "react";
import type { ReactNode } from "react";
import { apiFetch, setAccessToken } from "../api/client";
import type { AuthResponse, User } from "../types";

type AuthState = {
  user: User | null;
  isAuthenticated: boolean;
  setSession: (res: AuthResponse) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  // Sólo el usuario se guarda en localStorage (para mostrarlo al instante).
  // El access token vive en memoria (en api/client) y el refresh token en una
  // cookie HttpOnly: ninguno de los dos es accesible desde JavaScript.
  const [user, setUser] = useState<User | null>(() => {
    const raw = localStorage.getItem("user");
    if (!raw) return null;
    try {
      return JSON.parse(raw) as User;
    } catch {
      return null;
    }
  });

  function setSession(res: AuthResponse) {
    setAccessToken(res.token);
    localStorage.setItem("user", JSON.stringify(res.user));
    setUser(res.user);
  }

  async function logout() {
    try {
      await apiFetch("/auth/logout", { method: "POST" });
    } catch {
      // Logout es "best effort": aunque falle la red, limpiamos localmente.
    }
    setAccessToken(null);
    localStorage.removeItem("user");
    setUser(null);
  }

  // Al cargar la app pedimos un access token nuevo usando la cookie de refresh.
  // Si la cookie es válida, recuperamos la sesión sin volver a iniciar sesión;
  // si no, limpiamos cualquier usuario obsoleto guardado.
  useEffect(() => {
    let active = true;
    apiFetch<AuthResponse>("/auth/refresh", { method: "POST" })
      .then((res) => {
        if (!active) return;
        setAccessToken(res.token);
        localStorage.setItem("user", JSON.stringify(res.user));
        setUser(res.user);
      })
      .catch(() => {
        if (!active) return;
        setAccessToken(null);
        localStorage.removeItem("user");
        setUser(null);
      });
    return () => {
      active = false;
    };
  }, []);

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: user !== null, setSession, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth debe usarse dentro de <AuthProvider>");
  return ctx;
}
