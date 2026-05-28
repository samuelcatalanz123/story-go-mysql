import { useState } from "react";
import type { FormEvent } from "react";

export type ResourceFormValues = { title: string; text: string };

type Props = {
  initial?: ResourceFormValues;
  onSubmit: (values: ResourceFormValues) => void;
  onCancel: () => void;
  submitting?: boolean;
  error?: string | null;
};

export function ResourceForm({ initial, onSubmit, onCancel, submitting, error }: Props) {
  const [title, setTitle] = useState(initial?.title ?? "");
  const [text, setText] = useState(initial?.text ?? "");
  const [localError, setLocalError] = useState<string | null>(null);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (title.trim() === "") {
      setLocalError("El título es obligatorio");
      return;
    }
    setLocalError(null);
    onSubmit({ title: title.trim(), text });
  }

  const shownError = localError ?? error;

  return (
    <form onSubmit={handleSubmit}>
      {shownError && <p role="alert">{shownError}</p>}
      <div>
        <label htmlFor="rf-title">Título</label>
        <input
          id="rf-title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="rf-text">Texto</label>
        <textarea
          id="rf-text"
          value={text}
          onChange={(e) => setText(e.target.value)}
        />
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
