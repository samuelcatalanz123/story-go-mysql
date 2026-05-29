import { useEffect, useState } from "react";
import styles from "./SearchBar.module.css";

type Props = { onQueryChange: (q: string) => void; placeholder?: string };

// Input de búsqueda con debounce: espera ~300 ms tras la última tecla antes de
// avisar, para no consultar la API en cada pulsación.
export function SearchBar({ onQueryChange, placeholder = "Buscar…" }: Props) {
  const [value, setValue] = useState("");

  useEffect(() => {
    const timer = setTimeout(() => onQueryChange(value.trim()), 300);
    return () => clearTimeout(timer);
  }, [value, onQueryChange]);

  return (
    <input
      className={styles.search}
      type="search"
      value={value}
      placeholder={placeholder}
      aria-label="Buscar"
      onChange={(e) => setValue(e.target.value)}
    />
  );
}
