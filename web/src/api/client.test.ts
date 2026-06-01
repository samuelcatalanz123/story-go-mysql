import { describe, it, expect, vi, afterEach } from "vitest";
import { apiFetch, setAccessToken } from "./client";

afterEach(() => {
  vi.unstubAllGlobals();
  setAccessToken(null);
});

describe("apiFetch", () => {
  it("devuelve el JSON cuando la respuesta es exitosa", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ id: 1, title: "Asha" }),
      }),
    );
    const data = await apiFetch<{ id: number; title: string }>("/characters/1");
    expect(data).toEqual({ id: 1, title: "Asha" });
  });

  it("lanza el mensaje de error de la API cuando falla", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 409,
        json: async () => ({ error: "title already exists" }),
      }),
    );
    await expect(apiFetch("/characters")).rejects.toThrow("title already exists");
  });

  it("adjunta Authorization si hay access token en memoria", async () => {
    setAccessToken("abc123");
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({}),
    });
    vi.stubGlobal("fetch", fetchMock);
    await apiFetch("/characters");
    const opts = fetchMock.mock.calls[0][1];
    expect((opts.headers as Record<string, string>).Authorization).toBe("Bearer abc123");
  });

  it("ante un 401 intenta refrescar el token y reintenta una vez", async () => {
    const fetchMock = vi
      .fn()
      // 1) petición original → 401
      .mockResolvedValueOnce({ ok: false, status: 401, json: async () => ({}) })
      // 2) POST /auth/refresh → ok con token nuevo
      .mockResolvedValueOnce({ ok: true, status: 200, json: async () => ({ token: "nuevo" }) })
      // 3) reintento de la original → ok
      .mockResolvedValueOnce({ ok: true, status: 200, json: async () => ({ id: 1 }) });
    vi.stubGlobal("fetch", fetchMock);

    const data = await apiFetch<{ id: number }>("/characters/1");
    expect(data).toEqual({ id: 1 });
    expect(fetchMock).toHaveBeenCalledTimes(3);
  });
});
