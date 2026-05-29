import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { SearchBar } from "./SearchBar";

describe("SearchBar", () => {
  it("llama a onQueryChange con el texto tras el debounce", async () => {
    const onQueryChange = vi.fn();
    render(<SearchBar onQueryChange={onQueryChange} />);
    fireEvent.change(screen.getByLabelText("Buscar"), { target: { value: "asha" } });
    await waitFor(() => expect(onQueryChange).toHaveBeenCalledWith("asha"));
  });
});
