import { useEffect, useState } from 'react';
import { actualizarEstadoPedido, cancelarPedido, getPedidosAdmin, type PedidoResumen } from '../services/api';

const estados = ['procesando', 'preparando', 'enviado', 'entregado'];

const AdminPedidos = () => {
  const [pedidos, setPedidos] = useState<PedidoResumen[]>([]);
  const [buscar, setBuscar] = useState('');
  const [pagina, setPagina] = useState(1);
  const [total, setTotal] = useState(0);
  const [error, setError] = useState('');
  const [actualizando, setActualizando] = useState<number | null>(null);

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

  const cancelar = async (id: number) => {
    if (!window.confirm('¿Cancelar este pedido y devolver su stock?')) return;
    setActualizando(id);
    try { await cancelarPedido(id); await cargar(pagina, buscar); }
    catch (err) { setError(err instanceof Error ? err.message : 'No se pudo cancelar el pedido'); }
    finally { setActualizando(null); }
  };

  const paginas = Math.max(1, Math.ceil(total / 10));
  return (
    <div className="space-y-5">
      <div><h1 className="text-2xl font-bold text-slate-800">Administrar pedidos</h1><p className="text-slate-500 mt-1">Gestiona estados, búsqueda y cancelaciones.</p></div>
      <input value={buscar} onChange={(e) => { setBuscar(e.target.value); setPagina(1); }} placeholder="Buscar por número, correo o estado" className="w-full max-w-lg border rounded-lg p-3" />
      {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg">{error}</div>}
      <div className="space-y-3">
        {pedidos.map((pedido) => (
          <article key={pedido.id} className="bg-white border rounded-lg p-4">
            <div className="flex flex-wrap justify-between gap-3">
              <div><h2 className="font-bold">{pedido.numero} · {pedido.usuarioEmail}</h2><p className="text-sm text-slate-500">{new Date(pedido.fecha).toLocaleString('es-EC')} · ${pedido.total.toFixed(2)}</p></div>
              <div className="flex items-center gap-2">
                <select disabled={actualizando === pedido.id || pedido.estado === 'cancelado' || pedido.estado === 'entregado'} value={pedido.estado} onChange={(e) => void cambiarEstado(pedido.id, e.target.value)} className="border rounded p-2 text-sm">
                  <option value={pedido.estado}>{pedido.estado}</option>
                  {estados.filter((estado) => estado !== pedido.estado).map((estado) => <option key={estado} value={estado}>{estado}</option>)}
                </select>
                {pedido.estado !== 'cancelado' && pedido.estado !== 'entregado' && <button disabled={actualizando === pedido.id} onClick={() => void cancelar(pedido.id)} className="text-red-600 hover:underline text-sm">Cancelar</button>}
              </div>
            </div>
          </article>
        ))}
      </div>
      <div className="flex items-center gap-3"><button disabled={pagina <= 1} onClick={() => setPagina((actual) => actual - 1)} className="border px-3 py-2 rounded disabled:opacity-50">Anterior</button><span className="text-sm">Página {pagina} de {paginas}</span><button disabled={pagina >= paginas} onClick={() => setPagina((actual) => actual + 1)} className="border px-3 py-2 rounded disabled:opacity-50">Siguiente</button></div>
    </div>
  );
};

export default AdminPedidos;
