import type { ComponentType, SVGProps } from "react";
import { NavLink } from "react-router-dom";
import { CatalogoIcon, DashboardIcon, RedIcon } from "../icons";

// Props del Sidebar: isCollapsed = true -> ancho de 80px y solo iconos.
interface SidebarProps {
  isCollapsed: boolean;
  onNavigate?: () => void; // se ejecuta al elegir una opción (cierra el menú en celular)
}

interface MenuItem {
  to: string;
  label: string;
  Icon: ComponentType<SVGProps<SVGSVGElement>>;
}

// Opciones del menú (cada una con su icono)
const menuItems: MenuItem[] = [
  { to: "/", label: "Dashboard", Icon: DashboardIcon },
  { to: "/catalogo", label: "Catálogo", Icon: CatalogoIcon },
  { to: "/mi-red", label: "Mi Red", Icon: RedIcon },
];

const Sidebar = ({ isCollapsed, onNavigate }: SidebarProps) => {
  return (
    <aside
      // Celular: panel fijo que se desliza desde la izquierda (oculto si isCollapsed).
      // Escritorio (md+): parte del layout; ancho 80px si isCollapsed, 256px si no.
      className={`fixed inset-y-0 left-0 z-40 md:static md:z-auto ${
        isCollapsed ? "-translate-x-full w-64 md:translate-x-0 md:w-20" : "translate-x-0 w-64"
      } shrink-0 overflow-hidden bg-slate-900 text-white flex flex-col transition-all duration-300`}
    >
      {/* Cabecera: nombre completo o solo las iniciales si está colapsado */}
      <div
        className={`h-16 flex items-center border-b border-slate-700 font-bold whitespace-nowrap ${
          isCollapsed ? "justify-center text-xl" : "px-6 text-2xl"
        }`}
      >
        {isCollapsed ? "MC" : "MultiCatálogo"}
      </div>

      <nav className="flex-1 p-4 space-y-2">
        {menuItems.map(({ to, label, Icon }) => (
          <NavLink
            key={to}
            to={to}
            end={to === "/"}
            onClick={onNavigate}
            title={isCollapsed ? label : undefined} // tooltip cuando solo se ven iconos
            aria-label={label}
            className={({ isActive }) =>
              `flex items-center gap-3 p-3 rounded whitespace-nowrap transition ${
                isCollapsed ? "justify-center" : ""
              } ${isActive ? "bg-slate-800 text-white" : "text-slate-300 hover:bg-slate-800 hover:text-white"}`
            }
          >
            <Icon className="w-6 h-6 shrink-0" />
            {!isCollapsed && <span>{label}</span>}
          </NavLink>
        ))}
      </nav>
    </aside>
  );
};

export default Sidebar;
