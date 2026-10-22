import { useEffect, useState } from 'react';
import { BarChart3 } from 'lucide-react';
import { getReporteGeneral, type ReporteGeneral } from '../services/api';
import { EmptyState, PageHeader } from './Ui';

const AdminReportes = () => {
  const [reporte, setReporte] = useState<ReporteGeneral | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    let activo = true;
    getReporteGeneral()
      .then((datos) => { if (activo) setReporte(datos); })
      .catch((err: Error) => { if (activo) setError(err.message); });
    return () => { activo = false; };
  }, []);

  return (
    <div className="page page-stack">
      <PageHeader kicker="Administración" title="Reportes generales" description="Indicadores calculados directamente desde las ventas, pedidos y comisiones guardadas." />
      {error && <div className="alert-error" role="alert">{error}</div>}
      {!reporte && !error && <div className="metric-grid">{Array.from({ length: 3 }).map((_, index) => <div className="metric-card animate-pulse" key={index}><div className="h-3 w-1/2 rounded bg-white/5" /><div className="mt-5 h-8 w-2/3 rounded bg-white/5" /></div>)}</div>}
      {reporte && <>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-6">
          {[
            ['Ventas', `$${reporte.totalVentas.toFixed(2)}`],
            ['Pedidos', reporte.totalPedidos],
            ['Ticket promedio', `$${reporte.ticketPromedio.toFixed(2)}`],
            ['Comisiones', `$${reporte.totalComisiones.toFixed(2)}`],
            ['Pendiente de pago', `$${reporte.pendientePagar.toFixed(2)}`],
          ].map(([label, value], index) => <div key={label} className={`metric-card ${index === 0 ? 'glass-panel--accent sm:col-span-2 xl:col-span-2' : ''}`}><p className="metric-label">{label}</p><p className="metric-value">{value}</p></div>)}
        </div>
        {reporte.productosPrincipales.length === 0 && reporte.referidosPrincipales.length === 0 ? <EmptyState icon={BarChart3} title="Aún no hay información suficiente" description="Los rankings aparecerán cuando se registren ventas y comisiones." /> : <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <section className="data-shell"><h2 className="section-title border-b border-cyan-200/10 p-4">Productos principales</h2><table className="data-table min-w-full"><thead><tr><th>Producto</th><th className="text-right">Unidades</th><th className="text-right">Ventas</th></tr></thead><tbody>{reporte.productosPrincipales.map((item) => <tr key={item.nombre}><td>{item.nombre}</td><td className="text-right">{item.cantidad}</td><td className="mono text-right text-cyan-300">${item.ventas.toFixed(2)}</td></tr>)}</tbody></table></section>
          <section className="data-shell"><h2 className="section-title border-b border-cyan-200/10 p-4">Referidos principales</h2><table className="data-table min-w-full"><thead><tr><th>Referido</th><th className="text-right">Ventas</th><th className="text-right">Comisiones</th></tr></thead><tbody>{reporte.referidosPrincipales.map((item) => <tr key={item.nombre}><td>{item.nombre}</td><td className="mono text-right text-cyan-300">${item.ventas.toFixed(2)}</td><td className="mono text-right text-emerald-300">${item.comisiones.toFixed(2)}</td></tr>)}</tbody></table></section>
        </div>}
      </>}
    </div>
  );
};

export default AdminReportes;
