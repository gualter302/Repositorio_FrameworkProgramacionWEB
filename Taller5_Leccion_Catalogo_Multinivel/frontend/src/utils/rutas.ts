// src/utils/rutas.ts
// Ruta de inicio de cada rol (Tema 5): a dónde se envía al usuario al iniciar sesión
// o cuando intenta entrar a una ruta que su rol no puede ver.
import type { Rol } from '../context/AuthContext';

// Vista por defecto del cliente: la tienda a pantalla completa (Paso 7)
export const INICIO_CLIENTE = '/tienda';

// Admin -> Dashboard ("/"), cliente -> Tienda
export const rutaInicioPorRol = (rol: Rol): string => (rol === 'admin' ? '/' : INICIO_CLIENTE);
