// src/services/api.ts
// Único lugar donde el frontend habla con el backend (Go + Fiber).
// Ningún componente usa fetch directamente: todos llaman a las funciones de este archivo.
import type { Producto } from '../context/CartContext';

// URL base de la API. Se puede cambiar con la variable VITE_API_URL (archivo .env);
// si no existe, se usa el backend local en el puerto 3000.
export const API_URL: string = import.meta.env.VITE_API_URL ?? 'http://localhost:3000/api';

// Forma de la respuesta de POST /api/login (coincide con lo que devuelve authController.go)
export interface LoginResponse {
  token: string;
  email: string;
  rol: string; // "admin" | "cliente" (Tema 5)
}

// Petición genérica: hace el fetch, convierte a JSON y lanza un Error si algo sale mal.
const request = async <T>(path: string, options?: RequestInit): Promise<T> => {
  let response: Response;
  try {
    response = await fetch(`${API_URL}${path}`, options);
  } catch {
    // fetch falla (no responde) cuando el backend está apagado o CORS bloquea la petición
    throw new Error('No se pudo conectar con el servidor. ¿Está encendido el backend?');
  }

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    // El backend responde { "error": "mensaje" } cuando algo falla
    throw new Error(data?.error ?? `Error ${response.status} del servidor`);
  }
  return data as T;
};

// POST /api/login
export const loginRequest = (email: string, password: string) =>
  request<LoginResponse>('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

// GET /api/productos
export const getProductos = () => request<Producto[]>('/productos');
