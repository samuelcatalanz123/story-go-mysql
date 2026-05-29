import { apiFetch } from "./client";
import type {
  Character,
  Location,
  Scene,
  CharacterRequest,
  LocationRequest,
  SceneRequest,
  Paged,
} from "../types";

export type ListArgs = { q: string; page: number; pageSize: number };

function listQuery({ q, page, pageSize }: ListArgs): string {
  return new URLSearchParams({
    q,
    page: String(page),
    pageSize: String(pageSize),
  }).toString();
}

// --- Personajes ---
export const listCharacters = (args: ListArgs) =>
  apiFetch<Paged<Character>>(`/characters?${listQuery(args)}`);
export const createCharacter = (body: CharacterRequest) =>
  apiFetch<Character>("/characters", { method: "POST", body: JSON.stringify(body) });
export const updateCharacter = (id: number, body: CharacterRequest) =>
  apiFetch<Character>(`/characters/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteCharacter = (id: number) =>
  apiFetch<void>(`/characters/${id}`, { method: "DELETE" });

// --- Lugares ---
export const listLocations = (args: ListArgs) =>
  apiFetch<Paged<Location>>(`/locations?${listQuery(args)}`);
export const createLocation = (body: LocationRequest) =>
  apiFetch<Location>("/locations", { method: "POST", body: JSON.stringify(body) });
export const updateLocation = (id: number, body: LocationRequest) =>
  apiFetch<Location>(`/locations/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteLocation = (id: number) =>
  apiFetch<void>(`/locations/${id}`, { method: "DELETE" });

// --- Escenas ---
export const listScenes = (args: ListArgs) =>
  apiFetch<Paged<Scene>>(`/scenes?${listQuery(args)}`);
export const createScene = (body: SceneRequest) =>
  apiFetch<Scene>("/scenes", { method: "POST", body: JSON.stringify(body) });
export const updateScene = (id: number, body: SceneRequest) =>
  apiFetch<Scene>(`/scenes/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteScene = (id: number) =>
  apiFetch<void>(`/scenes/${id}`, { method: "DELETE" });

// Helpers que traen TODOS los elementos (para los desplegables del formulario
// de escenas), usando un pageSize grande.
const ALL: ListArgs = { q: "", page: 1, pageSize: 1000 };
export const listAllCharacters = () =>
  apiFetch<Paged<Character>>(`/characters?${listQuery(ALL)}`).then((p) => p.items);
export const listAllLocations = () =>
  apiFetch<Paged<Location>>(`/locations?${listQuery(ALL)}`).then((p) => p.items);
