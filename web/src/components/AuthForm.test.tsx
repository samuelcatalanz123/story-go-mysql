import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { AuthForm } from "./AuthForm";

describe("AuthForm", () => {
  it("muestra error y no envía si los campos están vacíos", () => {
    const onSubmit = vi.fn();
    render(
      <AuthForm title="Entrar" submitLabel="Entrar" onSubmit={onSubmit} footer={<span>pie</span>} />,
    );
    fireEvent.click(screen.getByRole("button", { name: "Entrar" }));
    expect(screen.getByRole("alert")).toHaveTextContent("obligatorios");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("envía email y contraseña cuando son válidos", () => {
    const onSubmit = vi.fn();
    render(
      <AuthForm title="Entrar" submitLabel="Entrar" onSubmit={onSubmit} footer={<span>pie</span>} />,
    );
    fireEvent.change(screen.getByLabelText("Email"), { target: { value: "a@b.com" } });
    fireEvent.change(screen.getByLabelText("Contraseña"), { target: { value: "password123" } });
    fireEvent.click(screen.getByRole("button", { name: "Entrar" }));
    expect(onSubmit).toHaveBeenCalledWith("a@b.com", "password123");
  });
});
