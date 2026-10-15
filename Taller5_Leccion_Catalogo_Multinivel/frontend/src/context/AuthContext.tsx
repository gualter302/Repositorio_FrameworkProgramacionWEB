// src/context/AuthContext.tsx
import { createContext, useContext, useState, type ReactNode } from 'react';

// 1. Tipos de rol que maneja la aplicación (Tema 5)
export type Rol = 'admin' | 'cliente';

// 2. Usuario autenticado: correo + rol (lo entrega la API en POST /api/login)
export interface Usuario {
  email: string;
  rol: Rol;
}

interface AuthContextType {
  isAuthenticated: boolean;
  user: Usuario | null; // Usuario completo (correo + rol)
  token: string | null; // Token que devuelve la API al iniciar sesión
  login: (usuario: Usuario, token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe ser usado dentro de un AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<Usuario | null>(null);
  const [token, setToken] = useState<string | null>(null);

  const login = (usuario: Usuario, token: string) => {
    setIsAuthenticated(true);
    setUser(usuario); // Guardamos correo y rol
    setToken(token); // Guardamos el token
  };

  const logout = () => {
    setIsAuthenticated(false);
    setUser(null); // Limpiamos el usuario al salir
    setToken(null); // Limpiamos el token al salir
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
