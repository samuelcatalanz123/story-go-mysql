import { NavLink, Outlet } from "react-router-dom";
import styles from "./Layout.module.css";

export function Layout() {
  const linkClass = ({ isActive }: { isActive: boolean }) =>
    isActive ? styles.active : "";
  return (
    <div className={styles.shell}>
      <nav className={styles.sidebar}>
        <h2>Story Admin</h2>
        <NavLink to="/characters" className={linkClass}>
          Personajes
        </NavLink>
        <NavLink to="/locations" className={linkClass}>
          Lugares
        </NavLink>
        <NavLink to="/scenes" className={linkClass}>
          Escenas
        </NavLink>
      </nav>
      <main className={styles.content}>
        <Outlet />
      </main>
    </div>
  );
}
