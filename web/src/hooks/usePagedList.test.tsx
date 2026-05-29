import { renderHook, act, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { usePagedList } from "./usePagedList";

describe("usePagedList", () => {
  it("al buscar vuelve a la página 1", async () => {
    const loader = vi.fn().mockResolvedValue({ items: [], total: 0, page: 1, pageSize: 20 });
    const { result } = renderHook(() => usePagedList(loader));
    await waitFor(() => expect(result.current.loading).toBe(false));

    act(() => result.current.setPage(3));
    expect(result.current.page).toBe(3);

    act(() => result.current.setQuery("asha"));
    expect(result.current.page).toBe(1);
    expect(result.current.query).toBe("asha");
  });
});
