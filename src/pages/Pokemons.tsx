import { useEffect, useState } from 'react';
import TarjetaPokemon from '../components/TarjetaPokemon';

interface Pokemon {
  id: number;
  name: string;
  sprites: { front_default: string };
  types: { type: { name: string } }[];
}

export default function Pokemons() {
  const [pokemones, setPokemones] = useState<Pokemon[]>([]);
  const [cargando, setCargando] = useState<boolean>(true);

  useEffect(() => {
    const obtenerPokemones = async () => {
      try {
        const respuesta = await fetch('https://pokeapi.co/api/v2/pokemon?limit=151');
        const datos = await respuesta.json();
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
  }, []);

  if (cargando) {
    return <p className="text-slate-500">Cargando Pokémon...</p>;
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Pokémons</h1>
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
        {pokemones.map((pokemon) => (
          <TarjetaPokemon
            key={pokemon.id}
            id={pokemon.id}
            nombre={pokemon.name}
            imagen={pokemon.sprites.front_default}
            tipos={pokemon.types.map((t) => t.type.name)}
          />
        ))}
      </div>
    </div>
  );
}