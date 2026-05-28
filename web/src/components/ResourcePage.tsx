import { useState } from "react";
import { useList } from "../hooks/useList";
import { DataTable } from "./DataTable";
import type { Column } from "./DataTable";
import { ResourceForm } from "./ResourceForm";
import type { ResourceFormValues } from "./ResourceForm";

// Forma mínima que deben cumplir los elementos de esta página.
type ResourceItem = {
  id: number;
  title: string;
  text: string | null;
  updatedAt: string;
};

type RequestBody = { title: string; text: string | null };

type Props<T extends ResourceItem> = {
  heading: string;
  list: () => Promise<T[]>;
  create: (body: RequestBody) => Promise<T>;
  update: (id: number, body: RequestBody) => Promise<T>;
  remove: (id: number) => Promise<void>;
};

// `null` = sólo lista; "new" = creando; un objeto = editando ese elemento.
type Editing<T> = null | "new" | T;

export function ResourcePage<T extends ResourceItem>({
  heading,
  list,
  create,
  update,
  remove,
}: Props<T>) {
  const { data, loading, error, reload } = useList(list);
  const [editing, setEditing] = useState<Editing<T>>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const columns: Column<T>[] = [
    { header: "ID", render: (r) => r.id },
    { header: "Título", render: (r) => r.title },
    { header: "Texto", render: (r) => r.text ?? "—" },
    {
      header: "Actualizado",
      render: (r) => new Date(r.updatedAt).toLocaleString(),
    },
  ];

  function toBody(values: ResourceFormValues): RequestBody {
    return { title: values.title, text: values.text.trim() === "" ? null : values.text };
  }

  async function handleSubmit(values: ResourceFormValues) {
    setSubmitting(true);
    setFormError(null);
    try {
      if (editing === "new") {
        await create(toBody(values));
      } else if (editing) {
        await update(editing.id, toBody(values));
      }
      setEditing(null);
      reload();
    } catch (e: unknown) {
      setFormError(e instanceof Error ? e.message : "Error desconocido");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(row: T) {
    if (!window.confirm(`¿Borrar "${row.title}"?`)) return;
    try {
      await remove(row.id);
      reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : "Error desconocido");
    }
  }

  if (editing !== null) {
    const initial =
      editing === "new"
        ? undefined
        : { title: editing.title, text: editing.text ?? "" };
    return (
      <section>
        <h1>{editing === "new" ? `Nuevo: ${heading}` : `Editar: ${heading}`}</h1>
        <ResourceForm
          initial={initial}
          onSubmit={handleSubmit}
          onCancel={() => {
            setEditing(null);
            setFormError(null);
          }}
          submitting={submitting}
          error={formError}
        />
      </section>
    );
  }

  return (
    <section>
      <h1>{heading}</h1>
      <button onClick={() => setEditing("new")}>Nuevo</button>
      {loading && <p>Cargando…</p>}
      {error && <p role="alert">Error: {error}</p>}
      {!loading && !error && (
        <DataTable
          columns={columns}
          rows={data}
          onEdit={(row) => setEditing(row)}
          onDelete={handleDelete}
        />
      )}
    </section>
  );
}
