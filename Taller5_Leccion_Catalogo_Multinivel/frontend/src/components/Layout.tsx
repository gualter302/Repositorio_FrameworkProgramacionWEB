import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Navbar from "./Navbar";
import { useSidebar } from "../context/SidebarContext";

const Layout = () => {
  // Leemos el estado global y se lo pasamos al Sidebar como prop
  const { isCollapsed, toggleSidebar, closeOnMobile } = useSidebar();

  return (
    <div className="flex h-dvh bg-slate-50">
      <Sidebar isCollapsed={isCollapsed} onNavigate={closeOnMobile} />

      {/* Fondo oscuro detrás del menú abierto (solo celular). Al tocarlo se cierra. */}
      {!isCollapsed && (
        <div
          className="fixed inset-0 z-30 bg-black/50 md:hidden"
          onClick={toggleSidebar}
          aria-hidden="true"
        />
      )}

      {/* Área de Contenido Principal */}
      <div className="flex-1 min-w-0 flex flex-col overflow-hidden">
        <Navbar />

        <main className="flex-1 overflow-x-hidden overflow-y-auto p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;
