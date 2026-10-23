import { useEffect, useState } from 'react';
import { actualizarEstadoPedido, cancelarPedido, getPedidosAdmin, type PedidoResumen } from '../services/api';
import { ChevronLeft, ChevronRight, ClipboardList, Search, XCircle } from 'lucide-react';
import { ConfirmDialog, EmptyState, PageHeader, StatusBadge } from './Ui';

const estados = ['procesando', 'preparando', 'enviado', 'entregado'];
const estadosCancelables = ['pendiente', 'confirmado', 'procesando'];

const AdminPedidos = () => {
  const [pedidos, setPedidos] = useState<PedidoResumen[]>([]);
  const [buscar, setBuscar] = useState('');
  const [pagina, setPagina] = useState(1);
  const [total, setTotal] = useState(0);
  const [error, setError] = useState('');
  const [actualizando, setActualizando] = useState<number | null>(null);
  const [pedidoACancelar, setPedidoACancelar] = useState<PedidoResumen | null>(null);

  const cargar = async (paginaActual: number, termino: string) => {
    setError('');
    try {
      const respuesta = await getPedidosAdmin(paginaActual, termino);
      setPedidos(respuesta.datos);
      setTotal(respuesta.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudieron cargar los pedidos');
    }
  };

  useEffect(() => {
    let activo = true;
    getPedidosAdmin(pagina, buscar)
      .then((respuesta) => { if (activo) { setPedidos(respuesta.datos); setTotal(respuesta.total); } })
      .catch((err: Error) => { if (activo) setError(err.message); });
    return () => { activo = false; };
  }, [pagina, buscar]);

  const cambiarEstado = async (id: number, estado: string) => {
    setActualizando(id);
    try { await actualizarEstadoPedido(id, estado); await cargar(pagina, buscar); }
    catch (err) { setError(err instanceof Error ? err.message : 'No se pudo actualizar el estado'); }
    finally { setActualizando(null); }
  };

  const cancelar = async () => {
    if (!pedidoACancelar) return;
    setActualizando(pedidoACancelar.id);
    try { await cancelarPedido(pedidoACancelar.id); await cargar(pagina, buscar); }
    catch (err) { setError(err instanceof Error ? err.message : 'No se pudo cancelar el pedido'); }
    finally { setActualizando(null); setPedidoACancelar(null); }
  };

  const paginas = Math.max(1, Math.ceil(total / 10));
  return (
    <div className="page page-stack">
      <PageHeader kicker="Administración" title="Pedidos" description="Busca compras, actualiza su avance y gestiona cancelaciones." />
      <label className="relative block max-w-xl"><span className="sr-only">Buscar pedidos</span><Search aria-hidden="true" className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} /><input value={buscar} onChange={(e) => { setBuscar(e.target.value); setPagina(1); }} placeholder="Buscar por número, correo o estado" className="form-input form-input--with-icon" /></label>
      {error && <div className="alert-error" role="alert">{error}</div>}
      {!error && pedidos.length === 0 && <EmptyState icon={ClipboardList} title="No se encontraron pedidos" description={buscar ? 'Prueba con otro número, correo o estado.' : 'Las compras registradas aparecerán aquí.'} />}
      <div className="space-y-3">
        {pedidos.map((pedido) => (
          <article key={pedido.id} className="glass-panel">
            <div className="flex flex-wrap justify-between gap-3">
              <div><div className="flex flex-wrap items-center gap-2"><h2 className="mono font-semibold text-white">{pedido.numero}</h2><StatusBadge>{pedido.estado}</StatusBadge></div><p className="mt-2 text-sm text-slate-400">{pedido.usuarioEmail}</p><p className="mt-1 text-xs text-slate-500">{new Date(pedido.fecha).toLocaleString('es-EC')} · <span className="mono text-cyan-300">${pedido.total.toFixed(2)}</span></p></div>
              <div className="flex flex-wrap items-center gap-2">
                <select aria-label={`Estado de ${pedido.numero}`} disabled={actualizando === pedido.id || pedido.estado === 'cancelado' || pedido.estado === 'entregado'} value={pedido.estado} onChange={(e) => void cambiarEstado(pedido.id, e.target.value)} className="form-input min-h-10 w-auto py-2 text-sm">
                  <option value={pedido.estado}>{pedido.estado}</option>
                  {estados.filter((estado) => estado !== pedido.estado).map((estado) => <option key={estado} value={estado}>{estado}</option>)}
                </select>
                {estadosCancelables.includes(pedido.estado) && <button disabled={actualizando === pedido.id} onClick={() => setPedidoACancelar(pedido)} className="btn-danger"><XCircle size={16} />Cancelar</button>}
              </div>
            </div>
          </article>
        ))}
      </div>
      <div className="flex items-center justify-center gap-3"><button disabled={pagina <= 1} onClick={() => setPagina((actual) => actual - 1)} className="btn-secondary"><ChevronLeft size={17} />Anterior</button><span className="mono text-xs text-slate-400">{pagina} / {paginas}</span><button disabled={pagina >= paginas} onClick={() => setPagina((actual) => actual + 1)} className="btn-secondary">Siguiente<ChevronRight size={17} /></button></div>
      <ConfirmDialog
        open={pedidoACancelar !== null}
        title={`Cancelar ${pedidoACancelar?.numero ?? 'pedido'}`}
        description="El pedido quedará cancelado y sus productos volverán al inventario. Esta acción no se puede deshacer."
        confirmLabel="Sí, cancelar pedido"
        onConfirm={() => void cancelar()}
        onCancel={() => setPedidoACancelar(null)}
        busy={actualizando === pedidoACancelar?.id}
      />
    </div>
  );
};

export default AdminPedidos;
