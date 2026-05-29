import type { ReactNode } from "react";
import styles from "./Field.module.css";

type Props = {
  label: string;
  htmlFor: string;
  error?: string | null;
  children: ReactNode;
};

export function Field({ label, htmlFor, error, children }: Props) {
  return (
    <div className={styles.field}>
      <label htmlFor={htmlFor} className={styles.label}>
        {label}
      </label>
      {children}
      {error && <span className={styles.error}>{error}</span>}
    </div>
  );
}
