import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { usePagedList } from "../hooks/usePagedList";
import { DataTable } from "./DataTable";
import type { Column } from "./DataTable";
import { ResourceForm } from "./ResourceForm";
import type { ResourceFormValues } from "./ResourceForm";
import { Modal } from "../ui/Modal";
import { ConfirmDialog } from "../ui/ConfirmDialog";
import { Button } from "../ui/Button";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { SearchBar } from "../ui/SearchBar";
import { Pagination } from "../ui/Pagination";
import { useToast } from "../ui/Toast";
import { useAuth } from "../auth/AuthContext";
import { ApiError } from "../api/client";
import type { Paged } from "../types";

type ResourceItem = {
  id: number;
  title: string;
  text: string | null;
  updatedAt: string;
};

type RequestBody = { title: string; text: string | null };
type ListArgs = { q: string; page: number; pageSize: number };

type Props<T extends ResourceItem> = {
  heading: string;
  list: (args: ListArgs) => Promise<Paged<T>>;
  create: (body: RequestBody) => Promise<T>;
  update: (id: number, body: RequestBody) => Promise<T>;
  remove: (id: number) => Promise<void>;
};

type Editing<T> = null | "new" | T;

function formatDate(value: string): string {
  const d = new Date(value);
  return isNaN(d.getTime()) ? value : d.toLocaleString();
}

export function ResourcePage<T extends ResourceItem>({
  heading,
  list,
  create,
  update,
  remove,
}: Props<T>) {
  const { data, total, page, pageSize, loading, error, setQuery, setPage, reload } =
    usePagedList(list);
  const toast = useToast();
  const { isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [editing, setEditing] = useState<Editing<T>>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState<T | null>(null);

  const columns: Column<T>[] = [
    { header: "ID", render: (r) => r.id },
    { header: "Título", render: (r) => r.title },
    { header: "Texto", render: (r) => r.text ?? "—" },
    { header: "Actualizado", render: (r) => formatDate(r.updatedAt) },
  ];

  function openNew() {
    setFormError(null);
    setEditing("new");
  }

  function isUnauthorized(e: unknown): boolean {
    if (e instanceof ApiError && e.status === 401) {
      logout();
      toast.error("Tu sesión expiró, inicia sesión de nuevo");
      navigate("/login");
      return true;
    }
    return false;
  }

  function toBody(values: ResourceFormValues): RequestBody {
    return { title: values.title, text: values.text.trim() === "" ? null : values.text };
  }

  async function handleSubmit(values: ResourceFormValues) {
    setSubmitting(true);
    setFormError(null);
    try {
      if (editing === "new") {
        await create(toBody(values));
        toast.success("Creado correctamente");
      } else if (editing) {
        await update(editing.id, toBody(values));
        toast.success("Cambios guardados");
      }
      setEditing(null);
      reload();
    } catch (e: unknown) {
      if (isUnauthorized(e)) return;
      const msg = e instanceof Error ? e.message : "Error desconocido";
      setFormError(msg);
      toast.error(msg);
    } finally {
      setSubmitting(false);
    }
  }

  async function confirmDelete() {
    if (!deleting) return;
    try {
      await remove(deleting.id);
      toast.success("Eliminado");
      reload();
    } catch (e: unknown) {
      if (!isUnauthorized(e)) {
        toast.error(e instanceof Error ? e.message : "Error desconocido");
      }
    } finally {
      setDeleting(null);
    }
  }

  const initial =
    editing && editing !== "new"
      ? { title: editing.title, text: editing.text ?? "" }
      : undefined;

  return (
    <section>
      <PageHeader
        title={heading}
        action={
          isAuthenticated ? (
            <Button onClick={openNew}>Nuevo</Button>
          ) : (
            <Link to="/login">Inicia sesión para gestionar</Link>
          )
        }
      />

      <SearchBar onQueryChange={setQuery} />

      {loading && <SkeletonRows rows={4} cols={4} />}
      {error && (
        <EmptyState
          title="No se pudo cargar"
          message={error}
          action={
            <Button variant="secondary" onClick={reload}>
              Reintentar
            </Button>
          }
        />
      )}
      {!loading && !error && data.length === 0 && (
        <EmptyState
          title="No hay resultados"
          message="Prueba con otra búsqueda o crea el primer elemento."
          action={isAuthenticated ? <Button onClick={openNew}>Nuevo</Button> : undefined}
        />
      )}
      {!loading && !error && data.length > 0 && (
        <>
          <DataTable
            columns={columns}
            rows={data}
            onEdit={
              isAuthenticated
                ? (row) => {
                    setFormError(null);
                    setEditing(row);
                  }
                : undefined
            }
            onDelete={isAuthenticated ? (row) => setDeleting(row) : undefined}
          />
          <Pagination page={page} pageSize={pageSize} total={total} onPage={setPage} />
        </>
      )}

      <Modal
        open={editing !== null}
        onClose={() => {
          setEditing(null);
          setFormError(null);
        }}
        title={editing === "new" ? `Nuevo: ${heading}` : `Editar: ${heading}`}
      >
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
      </Modal>

      <ConfirmDialog
        open={deleting !== null}
        title="Confirmar borrado"
        message={deleting ? `¿Seguro que quieres borrar "${deleting.title}"?` : ""}
        confirmLabel="Borrar"
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </section>
  );
}
