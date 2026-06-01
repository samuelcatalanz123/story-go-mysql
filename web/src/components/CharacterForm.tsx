import { useState } from "react";
import type { FormEvent } from "react";
import type { Organization } from "../types";
import { Field } from "../ui/Field";
import { Button } from "../ui/Button";
import styles from "./SceneForm.module.css";

export type CharacterFormValues = {
  title: string;
  text: string;
  organizationIds: number[];
};

type Props = {
  initial?: CharacterFormValues;
  organizationOptions: Organization[];
  onSubmit: (values: CharacterFormValues) => void;
  onCancel: () => void;
  submitting?: boolean;
  error?: string | null;
};

function selectedNumbers(select: HTMLSelectElement): number[] {
  return Array.from(select.selectedOptions).map((o) => Number(o.value));
}

export function CharacterForm({
  initial,
  organizationOptions,
  onSubmit,
  onCancel,
  submitting,
  error,
}: Props) {
  const [title, setTitle] = useState(initial?.title ?? "");
  const [text, setText] = useState(initial?.text ?? "");
  const [organizationIds, setOrganizationIds] = useState<number[]>(initial?.organizationIds ?? []);
  const [localError, setLocalError] = useState<string | null>(null);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (title.trim() === "") {
      setLocalError("El título es obligatorio");
      return;
    }
    setLocalError(null);
    onSubmit({ title: title.trim(), text, organizationIds });
  }

  const shownError = localError ?? error;

  return (
    <form onSubmit={handleSubmit}>
      {shownError && (
        <p role="alert" className={styles.alert}>
          {shownError}
        </p>
      )}
      <Field label="Título" htmlFor="cf-title">
        <input id="cf-title" value={title} onChange={(e) => setTitle(e.target.value)} />
      </Field>
      <Field label="Texto" htmlFor="cf-text">
        <textarea id="cf-text" value={text} onChange={(e) => setText(e.target.value)} />
      </Field>
      <Field label="Organizaciones" htmlFor="cf-organizations">
        <select
          id="cf-organizations"
          multiple
          value={organizationIds.map(String)}
          onChange={(e) => setOrganizationIds(selectedNumbers(e.target))}
        >
          {organizationOptions.map((o) => (
            <option key={o.id} value={o.id}>
              {o.title}
            </option>
          ))}
        </select>
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
