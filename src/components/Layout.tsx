import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Navbar from './Navbar';

// El marco de la app: el menú de la izquierda, la barra de arriba y el contenido.
export default function Layout() {
  // Guarda si el menú lateral está abierto. Esto solo importa en el celular.
  const [menuAbierto, setMenuAbierto] = useState(false);

  return (
    <div className="flex h-screen bg-slate-50">
      <Sidebar abierto={menuAbierto} onCerrar={() => setMenuAbierto(false)} />

      <div className="flex-1 flex flex-col overflow-hidden">
        <Navbar onAbrirMenu={() => setMenuAbierto(true)} />

        {/* Aquí adentro se muestra la página en la que estés (Inicio o Pokémon) */}
        <main className="flex-1 overflow-y-auto p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
