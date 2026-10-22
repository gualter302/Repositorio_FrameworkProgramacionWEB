import { Link, useLocation } from 'react-router-dom';
import type { PedidoConfirmado } from '../services/api';

const Confirmacion = () => {
  const location = useLocation();
  const pedido = (location.state as { pedido?: PedidoConfirmado } | null)?.pedido;

  if (!pedido) {
    return (
      <div className="bg-white p-10 rounded-lg border border-slate-200 text-center">
        <div className="text-5xl mb-4">✅</div>
        <h1 className="text-2xl font-bold text-slate-800 mb-2">No hay un pedido reciente</h1>
        <p className="text-slate-500 mb-6">Realiza una compra para ver aquí su confirmación.</p>
        <Link to="/catalogo" className="inline-block bg-indigo-600 text-white font-semibold px-6 py-3 rounded-lg">
          Ir al Catálogo
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="bg-green-600 text-white p-8 text-center">
          <div className="text-5xl mb-3">✓</div>
          <h1 className="text-2xl font-bold">¡Pedido confirmado!</h1>
          <p className="text-green-100 mt-1">Número: <span className="font-mono font-bold">{pedido.numero}</span></p>
        </div>
        <div className="p-8">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6 text-sm">
            <div><p className="text-slate-500">Cliente</p><p className="font-semibold">{pedido.cliente.nombre}</p></div>
            <div><p className="text-slate-500">Fecha</p><p className="font-semibold">{new Date(pedido.fecha).toLocaleString('es-EC')}</p></div>
            <div><p className="text-slate-500">Correo</p><p className="font-semibold">{pedido.cliente.email}</p></div>
            <div><p className="text-slate-500">Ciudad</p><p className="font-semibold">{pedido.cliente.ciudad}</p></div>
          </div>
          <h2 className="font-bold text-slate-800 mb-3">Detalle del pedido</h2>
          <div className="space-y-2 border-b border-slate-100 pb-4 mb-4">
            {pedido.items.map((item) => (
              <div key={item.productoId} className="flex justify-between text-sm">
                <span className="text-slate-600">{item.nombre} × {item.cantidad}</span>
                <span className="font-semibold">${item.subtotal.toFixed(2)}</span>
              </div>
            ))}
          </div>
          <div className="flex justify-between items-center mb-8">
            <span className="font-bold">Total pagado</span>
            <span className="text-2xl font-bold text-green-600">${pedido.total.toFixed(2)}</span>
          </div>
          <Link to="/catalogo" className="block text-center bg-indigo-600 text-white font-semibold py-3 rounded-lg">
            Seguir Comprando
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Confirmacion;
