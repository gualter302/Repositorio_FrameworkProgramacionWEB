import { useEffect, useState } from 'react';
import { getMisPedidos, type PedidoResumen } from '../services/api';
import { PackageOpen } from 'lucide-react';
import { EmptyState, PageHeader, StatusBadge } from './Ui';

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
    <div className="page page-stack">
      <PageHeader kicker="Historial" title="Mis pedidos" description="Consulta el estado, la fecha y el detalle de cada compra registrada." />
      {error && <div className="alert-error" role="alert">{error}</div>}
      {pedidos.length === 0 && !error && <EmptyState icon={PackageOpen} title="Todavía no tienes pedidos" description="Cuando completes una compra aparecerá aquí con su estado y detalle." />}
      {pedidos.map((pedido) => (
        <article key={pedido.id} className="glass-panel transition hover:border-cyan-300/30">
          <div className="flex flex-wrap justify-between gap-3 mb-4">
            <div><p className="page-kicker mb-1">Pedido</p><h2 className="mono font-semibold text-white">{pedido.numero}</h2><p className="mt-1 text-sm text-slate-500">{new Date(pedido.fecha).toLocaleString('es-EC')}</p></div>
            <div className="text-right"><StatusBadge>{pedido.estado}</StatusBadge><p className="mono mt-2 text-cyan-300">${pedido.total.toFixed(2)}</p></div>
          </div>
          <div className="space-y-2 border-t border-cyan-200/10 pt-4">
            {pedido.items.map((item) => <div key={item.productoId} className="flex justify-between gap-4 text-sm"><span className="text-slate-400">{item.nombre} × {item.cantidad}</span><span className="mono text-slate-200">${item.subtotal.toFixed(2)}</span></div>)}
          </div>
        </article>
      ))}
    </div>
  );
};

export default Pedidos;
