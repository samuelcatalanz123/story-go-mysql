import { apiFetch } from "./client";
import type { AuthResponse, LoginRequest, RegisterRequest } from "../types";

export const register = (body: RegisterRequest) =>
  apiFetch<AuthResponse>("/auth/register", { method: "POST", body: JSON.stringify(body) });

export const login = (body: LoginRequest) =>
  apiFetch<AuthResponse>("/auth/login", { method: "POST", body: JSON.stringify(body) });

export const oauthGoogle = (code: string, codeVerifier: string) =>
  apiFetch<AuthResponse>("/auth/oauth/google", {
    method: "POST",
    body: JSON.stringify({ code, codeVerifier }),
  });

type MessageResponse = { message: string };

export const forgotPassword = (email: string) =>
  apiFetch<MessageResponse>("/auth/forgot-password", {
    method: "POST",
    body: JSON.stringify({ email }),
  });

export const resetPassword = (token: string, newPassword: string) =>
  apiFetch<MessageResponse>("/auth/reset-password", {
    method: "POST",
    body: JSON.stringify({ token, newPassword }),
  });

export const verifyEmail = (token: string) =>
  apiFetch<MessageResponse>("/auth/verify-email", {
    method: "POST",
    body: JSON.stringify({ token }),
  });

export const resendVerification = () =>
  apiFetch<MessageResponse>("/auth/resend-verification", { method: "POST" });
