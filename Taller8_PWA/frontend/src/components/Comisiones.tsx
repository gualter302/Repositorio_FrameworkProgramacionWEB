import { useEffect, useMemo, useState } from 'react';
import { getMisComisiones, type Comision } from '../services/api';
import { BadgeDollarSign } from 'lucide-react';
import { EmptyState, PageHeader, StatusBadge } from './Ui';

const Comisiones = () => {
  const [comisiones, setComisiones] = useState<Comision[]>([]);
  const [error, setError] = useState('');
  useEffect(() => {
    let activo = true;
    getMisComisiones().then((datos) => { if (activo) setComisiones(datos); }).catch((err: Error) => { if (activo) setError(err.message); });
    return () => { activo = false; };
  }, []);
  const total = useMemo(() => comisiones.reduce((suma, item) => suma + item.monto, 0), [comisiones]);
  return <div className="page page-stack">
    <PageHeader kicker="Rendimiento" title="Mis comisiones" description="Consulta cuánto generó cada pedido, su nivel dentro de la red y el estado del pago." />
    {error && <div className="alert-error" role="alert">{error}</div>}
    <div className="metric-card glass-panel--accent max-w-xl"><BadgeDollarSign className="mb-5 text-emerald-300" size={22} /><p className="metric-label">Total histórico</p><p className="metric-value text-emerald-300">${total.toFixed(2)}</p></div>
    {!error && comisiones.length === 0 ? <EmptyState icon={BadgeDollarSign} title="Aún no hay comisiones" description="Las comisiones generadas por las compras de tu red aparecerán aquí." /> : <div className="data-shell"><table className="data-table"><thead><tr><th>Pedido</th><th>Nivel</th><th>Porcentaje</th><th>Monto</th><th>Estado</th></tr></thead><tbody>{comisiones.map((item) => <tr key={item.id}><td className="mono text-slate-300">{item.numeroPedido}</td><td>{item.nivel}</td><td>{(item.porcentaje * 100).toFixed(0)}%</td><td className="mono text-emerald-300">${item.monto.toFixed(2)}</td><td><StatusBadge>{item.estado}</StatusBadge></td></tr>)}</tbody></table></div>}
  </div>;
};
export default Comisiones;
