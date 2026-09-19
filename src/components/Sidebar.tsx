import { NavLink } from 'react-router-dom';

interface Props {
  abierto: boolean;
  onCerrar: () => void;
}

// El menú de la izquierda con los enlaces a las páginas.
export default function Sidebar({ abierto, onCerrar }: Props) {
  // Estilo de los enlaces. El de la página actual se pinta de otro color.
  const base = 'block p-3 rounded transition';
  const estilo = ({ isActive }: { isActive: boolean }) =>
    isActive
      ? `${base} bg-indigo-600 text-white`
      : `${base} text-slate-300 hover:bg-slate-800`;

  return (
    <>
      {/* Fondo oscuro que aparece detrás del menú cuando se abre en el celular */}
      {abierto && (
        <div
          onClick={onCerrar}
          className="fixed inset-0 bg-black/50 z-30 md:hidden"
        />
      )}

      {/* En celular el menú se esconde y se desliza; en computadora siempre se ve */}
      <aside
        className={`fixed md:static top-0 left-0 z-40 w-64 h-screen bg-slate-900 text-white flex flex-col transform transition-transform duration-300 ${
          abierto ? 'translate-x-0' : '-translate-x-full'
        } md:translate-x-0`}
      >
        <div className="p-6 text-2xl font-bold border-b border-slate-700">
          Pokédex
        </div>
        <nav className="flex-1 p-4 space-y-2">
          {/* Al tocar un enlace también cerramos el menú (útil en el celular) */}
          <NavLink to="/" end className={estilo} onClick={onCerrar}>
            Inicio
          </NavLink>
          <NavLink to="/pokemons" className={estilo} onClick={onCerrar}>
            Pokémons
          </NavLink>
        </nav>
      </aside>
    </>
  );
}
