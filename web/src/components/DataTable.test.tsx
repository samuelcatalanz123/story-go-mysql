import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { DataTable } from "./DataTable";

describe("DataTable", () => {
  it("renderiza una fila por cada elemento", () => {
    render(
      <DataTable
        columns={[{ header: "Título", render: (r) => r.title }]}
        rows={[
          { id: 1, title: "Asha" },
          { id: 2, title: "Bo" },
        ]}
      />,
    );
    expect(screen.getByText("Asha")).toBeInTheDocument();
    expect(screen.getByText("Bo")).toBeInTheDocument();
  });

  it("muestra un mensaje cuando no hay elementos", () => {
    render(<DataTable columns={[]} rows={[]} />);
    expect(screen.getByText("No hay elementos todavía.")).toBeInTheDocument();
  });
});
