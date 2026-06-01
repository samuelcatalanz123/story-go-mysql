import { useState } from "react";
import type { ChangeEvent } from "react";
import { useToast } from "../ui/Toast";
import styles from "./AvatarUploader.module.css";

type Props = {
  avatarPath: string | null;
  // onUpload sube el archivo y devuelve la entidad con su nuevo avatarPath.
  onUpload: (file: File) => Promise<{ avatarPath: string | null }>;
  // canEdit controla si se muestra el control para subir (sólo logueado).
  canEdit: boolean;
};

export function AvatarUploader({ avatarPath, onUpload, canEdit }: Props) {
  const [path, setPath] = useState<string | null>(avatarPath);
  const [busy, setBusy] = useState(false);
  const toast = useToast();

  async function handleChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setBusy(true);
    try {
      const updated = await onUpload(file);
      setPath(updated.avatarPath);
      toast.success("Imagen subida");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Error al subir la imagen");
    } finally {
      setBusy(false);
      e.target.value = ""; // permite volver a elegir el mismo archivo
    }
  }

  return (
    <div className={styles.wrapper}>
      {path ? (
        <img className={styles.avatar} src={path} alt="Avatar" />
      ) : (
        <div className={styles.placeholder}>Sin imagen</div>
      )}
      {canEdit && (
        <label className={styles.upload}>
          {busy ? "Subiendo…" : "Cambiar imagen"}
          <input type="file" accept="image/*" onChange={handleChange} disabled={busy} hidden />
        </label>
      )}
    </div>
  );
}
