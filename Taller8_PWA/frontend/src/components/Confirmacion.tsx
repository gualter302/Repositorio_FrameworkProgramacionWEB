import { Link, useLocation } from 'react-router-dom';
import type { PedidoConfirmado } from '../services/api';
import { ArrowRight, Check, PackageCheck } from 'lucide-react';
import { EmptyState } from './Ui';

const Confirmacion = () => {
  const location = useLocation();
  const pedido = (location.state as { pedido?: PedidoConfirmado } | null)?.pedido;

  if (!pedido) {
    return (
      <div className="page"><EmptyState icon={PackageCheck} title="No hay un pedido reciente" description="Realiza una compra para ver aquí su confirmación." action={<Link to="/catalogo" className="btn-primary">Ir al catálogo <ArrowRight size={17} /></Link>} /></div>
    );
  }

  return (
    <div className="page max-w-3xl">
      <div className="overflow-hidden rounded-[2rem] border border-cyan-200/20 bg-slate-950/65 shadow-2xl backdrop-blur-xl">
        <div className="border-b border-emerald-300/20 bg-emerald-400/10 p-8 text-center">
          <span className="mx-auto mb-4 grid h-16 w-16 place-items-center rounded-2xl border border-emerald-300/30 bg-emerald-400/10 text-emerald-300"><Check size={34} /></span>
          <p className="page-kicker text-emerald-300">Registro completo</p>
          <h1 className="page-title text-2xl">Pedido confirmado</h1>
          <p className="mt-2 text-emerald-100/70">Número: <span className="mono font-bold text-emerald-200">{pedido.numero}</span></p>
        </div>
        <div className="p-8">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6 text-sm">
            <div><p className="text-slate-500">Cliente</p><p className="font-semibold text-slate-200">{pedido.cliente.nombre}</p></div>
            <div><p className="text-slate-500">Fecha</p><p className="font-semibold text-slate-200">{new Date(pedido.fecha).toLocaleString('es-EC')}</p></div>
            <div><p className="text-slate-500">Correo</p><p className="font-semibold text-slate-200">{pedido.cliente.email}</p></div>
            <div><p className="text-slate-500">Ciudad</p><p className="font-semibold text-slate-200">{pedido.cliente.ciudad}</p></div>
          </div>
          <h2 className="section-title mb-3">Detalle del pedido</h2>
          <div className="mb-4 space-y-2 border-b border-cyan-200/10 pb-4">
            {pedido.items.map((item) => (
              <div key={item.productoId} className="flex justify-between text-sm">
                <span className="text-slate-400">{item.nombre} × {item.cantidad}</span>
                <span className="mono text-slate-200">${item.subtotal.toFixed(2)}</span>
              </div>
            ))}
          </div>
          <div className="flex justify-between items-center mb-8">
            <span className="font-bold text-white">Total registrado</span>
            <span className="mono text-2xl text-emerald-300">${pedido.total.toFixed(2)}</span>
          </div>
          <Link to="/catalogo" className="btn-primary w-full">
            Seguir comprando <ArrowRight size={17} />
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Confirmacion;
