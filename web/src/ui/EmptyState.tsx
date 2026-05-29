import type { ReactNode } from "react";
import styles from "./EmptyState.module.css";

type Props = { title: string; message?: string; action?: ReactNode };

export function EmptyState({ title, message, action }: Props) {
  return (
    <div className={styles.empty}>
      <h3 className={styles.title}>{title}</h3>
      {message && <p className={styles.message}>{message}</p>}
      {action && <div className={styles.action}>{action}</div>}
    </div>
  );
}
