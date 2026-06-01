import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { getLocation, uploadLocationAvatar } from "../api/resources";
import { ApiError } from "../api/client";
import type { Location } from "../types";
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

export function LocationDetailPage() {
  const { id } = useParams();
  const { isAuthenticated } = useAuth();
  const [location, setLocation] = useState<Location | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    setLoading(true);
    setNotFound(false);
    setError(null);

    getLocation(Number(id))
      .then((data) => {
        if (active) setLocation(data);
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

  if (loading) return <SkeletonRows rows={3} cols={1} />;
  if (notFound) return <NotFound />;
  if (error || !location) {
    return (
      <EmptyState
        title="No se pudo cargar"
        message={error ?? "Inténtalo de nuevo."}
        action={
          <Link to="/locations">
            <Button variant="secondary">Volver a la lista</Button>
          </Link>
        }
      />
    );
  }

  return (
    <section>
      <Breadcrumbs
        items={[
          { label: "Lugares", to: "/locations" },
          { label: location.title },
        ]}
      />

      <PageHeader
        title={location.title}
        action={
          <Link to="/locations">
            <Button variant="secondary">← Volver</Button>
          </Link>
        }
      />

      <AvatarUploader
        avatarPath={location.avatarPath}
        canEdit={isAuthenticated}
        onUpload={(file) => uploadLocationAvatar(location.id, file)}
      />

      <dl className={styles.fields}>
        <div className={styles.field}>
          <dt className={styles.label}>ID</dt>
          <dd className={styles.value}>{location.id}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Texto</dt>
          <dd className={styles.value}>{location.text ?? "—"}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Creado</dt>
          <dd className={styles.value}>{formatDate(location.createdAt)}</dd>
        </div>
        <div className={styles.field}>
          <dt className={styles.label}>Actualizado</dt>
          <dd className={styles.value}>{formatDate(location.updatedAt)}</dd>
        </div>
      </dl>
    </section>
  );
}
