import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { getScene } from "../api/resources";
import { ApiError } from "../api/client";
import type { Scene } from "../types";
import { Breadcrumbs } from "../components/Breadcrumbs";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { Badge } from "../ui/Badge";
import { Button } from "../ui/Button";
import { NotFound } from "./NotFound";
import styles from "./DetailPage.module.css";

function formatDate(value: string): string {
  const d = new Date(value);
  return isNaN(d.getTime()) ? value : d.toLocaleString();
}

export function SceneDetailPage() {
  const { id } = useParams();
  const [scene, setScene] = useState<Scene | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    setLoading(true);
    setNotFound(false);
    setError(null);

    getScene(Number(id))
      .then((data) => {
        if (active) setScene(data);
      })
      .catch((e: unknown) => {
        if (!active) return;
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

  if (loading) return <SkeletonRows rows={4} cols={1} />;
  if (notFound) return <NotFound />;
  if (error || !scene) {
    return (
      <EmptyState
        title="No se pudo cargar"
        message={error ?? "Inténtalo de nuevo."}
        action={
          <Link to="/scenes">
            <Button variant="secondary">Volver a la lista</Button>
          </Link>
        }
      />
    );
  }

  return (
    <section>
      <Breadcrumbs
        items={[{ label: "Escenas", to: "/scenes" }, { label: scene.title }]}
      />

      <PageHeader
        title={scene.title}
        action={
          <Link to="/scenes">
            <Button variant="secondary">← Volver</Button>
          </Link>
        }
      />

      <dl className={styles.fields}>
        <div className={styles.field}>
          <dt className={styles.label}>ID</dt>
          <dd className={styles.value}>{scene.id}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Texto</dt>
          <dd className={styles.value}>{scene.text ?? "—"}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Línea de tiempo</dt>
          <dd className={styles.value}>
            {scene.startTimeline} → {scene.endTimeline}
          </dd>
        </div>

        {/* Personajes relacionados: cada badge es un enlace a su detalle. */}
        <div className={styles.field}>
          <dt className={styles.label}>Personajes</dt>
          <dd className={styles.value}>
            {scene.characters.length ? (
              <div className={styles.tags}>
                {scene.characters.map((c) => (
                  <Link key={c.id} to={`/characters/${c.id}`}>
                    <Badge tone="primary">{c.title}</Badge>
                  </Link>
                ))}
              </div>
            ) : (
              "—"
            )}
          </dd>
        </div>

        {/* Lugares relacionados: igual, enlazan a su detalle. */}
        <div className={styles.field}>
          <dt className={styles.label}>Lugares</dt>
          <dd className={styles.value}>
            {scene.locations.length ? (
              <div className={styles.tags}>
                {scene.locations.map((l) => (
                  <Link key={l.id} to={`/locations/${l.id}`}>
                    <Badge>{l.title}</Badge>
                  </Link>
                ))}
              </div>
            ) : (
              "—"
            )}
          </dd>
        </div>

        <div className={styles.field}>
          <dt className={styles.label}>Creado</dt>
          <dd className={styles.value}>{formatDate(scene.createdAt)}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Actualizado</dt>
          <dd className={styles.value}>{formatDate(scene.updatedAt)}</dd>
        </div>
      </dl>
    </section>
  );
}
