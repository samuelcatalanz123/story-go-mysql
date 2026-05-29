import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { Modal } from "./Modal";

describe("Modal", () => {
  it("no renderiza nada si open=false", () => {
    render(
      <Modal open={false} onClose={() => {}} title="T">
        contenido
      </Modal>,
    );
    expect(screen.queryByText("contenido")).not.toBeInTheDocument();
  });

  it("llama a onClose al pulsar Escape", () => {
    const onClose = vi.fn();
    render(
      <Modal open onClose={onClose} title="T">
        contenido
      </Modal>,
    );
    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).toHaveBeenCalled();
  });

  it("llama a onClose al hacer clic en el overlay", () => {
    const onClose = vi.fn();
    render(
      <Modal open onClose={onClose} title="T">
        contenido
      </Modal>,
    );
    fireEvent.click(screen.getByTestId("modal-overlay"));
    expect(onClose).toHaveBeenCalled();
  });
});
