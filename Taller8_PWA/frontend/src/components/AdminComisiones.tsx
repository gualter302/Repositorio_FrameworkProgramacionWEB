import { useEffect, useState } from 'react';
import { BadgeDollarSign, HandCoins } from 'lucide-react';
import { getComisionesAdmin, liquidarComisiones, type Comision } from '../services/api';
import { EmptyState, PageHeader, StatusBadge } from './Ui';

const AdminComisiones = () => {
  const [comisiones, setComisiones] = useState<Comision[]>([]);
  const [seleccionadas, setSeleccionadas] = useState<number[]>([]);
  const [filtro, setFiltro] = useState('pendiente');
  const [error, setError] = useState('');
  const cargar = async () => { try { setComisiones(await getComisionesAdmin(filtro)); } catch (err) { setError(err instanceof Error ? err.message : 'No se pudieron cargar las comisiones'); } };
  useEffect(() => {
    let activo = true;
    getComisionesAdmin(filtro).then((datos) => { if (activo) setComisiones(datos); }).catch((err: Error) => { if (activo) setError(err.message); });
    return () => { activo = false; };
  }, [filtro]);
  const alternar = (id: number) => setSeleccionadas((actual) => actual.includes(id) ? actual.filter((item) => item !== id) : [...actual, id]);
  const liquidar = async () => { try { await liquidarComisiones(seleccionadas); setSeleccionadas([]); await cargar(); } catch (err) { setError(err instanceof Error ? err.message : 'No se pudo crear la liquidación'); } };

  return <div className="page page-stack">
    <PageHeader kicker="Administración" title="Libro de comisiones" description="Filtra movimientos pendientes o pagados y agrupa los que correspondan al mismo referido." />
    <div className="glass-panel flex flex-col gap-3 sm:flex-row sm:items-end">
      <label className="w-full sm:max-w-xs"><span className="field-label">Estado</span><select value={filtro} onChange={(e) => { setFiltro(e.target.value); setSeleccionadas([]); }} className="form-input"><option value="pendiente">Pendientes</option><option value="pagada">Pagadas</option><option value="">Todas</option></select></label>
      <button disabled={!seleccionadas.length || filtro === 'pagada'} onClick={() => void liquidar()} className="btn-primary"><HandCoins size={17} />Crear liquidación {seleccionadas.length > 0 && `(${seleccionadas.length})`}</button>
    </div>
    {error && <div className="alert-error" role="alert">{error}</div>}
    {!error && comisiones.length === 0 ? <EmptyState icon={BadgeDollarSign} title="No hay comisiones en este filtro" description="Cambia el estado seleccionado o espera a que se generen nuevas comisiones." /> : <div className="data-shell"><table className="data-table"><thead><tr><th>Seleccionar</th><th>Pedido</th><th>Referido</th><th>Monto</th><th>Estado</th></tr></thead><tbody>{comisiones.map((item) => <tr key={item.id}><td><input aria-label={`Seleccionar comisión ${item.id}`} className="h-4 w-4 accent-cyan-400" type="checkbox" disabled={item.estado !== 'pendiente'} checked={seleccionadas.includes(item.id)} onChange={() => alternar(item.id)} /></td><td className="mono text-slate-300">{item.numeroPedido}</td><td>{item.referido}</td><td className="mono text-emerald-300">${item.monto.toFixed(2)}</td><td><StatusBadge>{item.estado}</StatusBadge></td></tr>)}</tbody></table></div>}
  </div>;
};
export default AdminComisiones;
