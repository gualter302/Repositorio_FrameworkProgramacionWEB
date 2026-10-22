import { Link, useLocation, useNavigate } from "react-router-dom";
import { LogIn, LogOut, Menu, ShoppingCart, UserRound } from "lucide-react";
import { useCart } from "../context/CartContext";
import { useAuth } from "../context/AuthContext";
import { useSidebar } from "../context/SidebarContext";

const routeTitles: Record<string, string> = {
  '/': 'Vista general',
  '/catalogo': 'Catálogo',
  '/carrito': 'Carrito',
  '/checkout': 'Finalizar compra',
  '/confirmacion': 'Confirmación',
  '/pedidos': 'Mis pedidos',
  '/comisiones': 'Mis comisiones',
  '/mi-red': 'Mi red',
  '/admin/productos': 'Gestión de productos',
  '/admin/pedidos': 'Gestión de pedidos',
  '/admin/comisiones': 'Libro de comisiones',
  '/admin/reportes': 'Reportes',
};

const Navbar = () => {
  const { totalItems } = useCart();
  const { isAuthenticated, logout, userEmail, userName, role } = useAuth();
  const { toggleSidebar } = useSidebar(); // estado global del sidebar
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const initials = (userName || userEmail || 'Usuario').split(/\s|@/).filter(Boolean).slice(0, 2).map((part) => part[0]).join('').toUpperCase();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <header className="relative z-20 flex min-h-18 items-center justify-between border-b border-cyan-200/10 bg-slate-950/55 px-3 backdrop-blur-xl sm:px-5 md:px-8">
      <div className="flex items-center gap-2 md:gap-4 min-w-0">
        {/* Botón Toggle: colapsa / expande el Sidebar */}
        <button
          onClick={toggleSidebar}
          aria-label="Alternar menú lateral"
          title="Alternar menú lateral"
          className="icon-button"
        >
          <Menu size={22} aria-hidden="true" />
        </button>
        <div className="min-w-0">
          <p className="hidden text-[0.65rem] font-semibold uppercase tracking-[0.18em] text-cyan-300/70 sm:block">MultiCatálogo</p>
          <h2 className="truncate text-sm font-semibold text-slate-100 sm:text-base">{routeTitles[pathname] ?? 'Panel'}</h2>
        </div>
      </div>
      <div className="flex items-center gap-2 sm:gap-3">
        {isAuthenticated && <Link to="/carrito" className="icon-button relative" aria-label={`Carrito con ${totalItems} productos`} title="Carrito">
          <ShoppingCart size={21} aria-hidden="true" />
          {totalItems > 0 && <span className="absolute -right-1 -top-1 grid h-5 min-w-5 place-items-center rounded-full bg-cyan-300 px-1 text-[0.65rem] font-bold text-slate-950">{totalItems}</span>}
        </Link>}

        {isAuthenticated ? <details className="group relative">
          <summary className="flex list-none items-center gap-2 rounded-2xl border border-cyan-200/15 bg-white/5 p-1.5 pr-2.5 text-left transition hover:border-cyan-300/40 hover:bg-cyan-300/10 [&::-webkit-details-marker]:hidden">
            <span className="grid h-8 w-8 place-items-center rounded-xl bg-cyan-300/10 text-xs font-bold text-cyan-200">{initials}</span>
            <span className="hidden max-w-40 sm:block">
              <span className="block truncate text-xs font-semibold text-slate-100">{userName || userEmail}</span>
              <span className="block text-[0.65rem] uppercase tracking-wider text-slate-500">{role}</span>
            </span>
          </summary>
          <div className="absolute right-0 top-full z-50 mt-2 w-64 rounded-2xl border border-cyan-200/15 bg-slate-950/95 p-2 shadow-2xl backdrop-blur-xl">
            <div className="flex items-center gap-3 border-b border-white/10 px-3 py-3 sm:hidden"><UserRound size={18} /><span className="min-w-0 truncate text-sm">{userEmail}</span></div>
            <button onClick={handleLogout} className="flex w-full items-center gap-2 rounded-xl px-3 py-2.5 text-sm font-semibold text-rose-300 transition hover:bg-rose-500/10"><LogOut size={17} />Cerrar sesión</button>
          </div>
        </details> : <Link to="/login" className="btn-secondary"><LogIn size={17} /> <span className="hidden sm:inline">Iniciar sesión</span></Link>}
      </div>
    </header>
  );
};

export default Navbar;
