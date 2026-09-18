import { NavLink } from 'react-router-dom';

export default function Sidebar() {
  const base = 'block p-3 rounded transition';
  const estilo = ({ isActive }: { isActive: boolean }) =>
    isActive
      ? `${base} bg-indigo-600 text-white`
      : `${base} text-slate-300 hover:bg-slate-800`;

  return (
    <aside className="w-64 h-screen bg-slate-900 text-white flex flex-col">
      <div className="p-6 text-2xl font-bold border-b border-slate-700">
        Pokédex
      </div>
      <nav className="flex-1 p-4 space-y-2">
        <NavLink to="/" end className={estilo}>Inicio</NavLink>
        <NavLink to="/pokemons" className={estilo}>Pokémons</NavLink>
      </nav>
    </aside>
  );
}