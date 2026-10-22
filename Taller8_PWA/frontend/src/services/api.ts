// src/services/api.ts
// Único lugar donde el frontend habla con el backend (Go + Fiber).
// Ningún componente usa fetch directamente: todos llaman a las funciones de este archivo.
import type { Producto } from '../context/CartContext';
import { readSession, expireSession } from './session';

// URL base de la API. Se puede cambiar con la variable VITE_API_URL (archivo .env);
// si no existe, se usa el backend local en el puerto 3000.
export const API_URL: string = import.meta.env.VITE_API_URL ?? '/api';
export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) { super(message); this.status = status; }
}

// Forma de la respuesta de POST /api/login (coincide con lo que devuelve authController.go)
export interface LoginResponse {
  token: string;
  email: string;
  nombre: string;
  rol: string;
}

export interface RegisterRequest {
  nombre: string;
  email: string;
  password: string;
  ref?: string;
}

export interface Referido {
  id: number;
  nombre: string;
  codigo: string;
  nivel: number;
  ventas: number;
  comisiones: number;
  hijos: Referido[];
}

export interface CheckoutRequest {
  usuarioEmail: string;
  cliente: {
    nombre: string;
    email: string;
    direccion: string;
    ciudad: string;
    telefono: string;
  };
  metodoPago: string;
  items: { productoId: number; cantidad: number }[];
}

export interface PedidoConfirmado {
  id: number;
  numero: string;
  total: number;
  estado: string;
  fecha: string;
  cliente: CheckoutRequest['cliente'];
  items: {
    productoId: number;
    nombre: string;
    cantidad: number;
    precio: number;
    subtotal: number;
  }[];
}

export interface PedidoResumen extends PedidoConfirmado {
  usuarioEmail: string;
}

export interface Comision {
  id: number;
  pedidoId: number;
  numeroPedido: string;
  referidoId: number;
  referido: string;
  nivel: number;
  porcentaje: number;
  monto: number;
  estado: 'pendiente' | 'pagada';
  fecha: string;
  pagadaEn?: string;
}

export interface ReporteGeneral {
  totalVentas: number;
  totalPedidos: number;
  ticketPromedio: number;
  totalComisiones: number;
  pendientePagar: number;
  productosPrincipales: { nombre: string; cantidad: number; ventas: number }[];
  referidosPrincipales: { nombre: string; ventas: number; comisiones: number }[];
}

interface CatalogoGuardado {
  products: Producto[];
  date: number;
}

const CATALOGO_KEY = 'multicatalogo_catalogo';
const CATALOGO_FECHA_KEY = 'multicatalogo_catalogo_fecha';
const DURACION_CATALOGO = 86400000;

export const getProductosGuardados = (): CatalogoGuardado | null => {
  try {
    const saved = JSON.parse(localStorage.getItem(CATALOGO_KEY) ?? 'null');
    if (!Array.isArray(saved?.products) || typeof saved.date !== 'number') return null;
    if (Date.now() - saved.date >= DURACION_CATALOGO) return null;
    return { products: saved.products, date: saved.date };
  } catch {
    return null;
  }
};

// Petición genérica: hace el fetch, convierte a JSON y lanza un Error si algo sale mal.
const request = async <T>(path: string, options?: RequestInit): Promise<T> => {
  let response: Response;
  try {
    const token = readSession()?.token;
    const headers = new Headers(options?.headers);
    if (token) headers.set('Authorization', `Bearer ${token}`);
    response = await fetch(`${API_URL}${path}`, { ...options, headers, cache: 'no-store' });
  } catch {
    // fetch falla (no responde) cuando el backend está apagado o CORS bloquea la petición
    throw new Error('No se pudo conectar con el servidor. ¿Está encendido el backend?');
  }

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    if (response.status === 401 && path !== '/login') expireSession();
    // El backend responde { "error": "mensaje" } cuando algo falla
    throw new ApiError(data?.error ?? `Error ${response.status} del servidor`, response.status);
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

// POST /api/register
export const registerRequest = (datos: RegisterRequest) =>
  request<LoginResponse>('/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(datos),
  });

// GET /api/productos
export const getProductos = async (freshOnly = false): Promise<Producto[]> => {
  try {
    const products = await request<Producto[]>('/productos');
    const snapshot = JSON.stringify({ products: products.slice(0, 500), date: Date.now() });
    try {
      if (snapshot.length <= 1000000) {
        localStorage.setItem(CATALOGO_KEY, snapshot);
        localStorage.setItem(CATALOGO_FECHA_KEY, String(Date.now()));
      }
    } catch { /* El catálogo en línea funciona incluso si el almacenamiento está lleno. */ }
    window.dispatchEvent(new Event('catalogue-updated'));
    return products;
  } catch (error) {
    if (freshOnly) throw error;
    const saved = getProductosGuardados();
    if (saved) return saved.products;
    throw error;
  }
};

export const getProductosAdmin = () => request<Producto[]>('/admin/productos');

export const crearProducto = (producto: Omit<Producto, 'id' | 'activo'>) =>
  request<Producto>('/admin/productos', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(producto),
  });

export const actualizarProducto = (id: number, producto: Omit<Producto, 'id' | 'activo'>) =>
  request<Producto>(`/admin/productos/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(producto),
  });

export const desactivarProducto = (id: number) =>
  request<void>(`/admin/productos/${id}`, { method: 'DELETE' });

// GET /api/red/:email
export const getRed = () => request<Referido[]>('/red');

// POST /api/checkout
export const checkoutRequest = (checkout: CheckoutRequest, key: string) =>
  request<PedidoConfirmado>('/checkout', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'Idempotency-Key': key },
    body: JSON.stringify(checkout),
  });

export const getMisPedidos = () => request<PedidoResumen[]>('/pedidos');

export const getPedidosAdmin = (pagina = 1, buscar = '') =>
  request<{ datos: PedidoResumen[]; pagina: number; limite: number; total: number }>(
    `/admin/pedidos?pagina=${pagina}&limite=10&buscar=${encodeURIComponent(buscar)}`,
  );

export const actualizarEstadoPedido = (id: number, estado: string) =>
  request<void>(`/admin/pedidos/${id}/estado`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ estado }),
  });

export const cancelarPedido = (id: number) =>
  request<void>(`/admin/pedidos/${id}/cancelar`, { method: 'POST' });

export const getMisComisiones = () => request<Comision[]>('/comisiones');

export const getComisionesAdmin = (estado = '') =>
  request<Comision[]>(`/admin/comisiones?estado=${encodeURIComponent(estado)}`);

export const liquidarComisiones = (comisionIds: number[]) =>
  request<{ id: number; monto: number; fecha: string }>('/admin/liquidaciones', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ comisionIds }),
  });

export const getReporteGeneral = () => request<ReporteGeneral>('/admin/reportes/general');
