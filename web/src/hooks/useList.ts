import { useCallback, useEffect, useRef, useState } from "react";

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

  // Identifica la petición más reciente. Si llega la respuesta de una petición
  // antigua, o el componente ya se desmontó, la descartamos. Esto evita pisar
  // datos nuevos con otros viejos y actualizar el estado de un componente que
  // ya no está en pantalla (algo fácil de provocar con StrictMode).
  const requestId = useRef(0);

  const load = useCallback(() => {
    const id = ++requestId.current;
    setLoading(true);
    setError(null);
    loader()
      .then((items) => {
        if (id === requestId.current) setData(items);
      })
      .catch((e: unknown) => {
        if (id === requestId.current)
          setError(e instanceof Error ? e.message : "Error desconocido");
      })
      .finally(() => {
        if (id === requestId.current) setLoading(false);
      });
  }, [loader]);

  useEffect(() => {
    load();
    // Al desmontar (o reejecutar el efecto) invalidamos la petición en curso.
    return () => {
      requestId.current++;
    };
  }, [load]);

  return { data, loading, error, reload: load };
}
