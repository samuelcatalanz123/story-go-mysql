import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { Pagination } from "./Pagination";

describe("Pagination", () => {
  it("deshabilita Anterior en la primera página", () => {
    render(<Pagination page={1} pageSize={20} total={50} onPage={() => {}} />);
    expect(screen.getByRole("button", { name: "Anterior" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Siguiente" })).not.toBeDisabled();
  });

  it("deshabilita Siguiente en la última página", () => {
    render(<Pagination page={3} pageSize={20} total={50} onPage={() => {}} />);
    expect(screen.getByRole("button", { name: "Siguiente" })).toBeDisabled();
  });

  it("llama a onPage al avanzar", () => {
    const onPage = vi.fn();
    render(<Pagination page={1} pageSize={20} total={50} onPage={onPage} />);
    fireEvent.click(screen.getByRole("button", { name: "Siguiente" }));
    expect(onPage).toHaveBeenCalledWith(2);
  });
});
