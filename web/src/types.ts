// Tipos que reflejan los modelos JSON de la API de Go.
// En Go los campos *string (anulables) se representan como string | null.

export type Character = {
  id: number;
  title: string;
  text: string | null;
  createdAt: string;
  updatedAt: string;
};

export type Location = {
  id: number;
  title: string;
  text: string | null;
  createdAt: string;
  updatedAt: string;
};

export type Scene = {
  id: number;
  title: string;
  text: string | null;
  startTimeline: number;
  endTimeline: number;
  characters: Character[];
  locations: Location[];
  createdAt: string;
  updatedAt: string;
};

export type CharacterRequest = { title: string; text: string | null };
export type LocationRequest = { title: string; text: string | null };
export type SceneRequest = {
  title: string;
  text: string | null;
  startTimeline: number;
  endTimeline: number;
  characterIds: number[];
  locationIds: number[];
};

export type User = { id: number; email: string; createdAt: string };
export type AuthResponse = { token: string; user: User };
export type RegisterRequest = { email: string; password: string };
export type LoginRequest = { email: string; password: string };

export type Paged<T> = {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
};
