import { useState } from "react";
import { useList } from "../hooks/useList";
import { DataTable } from "../components/DataTable";
import type { Column } from "../components/DataTable";
import { SceneForm } from "../components/SceneForm";
import type { SceneFormValues } from "../components/SceneForm";
import type { Scene } from "../types";
import {
  listScenes,
  createScene,
  updateScene,
  deleteScene,
  listCharacters,
  listLocations,
} from "../api/resources";
import { Modal } from "../ui/Modal";
import { ConfirmDialog } from "../ui/ConfirmDialog";
import { Button } from "../ui/Button";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { Badge } from "../ui/Badge";
import { useToast } from "../ui/Toast";

type Editing = null | "new" | Scene;

export function ScenesPage() {
  const scenes = useList(listScenes);
  const characters = useList(listCharacters);
  const locations = useList(listLocations);
  const toast = useToast();

  const [editing, setEditing] = useState<Editing>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState<Scene | null>(null);

  const columns: Column<Scene>[] = [
    { header: "ID", render: (s) => s.id },
    { header: "Título", render: (s) => s.title },
    { header: "Inicio", render: (s) => s.startTimeline },
    { header: "Fin", render: (s) => s.endTimeline },
    {
      header: "Personajes",
      render: (s) =>
        s.characters.length
          ? s.characters.map((c) => (
              <Badge key={c.id} tone="primary">
                {c.title}
              </Badge>
            ))
          : "—",
    },
    {
      header: "Lugares",
      render: (s) =>
        s.locations.length
          ? s.locations.map((l) => <Badge key={l.id}>{l.title}</Badge>)
          : "—",
    },
  ];

  function openNew() {
    setFormError(null);
    setEditing("new");
  }

  async function handleSubmit(values: SceneFormValues) {
    setSubmitting(true);
    setFormError(null);
    const body = {
      title: values.title,
      text: values.text.trim() === "" ? null : values.text,
      startTimeline: values.startTimeline,
      endTimeline: values.endTimeline,
      characterIds: values.characterIds,
      locationIds: values.locationIds,
    };
    try {
      if (editing === "new") {
        await createScene(body);
        toast.success("Escena creada");
      } else if (editing) {
        await updateScene(editing.id, body);
        toast.success("Cambios guardados");
      }
      setEditing(null);
      scenes.reload();
    } catch (e: unknown) {
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
      await deleteScene(deleting.id);
      toast.success("Escena eliminada");
      scenes.reload();
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : "Error desconocido");
    } finally {
      setDeleting(null);
    }
  }

  const initial: SceneFormValues | undefined =
    editing && editing !== "new"
      ? {
          title: editing.title,
          text: editing.text ?? "",
          startTimeline: editing.startTimeline,
          endTimeline: editing.endTimeline,
          characterIds: editing.characters.map((c) => c.id),
          locationIds: editing.locations.map((l) => l.id),
        }
      : undefined;

  return (
    <section>
      <PageHeader title="Escenas" action={<Button onClick={openNew}>Nueva</Button>} />

      {scenes.loading && <SkeletonRows rows={4} cols={6} />}
      {scenes.error && (
        <EmptyState
          title="No se pudo cargar"
          message={scenes.error}
          action={
            <Button variant="secondary" onClick={scenes.reload}>
              Reintentar
            </Button>
          }
        />
      )}
      {!scenes.loading && !scenes.error && scenes.data.length === 0 && (
        <EmptyState
          title="Aún no hay escenas"
          message="Crea la primera escena para empezar."
          action={<Button onClick={openNew}>Nueva</Button>}
        />
      )}
      {!scenes.loading && !scenes.error && scenes.data.length > 0 && (
        <DataTable
          columns={columns}
          rows={scenes.data}
          onEdit={(row) => {
            setFormError(null);
            setEditing(row);
          }}
          onDelete={(row) => setDeleting(row)}
        />
      )}

      <Modal
        open={editing !== null}
        onClose={() => {
          setEditing(null);
          setFormError(null);
        }}
        title={editing === "new" ? "Nueva escena" : "Editar escena"}
      >
        <SceneForm
          initial={initial}
          characterOptions={characters.data}
          locationOptions={locations.data}
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
