export default function Inicio() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Inicio</h1>
      <div className="bg-white p-6 rounded-lg shadow-sm border border-slate-100 max-w-2xl">
        <p className="text-slate-600">
          Bienvenido a la Pokédex. Entra a la sección{' '}
          <span className="font-semibold text-indigo-600">Pokémons</span> en el
          menú lateral para ver los primeros 151 Pokémon consumidos desde la PokeAPI.
        </p>
      </div>
    </div>
  );
}