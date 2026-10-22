// src/components/Catalogo.tsx
import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Box, LogIn, PackageSearch, ShoppingCart } from 'lucide-react';
import { useCart, type Producto } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';
import { getProductos } from '../services/api';
import { EmptyState, PageHeader } from './Ui';

const Catalogo = () => {
  const { addToCart } = useCart(); // <-- Usamos la función del contexto
  const { isAuthenticated } = useAuth();

  // Los productos ya no están quemados en el código: vienen del backend (GET /api/productos)
  const [productos, setProductos] = useState<Producto[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  // useEffect con [] se ejecuta UNA vez, cuando el componente aparece en pantalla
  useEffect(() => {
    getProductos()
      .then(setProductos)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="page page-stack">
      <PageHeader kicker="Colección" title="Catálogo de productos" description="Explora la selección disponible. Los precios y existencias se verifican nuevamente antes de confirmar tu compra." />
      {!loading && productos.length > 0 && <div className="flex flex-wrap items-center gap-2 text-xs text-slate-500"><span className="status-pill">Copia disponible 24 h</span><span>Última consulta: {new Date(Number(localStorage.getItem('multicatalogo_catalogo_fecha') ?? 0)).toLocaleString('es-EC')}</span></div>}

      {loading && <div className="grid gap-5 sm:grid-cols-2 xl:grid-cols-4">{Array.from({ length: 4 }).map((_, index) => <div key={index} className="glass-panel h-80 animate-pulse"><div className="h-40 rounded-xl bg-white/5" /><div className="mt-5 h-4 w-2/3 rounded bg-white/5" /><div className="mt-3 h-4 w-1/3 rounded bg-white/5" /></div>)}</div>}

      {error && (
        <div className="alert-error" role="alert">
          {error}
        </div>
      )}

      {!loading && !error && productos.length === 0 && <EmptyState icon={PackageSearch} title="No hay productos disponibles" description="El catálogo todavía no tiene productos activos para mostrar." />}

      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-4">
        {productos.map((prod) => (
          <article key={prod.id} className="group flex flex-col overflow-hidden rounded-2xl border border-cyan-200/15 bg-slate-950/55 shadow-2xl backdrop-blur-xl transition duration-200 hover:-translate-y-1 hover:border-cyan-300/40">
            <div className="relative overflow-hidden">
              <img src={prod.img || '/placeholder.svg'} alt={prod.nombre} onError={(event) => { event.currentTarget.onerror = null; event.currentTarget.src = '/placeholder.svg'; }} className="h-48 w-full object-cover transition duration-500 group-hover:scale-[1.03]" />
              <div className="absolute inset-0 bg-gradient-to-t from-slate-950/80 via-transparent to-transparent" />
              <span className="absolute left-3 top-3 rounded-full border border-white/10 bg-slate-950/70 px-3 py-1 text-xs text-slate-200 backdrop-blur">{prod.categoria || 'General'}</span>
              <span className="absolute bottom-3 right-3 flex items-center gap-1.5 rounded-full border border-cyan-200/15 bg-slate-950/80 px-3 py-1 text-xs text-cyan-100"><Box size={13} />{prod.stock} disponibles</span>
            </div>
            <div className="flex flex-1 flex-col p-5">
              <h2 className="text-lg font-semibold text-white">{prod.nombre}</h2>
              {prod.descripcion && <p className="mt-2 line-clamp-2 text-sm leading-relaxed text-slate-400">{prod.descripcion}</p>}
              <p className="mono mb-5 mt-4 text-xl text-cyan-300">${prod.precio.toFixed(2)}</p>

              {isAuthenticated ? <button onClick={() => addToCart(prod)} disabled={prod.stock <= 0} className="btn-primary mt-auto w-full"><ShoppingCart size={17} />{prod.stock > 0 ? 'Añadir al carrito' : 'Sin existencias'}</button> : <Link to="/login" className="btn-secondary mt-auto w-full"><LogIn size={17} />Inicia sesión para comprar</Link>}
            </div>
          </article>
        ))}
      </div>
    </div>
  );
};

export default Catalogo;
