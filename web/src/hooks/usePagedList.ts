import { useCallback, useEffect, useRef, useState } from "react";
import type { Paged } from "../types";

const PAGE_SIZE = 20;

export type PagedListState<T> = {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  query: string;
  loading: boolean;
  error: string | null;
  setQuery: (q: string) => void;
  setPage: (page: number) => void;
  reload: () => void;
};

// Carga una lista paginada y con búsqueda. `loader` debe ser una referencia
// estable (una función de api/resources.ts). Descarta respuestas obsoletas.
export function usePagedList<T>(
  loader: (args: { q: string; page: number; pageSize: number }) => Promise<Paged<T>>,
): PagedListState<T> {
  const [data, setData] = useState<T[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [query, setQueryState] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const requestId = useRef(0);

  const load = useCallback(() => {
    const id = ++requestId.current;
    setLoading(true);
    setError(null);
    loader({ q: query, page, pageSize: PAGE_SIZE })
      .then((res) => {
        if (id !== requestId.current) return;
        setData(res.items);
        setTotal(res.total);
      })
      .catch((e: unknown) => {
        if (id === requestId.current)
          setError(e instanceof Error ? e.message : "Error desconocido");
      })
      .finally(() => {
        if (id === requestId.current) setLoading(false);
      });
  }, [loader, query, page]);

  useEffect(() => {
    load();
    return () => {
      requestId.current++;
    };
  }, [load]);

  const setQuery = useCallback((q: string) => {
    setQueryState(q);
    setPage(1); // nueva búsqueda → primera página
  }, []);

  return {
    data,
    total,
    page,
    pageSize: PAGE_SIZE,
    query,
    loading,
    error,
    setQuery,
    setPage,
    reload: load,
  };
}
