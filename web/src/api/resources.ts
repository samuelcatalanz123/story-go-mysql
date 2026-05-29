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
// de escenas). El backend limita pageSize a 100, así que paginamos hasta
// agotar el total en vez de pedir una sola página enorme.
async function fetchAll<T>(loader: (args: ListArgs) => Promise<Paged<T>>): Promise<T[]> {
  const pageSize = 100;
  const all: T[] = [];
  let page = 1;
  for (;;) {
    const res = await loader({ q: "", page, pageSize });
    all.push(...res.items);
    if (res.items.length === 0 || all.length >= res.total) break;
    page += 1;
  }
  return all;
}

export const listAllCharacters = () => fetchAll(listCharacters);
export const listAllLocations = () => fetchAll(listLocations);
