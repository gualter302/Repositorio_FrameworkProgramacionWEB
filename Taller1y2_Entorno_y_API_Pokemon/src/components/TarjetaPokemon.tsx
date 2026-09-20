interface Props {
  id: number;
  nombre: string;
  imagen: string;
  tipos: string[];
}

// Tarjeta de un Pokémon: muestra su número, su foto, su nombre y sus tipos.
export default function TarjetaPokemon({ id, nombre, imagen, tipos }: Props) {
  return (
    <div className="bg-white rounded-lg overflow-hidden border border-slate-200 shadow-sm flex flex-col">
      <div className="bg-slate-50 flex items-center justify-center py-2">
        <img src={imagen} alt={nombre} className="w-32 h-32 object-contain" />
      </div>
      <div className="p-4 flex flex-col flex-1">
        {/* El número con ceros adelante, por ejemplo #001 */}
        <span className="text-xs text-slate-400 font-semibold">
          #{String(id).padStart(3, '0')}
        </span>
        <h3 className="font-semibold text-slate-700 capitalize">{nombre}</h3>

        {/* Mostramos cada tipo dentro de una etiqueta */}
        <div className="flex flex-wrap gap-1 mt-2">
          {tipos.map((tipo) => (
            <span
              key={tipo}
              className="text-xs bg-slate-100 text-slate-600 px-2 py-0.5 rounded-full capitalize"
            >
              {tipo}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}
