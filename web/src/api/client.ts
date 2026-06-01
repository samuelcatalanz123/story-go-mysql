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

// Hace una petición a la API y devuelve el JSON ya parseado. Adjunta el token
// JWT (si existe) en la cabecera Authorization. Lanza ApiError si la respuesta
// no es exitosa.
export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string> | undefined),
  };
  const token = localStorage.getItem("token");
  if (token) headers.Authorization = `Bearer ${token}`;

  const res = await fetch(`${BASE}${path}`, { ...options, headers });

  if (!res.ok) {
    let message = `Error ${res.status}`;
    try {
      const body = await res.json();
      if (body && typeof body.error === "string") message = body.error;
    } catch {
      // La respuesta de error no traía JSON; usamos el mensaje por defecto.
    }
    throw new ApiError(res.status, message);
  }

  // 204 No Content: no hay cuerpo que parsear.
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

// apiUpload envía un FormData (archivos) por POST. A diferencia de apiFetch,
// NO fija Content-Type: el navegador lo pone solo con el "boundary" correcto
// de multipart/form-data. Adjunta el token JWT si existe.
export async function apiUpload<T>(path: string, formData: FormData): Promise<T> {
  const headers: Record<string, string> = {};
  const token = localStorage.getItem("token");
  if (token) headers.Authorization = `Bearer ${token}`;

  const res = await fetch(`${BASE}${path}`, { method: "POST", headers, body: formData });

  if (!res.ok) {
    let message = `Error ${res.status}`;
    try {
      const body = await res.json();
      if (body && typeof body.error === "string") message = body.error;
    } catch {
      // El cuerpo de error no traía JSON; usamos el mensaje por defecto.
    }
    throw new ApiError(res.status, message);
  }

  return (await res.json()) as T;
}
