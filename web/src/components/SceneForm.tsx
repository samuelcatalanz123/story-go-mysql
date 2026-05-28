import { useState } from "react";
import type { FormEvent } from "react";
import type { Character, Location } from "../types";

export type SceneFormValues = {
  title: string;
  text: string;
  startTimeline: number;
  endTimeline: number;
  characterIds: number[];
  locationIds: number[];
};

type Props = {
  initial?: SceneFormValues;
  characterOptions: Character[];
  locationOptions: Location[];
  onSubmit: (values: SceneFormValues) => void;
  onCancel: () => void;
  submitting?: boolean;
  error?: string | null;
};

// Lee los valores seleccionados de un <select multiple> como números.
function selectedNumbers(select: HTMLSelectElement): number[] {
  return Array.from(select.selectedOptions).map((o) => Number(o.value));
}

export function SceneForm({
  initial,
  characterOptions,
  locationOptions,
  onSubmit,
  onCancel,
  submitting,
  error,
}: Props) {
  const [title, setTitle] = useState(initial?.title ?? "");
  const [text, setText] = useState(initial?.text ?? "");
  const [start, setStart] = useState(String(initial?.startTimeline ?? 0));
  const [end, setEnd] = useState(String(initial?.endTimeline ?? 0));
  const [characterIds, setCharacterIds] = useState<number[]>(
    initial?.characterIds ?? [],
  );
  const [locationIds, setLocationIds] = useState<number[]>(
    initial?.locationIds ?? [],
  );
  const [localError, setLocalError] = useState<string | null>(null);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (title.trim() === "") {
      setLocalError("El título es obligatorio");
      return;
    }
    setLocalError(null);
    onSubmit({
      title: title.trim(),
      text,
      startTimeline: Number(start),
      endTimeline: Number(end),
      characterIds,
      locationIds,
    });
  }

  const shownError = localError ?? error;

  return (
    <form onSubmit={handleSubmit}>
      {shownError && <p role="alert">{shownError}</p>}
      <div>
        <label htmlFor="sf-title">Título</label>
        <input id="sf-title" value={title} onChange={(e) => setTitle(e.target.value)} />
      </div>
      <div>
        <label htmlFor="sf-text">Texto</label>
        <textarea id="sf-text" value={text} onChange={(e) => setText(e.target.value)} />
      </div>
      <div>
        <label htmlFor="sf-start">Inicio (timeline)</label>
        <input
          id="sf-start"
          type="number"
          value={start}
          onChange={(e) => setStart(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="sf-end">Fin (timeline)</label>
        <input
          id="sf-end"
          type="number"
          value={end}
          onChange={(e) => setEnd(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="sf-characters">Personajes</label>
        <select
          id="sf-characters"
          multiple
          value={characterIds.map(String)}
          onChange={(e) => setCharacterIds(selectedNumbers(e.target))}
        >
          {characterOptions.map((c) => (
            <option key={c.id} value={c.id}>
              {c.title}
            </option>
          ))}
        </select>
      </div>
      <div>
        <label htmlFor="sf-locations">Lugares</label>
        <select
          id="sf-locations"
          multiple
          value={locationIds.map(String)}
          onChange={(e) => setLocationIds(selectedNumbers(e.target))}
        >
          {locationOptions.map((l) => (
            <option key={l.id} value={l.id}>
              {l.title}
            </option>
          ))}
        </select>
      </div>
      <button type="submit" disabled={submitting}>
        Guardar
      </button>
      <button type="button" onClick={onCancel}>
        Cancelar
      </button>
    </form>
  );
}
