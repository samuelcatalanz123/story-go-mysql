import { useEffect, useState } from "react";
import { PageHeader } from "../ui/PageHeader";
import { Button } from "../ui/Button";

// Consulta GraphQL de ejemplo: pedimos SOLO los campos que queremos.
const QUERY = `{
  characters { id title }
  scenes { id title characters { title } }
}`;

// GraphQLDemoPage muestra cómo se consulta GraphQL con fetch plano (sin Apollo
// ni librerías): un POST a /api/graphql con { query }. Compara con el REST.
export function GraphQLDemoPage() {
  const [result, setResult] = useState<string>("Cargando…");

  async function run() {
    setResult("Cargando…");
    try {
      const res = await fetch("/api/graphql", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query: QUERY }),
      });
      const data = await res.json();
      setResult(JSON.stringify(data, null, 2));
    } catch (e) {
      setResult("Error: " + (e instanceof Error ? e.message : "desconocido"));
    }
  }

  useEffect(() => {
    void run();
  }, []);

  return (
    <section>
      <PageHeader title="GraphQL (demo)" action={<Button onClick={() => void run()}>Reejecutar</Button>} />
      <p>
        Esta es la MISMA data que el REST, pero pedida con una sola consulta GraphQL. Para
        explorar el API de forma interactiva, abre el{" "}
        <a href="/api/playground" target="_blank" rel="noreferrer">
          playground (GraphiQL)
        </a>
        .
      </p>
      <h3>Consulta</h3>
      <pre style={{ background: "var(--color-surface)", border: "1px solid var(--color-border)", borderRadius: 8, padding: 12, overflowX: "auto" }}>
        {QUERY}
      </pre>
      <h3>Respuesta</h3>
      <pre style={{ background: "var(--color-surface)", border: "1px solid var(--color-border)", borderRadius: 8, padding: 12, overflowX: "auto" }}>
        {result}
      </pre>
    </section>
  );
}
