import { Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "./components/Layout";
import { CharactersPage } from "./pages/CharactersPage";
import { LocationsPage } from "./pages/LocationsPage";
import { ScenesPage } from "./pages/ScenesPage";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/" element={<Layout />}>
        <Route index element={<Navigate to="/characters" replace />} />
        <Route path="characters" element={<CharactersPage />} />
        <Route path="locations" element={<LocationsPage />} />
        <Route path="scenes" element={<ScenesPage />} />
      </Route>
    </Routes>
  );
}
