import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { ResourceForm } from "./ResourceForm";

describe("ResourceForm", () => {
  it("muestra error y no envía si el título está vacío", () => {
    const onSubmit = vi.fn();
    render(<ResourceForm onSubmit={onSubmit} onCancel={() => {}} />);
    fireEvent.click(screen.getByText("Guardar"));
    expect(screen.getByRole("alert")).toHaveTextContent("El título es obligatorio");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("envía título y texto cuando son válidos", () => {
    const onSubmit = vi.fn();
    render(<ResourceForm onSubmit={onSubmit} onCancel={() => {}} />);
    fireEvent.change(screen.getByLabelText("Título"), {
      target: { value: "Asha" },
    });
    fireEvent.change(screen.getByLabelText("Texto"), {
      target: { value: "Una piloto" },
    });
    fireEvent.click(screen.getByText("Guardar"));
    expect(onSubmit).toHaveBeenCalledWith({ title: "Asha", text: "Una piloto" });
  });
});
