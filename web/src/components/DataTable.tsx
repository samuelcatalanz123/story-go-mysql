import type { ReactNode } from "react";
import { Button } from "../ui/Button";
import styles from "./DataTable.module.css";

export type Column<T> = {
  header: string;
  render: (row: T) => ReactNode;
};

type Props<T> = {
  columns: Column<T>[];
  rows: T[];
  onEdit?: (row: T) => void;
  onDelete?: (row: T) => void;
};

export function DataTable<T extends { id: number }>({
  columns,
  rows,
  onEdit,
  onDelete,
}: Props<T>) {
  if (rows.length === 0) {
    return <p>No hay elementos todavía.</p>;
  }

  const hasActions = Boolean(onEdit || onDelete);

  return (
    <table className={styles.table}>
      <thead>
        <tr>
          {columns.map((c) => (
            <th key={c.header}>{c.header}</th>
          ))}
          {hasActions && <th>Acciones</th>}
        </tr>
      </thead>
      <tbody>
        {rows.map((row) => (
          <tr key={row.id}>
            {columns.map((c) => (
              <td key={c.header}>{c.render(row)}</td>
            ))}
            {hasActions && (
              <td>
                <div className={styles.actions}>
                  {onEdit && (
                    <Button variant="ghost" size="sm" onClick={() => onEdit(row)}>
                      Editar
                    </Button>
                  )}
                  {onDelete && (
                    <Button variant="danger" size="sm" onClick={() => onDelete(row)}>
                      Borrar
                    </Button>
                  )}
                </div>
              </td>
            )}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
