// src/context/SidebarContext.tsx
// Estado GLOBAL del sidebar: cualquier componente (Navbar, Layout, Sidebar...) puede leerlo o cambiarlo.
import { createContext, useContext, useState, type ReactNode } from 'react';

interface SidebarContextType {
  isCollapsed: boolean;      // true = sidebar reducido (solo iconos)
  toggleSidebar: () => void; // alterna entre colapsado y expandido
  closeOnMobile: () => void; // cierra el menú deslizable (solo aplica en celular)
}

// Debe coincidir con el breakpoint "md" de Tailwind (768px)
const MOBILE_BREAKPOINT = 768;

const SidebarContext = createContext<SidebarContextType | undefined>(undefined);

export const useSidebar = () => {
  const context = useContext(SidebarContext);
  if (!context) {
    throw new Error("useSidebar debe ser usado dentro de un SidebarProvider");
  }
  return context;
};

interface SidebarProviderProps {
  children: ReactNode;
}

export const SidebarProvider = ({ children }: SidebarProviderProps) => {
  // En celular el sidebar empieza cerrado; en escritorio, expandido.
  const [isCollapsed, setIsCollapsed] = useState<boolean>(
    () => window.innerWidth < MOBILE_BREAKPOINT
  );

  const toggleSidebar = () => setIsCollapsed((prev) => !prev);

  // En celular, al elegir una opción del menú lo cerramos; en escritorio no hace nada.
  const closeOnMobile = () => {
    if (window.innerWidth < MOBILE_BREAKPOINT) setIsCollapsed(true);
  };

  return (
    <SidebarContext.Provider value={{ isCollapsed, toggleSidebar, closeOnMobile }}>
      {children}
    </SidebarContext.Provider>
  );
};
