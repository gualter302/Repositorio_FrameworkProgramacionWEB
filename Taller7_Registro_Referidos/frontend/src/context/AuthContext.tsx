import { createContext, useContext, useState, type ReactNode } from 'react';

export const SESSION_KEY = 'multicatalogo_sesion';

interface Session {
  email: string;
  nombre: string;
  rol: string;
  token: string;
}

interface AuthContextType {
  isAuthenticated: boolean;
  userEmail: string | null;
  userName: string | null;
  role: string | null;
  token: string | null;
  login: (session: Session) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const cargarSesion = (): Session | null => {
  try {
    const guardada = localStorage.getItem(SESSION_KEY);
    return guardada ? JSON.parse(guardada) as Session : null;
  } catch {
    localStorage.removeItem(SESSION_KEY);
    return null;
  }
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth debe ser usado dentro de un AuthProvider');
  return context;
};

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [session, setSession] = useState<Session | null>(cargarSesion);

  const login = (nuevaSesion: Session) => {
    localStorage.setItem(SESSION_KEY, JSON.stringify(nuevaSesion));
    setSession(nuevaSesion);
  };

  const logout = () => {
    localStorage.removeItem(SESSION_KEY);
    setSession(null);
  };

  return (
    <AuthContext.Provider value={{
      isAuthenticated: Boolean(session?.token),
      userEmail: session?.email ?? null,
      userName: session?.nombre ?? null,
      role: session?.rol ?? null,
      token: session?.token ?? null,
      login,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  );
};
