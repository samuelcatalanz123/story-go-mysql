import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useList } from "../hooks/useList";
import { usePagedList } from "../hooks/usePagedList";
import { DataTable } from "../components/DataTable";
import type { Column } from "../components/DataTable";
import { CharacterForm } from "../components/CharacterForm";
import type { CharacterFormValues } from "../components/CharacterForm";
import type { Character } from "../types";
import {
  listCharacters,
  createCharacter,
  updateCharacter,
  deleteCharacter,
  listAllOrganizations,
} from "../api/resources";
import { Modal } from "../ui/Modal";
import { ConfirmDialog } from "../ui/ConfirmDialog";
import { Button } from "../ui/Button";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { Badge } from "../ui/Badge";
import { SearchBar } from "../ui/SearchBar";
import { Pagination } from "../ui/Pagination";
import { useToast } from "../ui/Toast";
import { useAuth } from "../auth/AuthContext";
import { ApiError } from "../api/client";

type Editing = null | "new" | Character;

function formatDate(value: string): string {
  const d = new Date(value);
  return isNaN(d.getTime()) ? value : d.toLocaleString();
}

export function CharactersPage() {
  const characters = usePagedList(listCharacters);
  const organizations = useList(listAllOrganizations);
  const toast = useToast();
  const { isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const [editing, setEditing] = useState<Editing>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState<Character | null>(null);

  const columns: Column<Character>[] = [
    { header: "ID", render: (c) => c.id },
    {
      header: "Avatar",
      render: (c) =>
        c.avatarPath ? (
          <img
            src={c.avatarPath}
            alt={c.title}
            style={{ width: 32, height: 32, objectFit: "cover", borderRadius: 6 }}
          />
        ) : (
          "—"
        ),
    },
    { header: "Título", render: (c) => <Link to={`/characters/${c.id}`}>{c.title}</Link> },
    { header: "Texto", render: (c) => c.text ?? "—" },
    {
      header: "Organizaciones",
      render: (c) =>
        c.organizations && c.organizations.length
          ? c.organizations.map((o) => (
              <Badge key={o.id} tone="primary">
                {o.title}
              </Badge>
            ))
          : "—",
    },
    { header: "Actualizado", render: (c) => formatDate(c.updatedAt) },
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

  async function handleSubmit(values: CharacterFormValues) {
    setSubmitting(true);
    setFormError(null);
    const body = {
      title: values.title,
      text: values.text.trim() === "" ? null : values.text,
      organizationIds: values.organizationIds,
    };
    try {
      if (editing === "new") {
        await createCharacter(body);
        toast.success("Personaje creado");
      } else if (editing) {
        await updateCharacter(editing.id, body);
        toast.success("Cambios guardados");
      }
      setEditing(null);
      characters.reload();
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
      await deleteCharacter(deleting.id);
      toast.success("Personaje eliminado");
      characters.reload();
    } catch (e: unknown) {
      if (!isUnauthorized(e)) {
        toast.error(e instanceof Error ? e.message : "Error desconocido");
      }
    } finally {
      setDeleting(null);
    }
  }

  const initial: CharacterFormValues | undefined =
    editing && editing !== "new"
      ? {
          title: editing.title,
          text: editing.text ?? "",
          organizationIds: (editing.organizations ?? []).map((o) => o.id),
        }
      : undefined;

  return (
    <section>
      <PageHeader
        title="Personajes"
        action={
          isAuthenticated ? (
            <Button onClick={openNew}>Nuevo</Button>
          ) : (
            <Link to="/login">Inicia sesión para gestionar</Link>
          )
        }
      />

      <SearchBar onQueryChange={characters.setQuery} />

      {characters.loading && <SkeletonRows rows={4} cols={5} />}
      {characters.error && (
        <EmptyState
          title="No se pudo cargar"
          message={characters.error}
          action={
            <Button variant="secondary" onClick={characters.reload}>
              Reintentar
            </Button>
          }
        />
      )}
      {!characters.loading && !characters.error && characters.data.length === 0 && (
        <EmptyState
          title="No hay resultados"
          message="Prueba con otra búsqueda o crea el primer personaje."
          action={isAuthenticated ? <Button onClick={openNew}>Nuevo</Button> : undefined}
        />
      )}
      {!characters.loading && !characters.error && characters.data.length > 0 && (
        <>
          <DataTable
            columns={columns}
            rows={characters.data}
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
          <Pagination
            page={characters.page}
            pageSize={characters.pageSize}
            total={characters.total}
            onPage={characters.setPage}
          />
        </>
      )}

      <Modal
        open={editing !== null}
        onClose={() => {
          setEditing(null);
          setFormError(null);
        }}
        title={editing === "new" ? "Nuevo personaje" : "Editar personaje"}
      >
        <CharacterForm
          initial={initial}
          organizationOptions={organizations.data}
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
