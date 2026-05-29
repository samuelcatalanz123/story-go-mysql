import type { ReactNode } from "react";
import styles from "./PageHeader.module.css";

type Props = { title: string; action?: ReactNode };

export function PageHeader({ title, action }: Props) {
  return (
    <header className={styles.header}>
      <h1 className={styles.title}>{title}</h1>
      {action && <div>{action}</div>}
    </header>
  );
}
