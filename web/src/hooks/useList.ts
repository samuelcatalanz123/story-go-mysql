import { useCallback, useEffect, useState } from "react";

export type ListState<T> = {
  data: T[];
  loading: boolean;
  error: string | null;
  reload: () => void;
};

// Carga una lista al montar el componente y expone estado de
// carga/error/datos, además de una función para recargar.
// `loader` debe ser una referencia estable (p. ej. una función de api/resources.ts).
export function useList<T>(loader: () => Promise<T[]>): ListState<T> {
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError(null);
    loader()
      .then((items) => setData(items))
      .catch((e: unknown) =>
        setError(e instanceof Error ? e.message : "Error desconocido"),
      )
      .finally(() => setLoading(false));
  }, [loader]);

  useEffect(() => {
    load();
  }, [load]);

  return { data, loading, error, reload: load };
}
