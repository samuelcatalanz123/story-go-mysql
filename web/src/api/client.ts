const BASE = "/api";

// ApiError lleva el código HTTP para poder distinguir, por ejemplo, un 401.
export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

// El access token vive SÓLO en memoria (no en localStorage), así un ataque XSS
// no puede robarlo. El refresh token va en una cookie HttpOnly que el navegador
// envía solo (no la toca JavaScript).
let accessToken: string | null = null;

// setAccessToken guarda el access token en memoria. Lo llama el AuthContext
// tras login o al refrescar la sesión.
export function setAccessToken(token: string | null): void {
  accessToken = token;
}

// buildError lee el mensaje de error JSON de una respuesta fallida.
async function buildError(res: Response): Promise<ApiError> {
  let message = `Error ${res.status}`;
  try {
    const body = await res.json();
    if (body && typeof body.error === "string") message = body.error;
  } catch {
    // La respuesta de error no traía JSON; usamos el mensaje por defecto.
  }
  return new ApiError(res.status, message);
}

// rawFetch hace la petición adjuntando el access token (si hay) y enviando las
// cookies (credentials: "include") para que viaje el refresh token.
function rawFetch(path: string, options: RequestInit): Promise<Response> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string> | undefined),
  };
  if (accessToken) headers.Authorization = `Bearer ${accessToken}`;
  return fetch(`${BASE}${path}`, { ...options, headers, credentials: "include" });
}

// tryRefresh pide un nuevo access token usando la cookie de refresh. Devuelve
// true si lo consiguió. No usa apiFetch para evitar un bucle de reintentos.
async function tryRefresh(): Promise<boolean> {
  const res = await fetch(`${BASE}/auth/refresh`, { method: "POST", credentials: "include" });
  if (!res.ok) return false;
  try {
    const data = (await res.json()) as { token?: string };
    if (data.token) {
      accessToken = data.token;
      return true;
    }
  } catch {
    // respuesta sin JSON válido
  }
  return false;
}

// noRetry son las rutas de auth: si fallan con 401 NO intentamos refrescar
// (sería un bucle, y un login fallido debe propagarse tal cual).
const noRetry = new Set(["/auth/refresh", "/auth/login", "/auth/register", "/auth/logout"]);

// apiFetch hace una petición JSON. Si recibe 401, intenta refrescar el access
// token UNA vez con la cookie de refresh y reintenta la petición original.
export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  let res = await rawFetch(path, options);

  if (res.status === 401 && !noRetry.has(path) && (await tryRefresh())) {
    res = await rawFetch(path, options);
  }

  if (!res.ok) throw await buildError(res);

  // 204 No Content: no hay cuerpo que parsear.
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

// apiUpload envía un FormData (archivos) por POST. A diferencia de apiFetch, NO
// fija Content-Type: el navegador lo pone solo con el "boundary" de
// multipart/form-data. Reintenta una vez si el access token expiró (401).
export async function apiUpload<T>(path: string, formData: FormData): Promise<T> {
  const send = () => {
    const headers: Record<string, string> = {};
    if (accessToken) headers.Authorization = `Bearer ${accessToken}`;
    return fetch(`${BASE}${path}`, { method: "POST", headers, body: formData, credentials: "include" });
  };

  let res = await send();
  if (res.status === 401 && (await tryRefresh())) {
    res = await send();
  }

  if (!res.ok) throw await buildError(res);
  return (await res.json()) as T;
}
