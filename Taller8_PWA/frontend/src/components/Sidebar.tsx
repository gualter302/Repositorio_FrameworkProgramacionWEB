import type { ComponentType, SVGProps } from "react";
import { NavLink } from "react-router-dom";
import { BadgeDollarSign, BarChart3, Boxes, ClipboardList, LayoutDashboard, LogIn, Network, PackagePlus, ReceiptText, UserPlus } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { BrandMark } from './Ui';

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
  { to: "/", label: "Vista general", Icon: LayoutDashboard },
  { to: "/catalogo", label: "Catálogo", Icon: Boxes },
  { to: "/pedidos", label: "Mis pedidos", Icon: ReceiptText },
  { to: "/comisiones", label: "Mis comisiones", Icon: BadgeDollarSign },
  { to: "/mi-red", label: "Mi red", Icon: Network },
];

const guestItems: MenuItem[] = [
  { to: '/catalogo', label: 'Catálogo', Icon: Boxes },
  { to: '/login', label: 'Iniciar sesión', Icon: LogIn },
  { to: '/registro', label: 'Crear cuenta', Icon: UserPlus },
];

const Sidebar = ({ isCollapsed, onNavigate }: SidebarProps) => {
  const { isAuthenticated, role } = useAuth();
  const items = !isAuthenticated ? guestItems : role === "admin"
    ? [...menuItems, { to: "/admin/productos", label: "Productos", Icon: PackagePlus }, { to: "/admin/pedidos", label: "Pedidos", Icon: ClipboardList }, { to: "/admin/comisiones", label: "Liquidaciones", Icon: BadgeDollarSign }, { to: "/admin/reportes", label: "Reportes", Icon: BarChart3 }]
    : menuItems;
  return (
    <aside
      // Celular: panel fijo que se desliza desde la izquierda (oculto si isCollapsed).
      // Escritorio (md+): parte del layout; ancho 80px si isCollapsed, 256px si no.
      className={`fixed inset-y-0 left-0 z-40 border-r border-cyan-200/10 md:static md:z-auto ${
        isCollapsed ? "-translate-x-full w-64 md:translate-x-0 md:w-20" : "translate-x-0 w-64"
      } shrink-0 overflow-hidden bg-slate-950/92 text-white flex flex-col backdrop-blur-xl transition-all duration-300`}
    >
      {/* Cabecera: nombre completo o solo las iniciales si está colapsado */}
      <div
        className={`min-h-18 flex items-center border-b border-cyan-200/10 whitespace-nowrap ${
          isCollapsed ? "justify-center" : "px-5"
        }`}
      >
        <BrandMark compact={isCollapsed} />
      </div>

      <nav className="flex-1 space-y-1.5 overflow-y-auto p-3" aria-label="Navegación principal">
        {items.map(({ to, label, Icon }) => (
          <NavLink
            key={to}
            to={to}
            end={to === "/"}
            onClick={onNavigate}
            title={isCollapsed ? label : undefined} // tooltip cuando solo se ven iconos
            aria-label={label}
            className={({ isActive }) =>
              `group flex min-h-12 items-center gap-3 rounded-xl border px-3 py-2.5 whitespace-nowrap transition ${
                isCollapsed ? "justify-center" : ""
              } ${isActive ? "border-cyan-300/35 bg-cyan-300/10 text-cyan-200" : "border-transparent text-slate-400 hover:border-cyan-200/10 hover:bg-white/5 hover:text-white"}`
            }
          >
            <Icon className="h-5 w-5 shrink-0" />
            {!isCollapsed && <span className="text-sm font-medium">{label}</span>}
          </NavLink>
        ))}
      </nav>
      {!isCollapsed && <div className="m-3 rounded-xl border border-cyan-200/10 bg-white/[0.03] p-3 text-xs leading-relaxed text-slate-500">PWA preparada para trabajar con un catálogo guardado cuando no hay conexión.</div>}
    </aside>
  );
};

export default Sidebar;
