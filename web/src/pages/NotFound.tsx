import { Link } from "react-router-dom";
import { EmptyState } from "../ui/EmptyState";
import { Button } from "../ui/Button";

// NotFound es la página 404. Se muestra en dos casos:
//  1) El usuario escribe una URL que no existe (ruta comodín "*" en App.tsx).
//  2) Pide el detalle de un elemento cuyo id no existe en el backend (404).
export function NotFound() {
  return (
    <EmptyState
      title="404 — No encontrado"
      message="La página o el elemento que buscas no existe."
      action={
        <Link to="/characters">
          <Button>Volver al inicio</Button>
        </Link>
      }
    />
  );
}
