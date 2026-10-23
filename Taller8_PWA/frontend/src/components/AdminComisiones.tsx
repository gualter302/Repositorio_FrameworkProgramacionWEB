import { useEffect, useState } from 'react';
import { BadgeDollarSign, HandCoins } from 'lucide-react';
import { getComisionesAdmin, liquidarComisiones, type Comision } from '../services/api';
import { ConfirmDialog, EmptyState, PageHeader, StatusBadge } from './Ui';

const AdminComisiones = () => {
  const [comisiones, setComisiones] = useState<Comision[]>([]);
  const [seleccionadas, setSeleccionadas] = useState<number[]>([]);
  const [filtro, setFiltro] = useState('pendiente');
  const [error, setError] = useState('');
  const [mensaje, setMensaje] = useState('');
  const [confirmando, setConfirmando] = useState(false);
  const [procesando, setProcesando] = useState(false);
  const cargar = async () => { try { setError(''); setComisiones(await getComisionesAdmin(filtro)); } catch (err) { setError(err instanceof Error ? err.message : 'No se pudieron cargar las comisiones'); } };
  useEffect(() => {
    let activo = true;
    getComisionesAdmin(filtro).then((datos) => { if (activo) setComisiones(datos); }).catch((err: Error) => { if (activo) setError(err.message); });
    return () => { activo = false; };
  }, [filtro]);
  const referidoSeleccionado = comisiones.find((item) => seleccionadas.includes(item.id))?.referidoId;
  const alternar = (comision: Comision) => setSeleccionadas((actual) => {
    if (actual.includes(comision.id)) return actual.filter((item) => item !== comision.id);
    const referidoActual = comisiones.find((item) => actual.includes(item.id))?.referidoId;
    return referidoActual === undefined || referidoActual === comision.referidoId ? [...actual, comision.id] : actual;
  });
  const liquidar = async () => {
    setProcesando(true);
    setError('');
    setMensaje('');
    try {
      const liquidacion = await liquidarComisiones(seleccionadas);
      setSeleccionadas([]);
      setConfirmando(false);
      setMensaje(`Liquidación #${liquidacion.id} creada por $${liquidacion.monto.toFixed(2)}.`);
      await cargar();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo crear la liquidación');
      setConfirmando(false);
    } finally {
      setProcesando(false);
    }
  };

  const totalSeleccionado = comisiones
    .filter((item) => seleccionadas.includes(item.id))
    .reduce((total, item) => total + item.monto, 0);

  return <div className="page page-stack">
    <PageHeader kicker="Administración" title="Libro de comisiones" description="Filtra movimientos pendientes o pagados y agrupa los que correspondan al mismo referido." />
    <div className="glass-panel flex flex-col gap-3 sm:flex-row sm:items-end">
      <label className="w-full sm:max-w-xs"><span className="field-label">Estado</span><select value={filtro} onChange={(e) => { setFiltro(e.target.value); setSeleccionadas([]); setMensaje(''); setError(''); }} className="form-input"><option value="pendiente">Pendientes</option><option value="pagada">Pagadas</option><option value="">Todas</option></select></label>
      <button disabled={!seleccionadas.length || filtro === 'pagada'} onClick={() => setConfirmando(true)} className="btn-primary"><HandCoins size={17} />Crear liquidación {seleccionadas.length > 0 && `(${seleccionadas.length})`}</button>
    </div>
    {mensaje && <div className="alert-success">{mensaje}</div>}
    {error && <div className="alert-error" role="alert">{error}</div>}
    {!error && comisiones.length === 0 ? <EmptyState icon={BadgeDollarSign} title="No hay comisiones en este filtro" description="Cambia el estado seleccionado o espera a que se generen nuevas comisiones." /> : <div className="data-shell"><table className="data-table"><thead><tr><th>Seleccionar</th><th>Pedido</th><th>Referido</th><th>Monto</th><th>Estado</th></tr></thead><tbody>{comisiones.map((item) => {
      const perteneceAOtroReferido = referidoSeleccionado !== undefined && referidoSeleccionado !== item.referidoId;
      return <tr key={item.id}><td><input aria-label={`Seleccionar comisión ${item.id}`} className="h-4 w-4 accent-cyan-400" type="checkbox" disabled={item.estado !== 'pendiente' || perteneceAOtroReferido} checked={seleccionadas.includes(item.id)} onChange={() => alternar(item)} /></td><td className="mono text-slate-300">{item.numeroPedido}</td><td>{item.referido}</td><td className="mono text-emerald-300">${item.monto.toFixed(2)}</td><td><StatusBadge>{item.estado}</StatusBadge></td></tr>;
    })}</tbody></table></div>}
    <ConfirmDialog
      open={confirmando}
      title="Crear liquidación"
      description={`Se marcarán ${seleccionadas.length} comisiones como pagadas por un total de $${totalSeleccionado.toFixed(2)}. Todas pertenecen al mismo referido.`}
      confirmLabel="Confirmar liquidación"
      onConfirm={() => void liquidar()}
      onCancel={() => setConfirmando(false)}
      busy={procesando}
      tone="primary"
    />
  </div>;
};
export default AdminComisiones;
