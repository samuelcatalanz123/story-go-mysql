import { describe, it, expect, vi, afterEach } from "vitest";
import { apiFetch } from "./client";

afterEach(() => vi.unstubAllGlobals());

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

  it("adjunta Authorization si hay token guardado", async () => {
    localStorage.setItem("token", "abc123");
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({}),
    });
    vi.stubGlobal("fetch", fetchMock);
    await apiFetch("/characters");
    const opts = fetchMock.mock.calls[0][1];
    expect((opts.headers as Record<string, string>).Authorization).toBe("Bearer abc123");
    localStorage.removeItem("token");
  });
});
