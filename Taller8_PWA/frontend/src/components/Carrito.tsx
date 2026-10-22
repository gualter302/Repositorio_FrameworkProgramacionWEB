// src/components/Carrito.tsx
import { useCart } from '../context/CartContext';
import { Link } from 'react-router-dom';
import { ArrowRight, ShoppingBag, Trash2 } from 'lucide-react';
import { EmptyState, PageHeader } from './Ui';

const Carrito = () => {
  const { cart, removeFromCart, totalPrice } = useCart();

  return (
    <div className="page page-stack">
      <PageHeader kicker="Selección" title="Tu carrito" description="Revisa los productos elegidos antes de continuar con los datos de envío." />

      {cart.length === 0 ? (
        <EmptyState icon={ShoppingBag} title="Tu carrito está vacío" description="Tu carrito está vacío actualmente. Explora el catálogo y añade los productos que quieras comprar." action={<Link to="/catalogo" className="btn-primary">Ir al catálogo <ArrowRight size={17} /></Link>} />
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Lista de productos */}
          <div className="lg:col-span-2 space-y-4">
            {cart.map((item) => (
              <article key={item.id} className="glass-panel flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-center gap-4">
                  <img src={item.img || '/placeholder.svg'} onError={(event) => { event.currentTarget.onerror = null; event.currentTarget.src = '/placeholder.svg'; }} alt={item.nombre} className="h-18 w-18 rounded-xl border border-white/10 object-cover" />
                  <div className="min-w-0">
                    <h2 className="truncate font-semibold text-white">{item.nombre}</h2>
                    <p className="mt-1 text-sm text-slate-400">Cantidad: {item.cantidad}</p>
                  </div>
                </div>
                <div className="flex items-center gap-6">
                  <p className="mono font-medium text-cyan-300">${(item.precio * item.cantidad).toFixed(2)}</p>
                  <button
                    onClick={() => removeFromCart(item.id)}
                    className="icon-button text-rose-300 hover:border-rose-300/30 hover:bg-rose-500/10 hover:text-rose-200"
                    title="Eliminar producto"
                    aria-label={`Eliminar ${item.nombre}`}
                  >
                    <Trash2 size={18} aria-hidden="true" />
                  </button>
                </div>
              </article>
            ))}
          </div>

          {/* Resumen de pago */}
          <aside className="glass-panel glass-panel--accent h-fit lg:sticky lg:top-4">
            <p className="page-kicker">Pedido</p>
            <h2 className="section-title mb-5">Resumen</h2>
            <div className="mb-4 flex justify-between border-b border-cyan-200/10 pb-4">
              <span className="text-slate-400">Subtotal</span>
              <span className="mono font-semibold text-slate-200">${totalPrice.toFixed(2)}</span>
            </div>
            <div className="flex justify-between items-center mb-6">
              <span className="font-bold text-white">Total a pagar</span>
              <span className="mono text-2xl text-cyan-300">${totalPrice.toFixed(2)}</span>
            </div>
            <Link
              to="/checkout"
              className="btn-primary w-full"
            >
              Continuar <ArrowRight size={17} />
            </Link>
          </aside>
        </div>
      )}
    </div>
  );
};

export default Carrito;
