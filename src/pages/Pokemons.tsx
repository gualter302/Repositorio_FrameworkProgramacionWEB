import { useEffect, useMemo, useState } from 'react';
import TarjetaPokemon from '../components/TarjetaPokemon';

// La forma que tiene cada Pokémon que nos llega de la API.
interface Pokemon {
  id: number;
  name: string;
  sprites: { front_default: string };
  types: { type: { name: string } }[];
}

// Las letras del abecedario para los botones de filtro.
const LETRAS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('');

export default function Pokemons() {
  // Aquí guardamos la lista de Pokémon y si todavía se están cargando.
  const [pokemones, setPokemones] = useState<Pokemon[]>([]);
  const [cargando, setCargando] = useState<boolean>(true);

  // Lo que el usuario elige en cada filtro.
  const [filtroTipo, setFiltroTipo] = useState<string>('');
  const [filtroLetra, setFiltroLetra] = useState<string>('');
  const [orden, setOrden] = useState<'asc' | 'desc'>('asc');

  // Al abrir la página, le pedimos los Pokémon a la API.
  useEffect(() => {
    const obtenerPokemones = async () => {
      try {
        // Primero pedimos la lista (solo trae el nombre y un enlace de cada uno).
        const respuesta = await fetch('https://pokeapi.co/api/v2/pokemon?limit=151');
        const datos = await respuesta.json();

        // Luego pedimos los datos completos de cada Pokémon (foto, tipos, etc.).
        const detalles: Pokemon[] = await Promise.all(
          datos.results.map((p: { url: string }) =>
            fetch(p.url).then((r) => r.json())
          )
        );
        setPokemones(detalles);
      } catch (error) {
        console.error('Error al consumir la PokeAPI', error);
      } finally {
        setCargando(false);
      }
    };
    obtenerPokemones();
  }, []); // Se hace una sola vez, al entrar a la página.

  // Sacamos la lista de tipos que existen, para armar el menú de tipos.
  const tipos = useMemo(() => {
    const set = new Set<string>();
    pokemones.forEach((p) => p.types.forEach((t) => set.add(t.type.name)));
    return Array.from(set).sort();
  }, [pokemones]);

  // Aplicamos los filtros y ordenamos la lista según lo que eligió el usuario.
  const filtrados = useMemo(() => {
    return pokemones
      // Dejamos solo los del tipo elegido (o todos si no eligió ninguno).
      .filter((p) =>
        filtroTipo ? p.types.some((t) => t.type.name === filtroTipo) : true
      )
      // Dejamos solo los que empiezan por la letra elegida.
      .filter((p) =>
        filtroLetra
          ? p.name.toLowerCase().startsWith(filtroLetra.toLowerCase())
          : true
      )
      // Ordenamos por nombre (de A a Z, o de Z a A).
      .sort((a, b) =>
        orden === 'asc'
          ? a.name.localeCompare(b.name)
          : b.name.localeCompare(a.name)
      );
  }, [pokemones, filtroTipo, filtroLetra, orden]);

  // Mientras se cargan, mostramos un mensaje.
  if (cargando) {
    return <p className="text-slate-500">Cargando Pokémon...</p>;
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Pokémons</h1>

      {/* Zona de filtros */}
      <div className="flex flex-col gap-4 mb-6">
        <div className="flex flex-wrap items-center gap-3">
          {/* Filtro por tipo */}
          <select
            value={filtroTipo}
            onChange={(e) => setFiltroTipo(e.target.value)}
            className="border border-slate-300 rounded-lg px-3 py-2 text-sm bg-white text-slate-700 capitalize"
          >
            <option value="">Todos los tipos</option>
            {tipos.map((tipo) => (
              <option key={tipo} value={tipo}>
                {tipo}
              </option>
            ))}
          </select>

          {/* Orden por nombre */}
          <select
            value={orden}
            onChange={(e) => setOrden(e.target.value as 'asc' | 'desc')}
            className="border border-slate-300 rounded-lg px-3 py-2 text-sm bg-white text-slate-700"
          >
            <option value="asc">Nombre: A → Z</option>
            <option value="desc">Nombre: Z → A</option>
          </select>

          {/* Botón para quitar los filtros, solo aparece si hay alguno puesto */}
          {(filtroTipo || filtroLetra) && (
            <button
              onClick={() => {
                setFiltroTipo('');
                setFiltroLetra('');
              }}
              className="text-sm text-indigo-600 hover:underline"
            >
              Limpiar filtros
            </button>
          )}
        </div>

        {/* Filtro por letra inicial */}
        <div className="flex flex-wrap gap-1">
          <button
            onClick={() => setFiltroLetra('')}
            className={`px-3 h-8 rounded text-sm font-medium transition ${
              filtroLetra === ''
                ? 'bg-indigo-600 text-white'
                : 'bg-white border border-slate-200 text-slate-600 hover:bg-slate-100'
            }`}
          >
            Todas
          </button>
          {LETRAS.map((letra) => (
            <button
              key={letra}
              onClick={() => setFiltroLetra(letra)}
              className={`w-8 h-8 rounded text-sm font-medium transition ${
                filtroLetra === letra
                  ? 'bg-indigo-600 text-white'
                  : 'bg-white border border-slate-200 text-slate-600 hover:bg-slate-100'
              }`}
            >
              {letra}
            </button>
          ))}
        </div>
      </div>

      {/* La lista de tarjetas (o un aviso si no hay coincidencias) */}
      {filtrados.length === 0 ? (
        <p className="text-slate-500">
          No hay Pokémon que coincidan con el filtro.
        </p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
          {filtrados.map((pokemon) => (
            <TarjetaPokemon
              key={pokemon.id}
              id={pokemon.id}
              nombre={pokemon.name}
              imagen={pokemon.sprites.front_default}
              tipos={pokemon.types.map((t) => t.type.name)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
