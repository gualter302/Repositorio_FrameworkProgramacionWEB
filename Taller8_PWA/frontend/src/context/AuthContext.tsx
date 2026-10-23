import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { SESSION_KEY, SESSION_EVENT, readSession } from '../services/session';


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
  login: (session: Session) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const cargarSesion = (): Session | null => {
  try {
    return readSession();
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
  useEffect(() => {
    const check = () => setSession(readSession());
    window.addEventListener(SESSION_EVENT, check);
    const timer = window.setInterval(check, 30000);
    return () => { window.removeEventListener(SESSION_EVENT, check); window.clearInterval(timer); };
  }, []);

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
      login,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  );
};
