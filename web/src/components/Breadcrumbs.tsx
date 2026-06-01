import { Link } from "react-router-dom";
import styles from "./Breadcrumbs.module.css";

// Un "crumb" (miga) puede ser un enlace (lleva `to`) o solo texto (la página
// actual, sin enlace porque ya estás en ella).
export type Crumb = { label: string; to?: string };

type Props = { items: Crumb[] };

// Breadcrumbs muestra la ruta de navegación: Personajes › Asha
// Ayuda al usuario a saber dónde está y a volver atrás con un clic.
export function Breadcrumbs({ items }: Props) {
  return (
    <nav className={styles.breadcrumbs} aria-label="Ruta de navegación">
      {items.map((item, index) => {
        const isLast = index === items.length - 1;
        return (
          <span key={index} className={styles.crumb}>
            {item.to && !isLast ? (
              <Link to={item.to} className={styles.link}>
                {item.label}
              </Link>
            ) : (
              <span className={styles.current}>{item.label}</span>
            )}
            {!isLast && <span className={styles.separator}>›</span>}
          </span>
        );
      })}
    </nav>
  );
}
