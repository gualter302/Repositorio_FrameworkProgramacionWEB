import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import Inicio from './pages/Inicio';
import Pokemons from './pages/Pokemons';

// Aquí decidimos qué página mostrar según la dirección en la que esté el usuario.
function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Layout es el marco (menú + barra de arriba) que rodea a todas las páginas */}
        <Route path="/" element={<Layout />}>
          <Route index element={<Inicio />} />            {/* Página de inicio */}
          <Route path="pokemons" element={<Pokemons />} /> {/* Página de Pokémon */}
        </Route>

        {/* Si escriben una dirección que no existe, los mandamos al inicio */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
