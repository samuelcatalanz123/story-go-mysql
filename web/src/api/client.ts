const BASE = "/api";

// Hace una petición a la API y devuelve el JSON ya parseado.
// Si la respuesta no es exitosa, lanza un Error con el mensaje de la API.
export async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!res.ok) {
    let message = `Error ${res.status}`;
    try {
      const body = await res.json();
      if (body && typeof body.error === "string") message = body.error;
    } catch {
      // La respuesta de error no traía JSON; usamos el mensaje por defecto.
    }
    throw new Error(message);
  }

  // 204 No Content: no hay cuerpo que parsear. (En esta API el DELETE
  // responde 200 con un pequeño cuerpo JSON, que simplemente devolvemos.)
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}
