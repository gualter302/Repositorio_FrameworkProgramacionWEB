import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Navbar from "./Navbar";
import { useSidebar } from "../context/SidebarContext";

const Layout = () => {
  // Leemos el estado global y se lo pasamos al Sidebar como prop
  const { isCollapsed } = useSidebar();

  return (
    <div className="flex h-screen bg-slate-50">
      <Sidebar isCollapsed={isCollapsed} />

      {/* Área de Contenido Principal */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <Navbar />

        <main className="flex-1 overflow-x-hidden overflow-y-auto p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;
