import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { Button } from "./Button";

describe("Button", () => {
  it("aplica la clase de la variante indicada", () => {
    render(<Button variant="danger">Borrar</Button>);
    expect(screen.getByRole("button", { name: "Borrar" }).className).toMatch(/danger/);
  });

  it("respeta el atributo disabled", () => {
    render(<Button disabled>Guardar</Button>);
    expect(screen.getByRole("button", { name: "Guardar" })).toBeDisabled();
  });
});
