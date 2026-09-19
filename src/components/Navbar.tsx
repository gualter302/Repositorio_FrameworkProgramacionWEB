interface Props {
  onAbrirMenu: () => void;
}

// La barra de arriba. En el celular trae el botón para abrir el menú.
export default function Navbar({ onAbrirMenu }: Props) {
  return (
    <header className="h-16 bg-white border-b border-slate-200 flex items-center gap-3 px-4 md:px-8">
      {/* Este botón (☰) solo se ve en el celular */}
      <button
        onClick={onAbrirMenu}
        aria-label="Abrir menú"
        className="md:hidden p-2 rounded hover:bg-slate-100 transition text-slate-600 text-xl"
      >
        ☰
      </button>
      <h2 className="text-slate-600 font-medium text-base md:text-lg">
        Consumo de PokeAPI
      </h2>
    </header>
  );
}
