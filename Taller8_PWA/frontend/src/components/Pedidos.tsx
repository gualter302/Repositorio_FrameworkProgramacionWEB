import { useEffect, useState } from 'react';
import { getMisPedidos, type PedidoResumen } from '../services/api';

const Pedidos = () => {
  const [pedidos, setPedidos] = useState<PedidoResumen[]>([]);
  const [error, setError] = useState('');

  useEffect(() => {
    let activo = true;
    getMisPedidos()
      .then((datos) => { if (activo) setPedidos(datos); })
      .catch((err: Error) => { if (activo) setError(err.message); });
    return () => { activo = false; };
  }, []);

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-2xl font-bold text-slate-800">Mis pedidos</h1>
        <p className="text-slate-500 mt-1">Consulta el estado y detalle de tus compras.</p>
      </div>
      {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg">{error}</div>}
      {pedidos.length === 0 && !error && <div className="bg-white p-6 rounded-lg border text-slate-500">Todavía no tienes pedidos.</div>}
      {pedidos.map((pedido) => (
        <article key={pedido.id} className="bg-white border border-slate-200 rounded-lg p-5">
          <div className="flex flex-wrap justify-between gap-3 mb-4">
            <div><h2 className="font-bold text-slate-800">{pedido.numero}</h2><p className="text-sm text-slate-500">{new Date(pedido.fecha).toLocaleString('es-EC')}</p></div>
            <div className="text-right"><span className="inline-block rounded-full bg-indigo-50 text-indigo-700 px-3 py-1 text-sm font-semibold">{pedido.estado}</span><p className="font-bold mt-1">${pedido.total.toFixed(2)}</p></div>
          </div>
          <div className="space-y-2 border-t pt-3">
            {pedido.items.map((item) => <div key={item.productoId} className="flex justify-between text-sm"><span>{item.nombre} × {item.cantidad}</span><span>${item.subtotal.toFixed(2)}</span></div>)}
          </div>
        </article>
      ))}
    </div>
  );
};

export default Pedidos;
