import { apiFetch } from "./client";
import type { AuthResponse, LoginRequest, RegisterRequest } from "../types";

export const register = (body: RegisterRequest) =>
  apiFetch<AuthResponse>("/auth/register", { method: "POST", body: JSON.stringify(body) });

export const login = (body: LoginRequest) =>
  apiFetch<AuthResponse>("/auth/login", { method: "POST", body: JSON.stringify(body) });
