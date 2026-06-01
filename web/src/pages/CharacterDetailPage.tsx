import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { getCharacter, uploadCharacterAvatar } from "../api/resources";
import { ApiError } from "../api/client";
import type { Character } from "../types";
import { Breadcrumbs } from "../components/Breadcrumbs";
import { AvatarUploader } from "../components/AvatarUploader";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { Button } from "../ui/Button";
import { useAuth } from "../auth/AuthContext";
import { NotFound } from "./NotFound";
import styles from "./DetailPage.module.css";

function formatDate(value: string): string {
  const d = new Date(value);
  return isNaN(d.getTime()) ? value : d.toLocaleString();
}

export function CharacterDetailPage() {
  // useParams() lee los valores dinámicos de la URL. Como la ruta es
  // "/characters/:id", aquí `id` es el texto que venga después de la barra.
  const { id } = useParams();
  const { isAuthenticated } = useAuth();

  // 4 estados para cubrir todo lo que puede pasar al pedir datos:
  const [character, setCharacter] = useState<Character | null>(null);
  const [loading, setLoading] = useState(true); // mientras esperamos la API
  const [notFound, setNotFound] = useState(false); // el id no existe (404)
  const [error, setError] = useState<string | null>(null); // otro error

  // useEffect se ejecuta cuando el componente aparece y cada vez que cambia
  // `id`. Aquí disparamos la petición GET /characters/:id al backend.
  useEffect(() => {
    // `active` evita actualizar el estado si el usuario ya se fue a otra ruta
    // antes de que la respuesta llegara (buena práctica para evitar warnings).
    let active = true;
    setLoading(true);
    setNotFound(false);
    setError(null);

    getCharacter(Number(id))
      .then((data) => {
        if (active) setCharacter(data);
      })
      .catch((e: unknown) => {
        if (!active) return;
        // Si la API responde 404, mostramos la página de "no encontrado".
        if (e instanceof ApiError && e.status === 404) {
          setNotFound(true);
        } else {
          setError(e instanceof Error ? e.message : "Error desconocido");
        }
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [id]);

  if (loading) return <SkeletonRows rows={3} cols={1} />;
  if (notFound) return <NotFound />;
  if (error || !character) {
    return (
      <EmptyState
        title="No se pudo cargar"
        message={error ?? "Inténtalo de nuevo."}
        action={
          <Link to="/characters">
            <Button variant="secondary">Volver a la lista</Button>
          </Link>
        }
      />
    );
  }

  // Si llegamos aquí, tenemos el personaje cargado: lo mostramos.
  return (
    <section>
      <Breadcrumbs
        items={[
          { label: "Personajes", to: "/characters" },
          { label: character.title },
        ]}
      />

      <PageHeader
        title={character.title}
        action={
          <Link to="/characters">
            <Button variant="secondary">← Volver</Button>
          </Link>
        }
      />

      <AvatarUploader
        avatarPath={character.avatarPath}
        canEdit={isAuthenticated}
        onUpload={(file) => uploadCharacterAvatar(character.id, file)}
      />

      <dl className={styles.fields}>
        <div className={styles.field}>
          <dt className={styles.label}>ID</dt>
          <dd className={styles.value}>{character.id}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Texto</dt>
          <dd className={styles.value}>{character.text ?? "—"}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Creado</dt>
          <dd className={styles.value}>{formatDate(character.createdAt)}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Actualizado</dt>
          <dd className={styles.value}>{formatDate(character.updatedAt)}</dd>
        </div>
      </dl>
    </section>
  );
}
