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

type Editing = null | "new" | Scene;

export function ScenesPage() {
  const scenes = useList(listScenes);
  const characters = useList(listCharacters);
  const locations = useList(listLocations);

  const [editing, setEditing] = useState<Editing>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const columns: Column<Scene>[] = [
    { header: "ID", render: (s) => s.id },
    { header: "Título", render: (s) => s.title },
    { header: "Inicio", render: (s) => s.startTimeline },
    { header: "Fin", render: (s) => s.endTimeline },
    {
      header: "Personajes",
      render: (s) => s.characters.map((c) => c.title).join(", ") || "—",
    },
    {
      header: "Lugares",
      render: (s) => s.locations.map((l) => l.title).join(", ") || "—",
    },
  ];

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
      } else if (editing) {
        await updateScene(editing.id, body);
      }
      setEditing(null);
      scenes.reload();
    } catch (e: unknown) {
      setFormError(e instanceof Error ? e.message : "Error desconocido");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(scene: Scene) {
    if (!window.confirm(`¿Borrar "${scene.title}"?`)) return;
    try {
      await deleteScene(scene.id);
      scenes.reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : "Error desconocido");
    }
  }

  if (editing !== null) {
    const initial: SceneFormValues | undefined =
      editing === "new"
        ? undefined
        : {
            title: editing.title,
            text: editing.text ?? "",
            startTimeline: editing.startTimeline,
            endTimeline: editing.endTimeline,
            characterIds: editing.characters.map((c) => c.id),
            locationIds: editing.locations.map((l) => l.id),
          };
    return (
      <section>
        <h1>{editing === "new" ? "Nueva escena" : "Editar escena"}</h1>
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
      </section>
    );
  }

  return (
    <section>
      <h1>Escenas</h1>
      <button onClick={() => setEditing("new")}>Nueva</button>
      {scenes.loading && <p>Cargando…</p>}
      {scenes.error && <p role="alert">Error: {scenes.error}</p>}
      {!scenes.loading && !scenes.error && (
        <DataTable
          columns={columns}
          rows={scenes.data}
          onEdit={(row) => setEditing(row)}
          onDelete={handleDelete}
        />
      )}
    </section>
  );
}
