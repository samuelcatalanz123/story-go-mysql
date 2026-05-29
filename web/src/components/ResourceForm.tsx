import { useState } from "react";
import type { FormEvent } from "react";
import { Field } from "../ui/Field";
import { Button } from "../ui/Button";
import styles from "./ResourceForm.module.css";

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
      {shownError && (
        <p role="alert" className={styles.alert}>
          {shownError}
        </p>
      )}
      <Field label="Título" htmlFor="rf-title">
        <input id="rf-title" value={title} onChange={(e) => setTitle(e.target.value)} />
      </Field>
      <Field label="Texto" htmlFor="rf-text">
        <textarea id="rf-text" value={text} onChange={(e) => setText(e.target.value)} />
      </Field>
      <div className={styles.actions}>
        <Button type="button" variant="secondary" onClick={onCancel}>
          Cancelar
        </Button>
        <Button type="submit" disabled={submitting}>
          Guardar
        </Button>
      </div>
    </form>
  );
}
