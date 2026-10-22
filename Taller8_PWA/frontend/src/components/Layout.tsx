import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Navbar from "./Navbar";
import { useSidebar } from "../context/SidebarContext";

const Layout = () => {
  // Leemos el estado global y se lo pasamos al Sidebar como prop
  const { isCollapsed, toggleSidebar, closeOnMobile } = useSidebar();

  return (
    <div className="app-shell">
      <Sidebar isCollapsed={isCollapsed} onNavigate={closeOnMobile} />

      {/* Fondo oscuro detrás del menú abierto (solo celular). Al tocarlo se cierra. */}
      {!isCollapsed && (
        <div
          className="fixed inset-0 z-30 bg-slate-950/80 backdrop-blur-sm md:hidden"
          onClick={toggleSidebar}
          aria-hidden="true"
        />
      )}

      {/* Área de Contenido Principal */}
      <div className="content-column">
        <Navbar />

        <main className="content-scroll">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;
