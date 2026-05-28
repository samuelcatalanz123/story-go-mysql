import { apiFetch } from "./client";
import type {
  Character,
  Location,
  Scene,
  CharacterRequest,
  LocationRequest,
  SceneRequest,
} from "../types";

// --- Personajes ---
export const listCharacters = () => apiFetch<Character[]>("/characters");
export const createCharacter = (body: CharacterRequest) =>
  apiFetch<Character>("/characters", { method: "POST", body: JSON.stringify(body) });
export const updateCharacter = (id: number, body: CharacterRequest) =>
  apiFetch<Character>(`/characters/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteCharacter = (id: number) =>
  apiFetch<void>(`/characters/${id}`, { method: "DELETE" });

// --- Lugares ---
export const listLocations = () => apiFetch<Location[]>("/locations");
export const createLocation = (body: LocationRequest) =>
  apiFetch<Location>("/locations", { method: "POST", body: JSON.stringify(body) });
export const updateLocation = (id: number, body: LocationRequest) =>
  apiFetch<Location>(`/locations/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteLocation = (id: number) =>
  apiFetch<void>(`/locations/${id}`, { method: "DELETE" });

// --- Escenas ---
export const listScenes = () => apiFetch<Scene[]>("/scenes");
export const createScene = (body: SceneRequest) =>
  apiFetch<Scene>("/scenes", { method: "POST", body: JSON.stringify(body) });
export const updateScene = (id: number, body: SceneRequest) =>
  apiFetch<Scene>(`/scenes/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteScene = (id: number) =>
  apiFetch<void>(`/scenes/${id}`, { method: "DELETE" });
