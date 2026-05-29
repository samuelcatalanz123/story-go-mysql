import styles from "./Skeleton.module.css";

type Props = { width?: string; height?: string };

export function Skeleton({ width = "100%", height = "1rem" }: Props) {
  return <span className={styles.skeleton} style={{ width, height }} />;
}

export function SkeletonRows({ rows = 3, cols = 4 }: { rows?: number; cols?: number }) {
  return (
    <div className={styles.rows}>
      {Array.from({ length: rows }).map((_, r) => (
        <div key={r} className={styles.row}>
          {Array.from({ length: cols }).map((_, c) => (
            <Skeleton key={c} height="1.25rem" />
          ))}
        </div>
      ))}
    </div>
  );
}
