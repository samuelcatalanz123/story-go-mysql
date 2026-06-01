import { Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "./components/Layout";
import { CharactersPage } from "./pages/CharactersPage";
import { CharacterDetailPage } from "./pages/CharacterDetailPage";
import { LocationsPage } from "./pages/LocationsPage";
import { LocationDetailPage } from "./pages/LocationDetailPage";
import { ScenesPage } from "./pages/ScenesPage";
import { SceneDetailPage } from "./pages/SceneDetailPage";
import { OrganizationsPage } from "./pages/OrganizationsPage";
import { ConflictsPage } from "./pages/ConflictsPage";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { NotFound } from "./pages/NotFound";

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/" element={<Layout />}>
        <Route index element={<Navigate to="/characters" replace />} />
        <Route path="characters" element={<CharactersPage />} />
        <Route path="characters/:id" element={<CharacterDetailPage />} />
        <Route path="locations" element={<LocationsPage />} />
        <Route path="locations/:id" element={<LocationDetailPage />} />
        <Route path="scenes" element={<ScenesPage />} />
        <Route path="scenes/:id" element={<SceneDetailPage />} />
        <Route path="organizations" element={<OrganizationsPage />} />
        <Route path="conflicts" element={<ConflictsPage />} />
        {/* Ruta comodín: cualquier URL que no coincida arriba cae aquí (404). */}
        <Route path="*" element={<NotFound />} />
      </Route>
    </Routes>
  );
}
