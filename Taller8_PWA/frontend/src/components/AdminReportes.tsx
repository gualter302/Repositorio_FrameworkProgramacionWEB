import { useEffect, useState } from 'react';
import { getReporteGeneral, type ReporteGeneral } from '../services/api';

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
    <div className="space-y-6">
      <div><h1 className="text-2xl font-bold text-slate-800">Reportes generales</h1><p className="text-slate-500 mt-1">Indicadores calculados desde PostgreSQL.</p></div>
      {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg">{error}</div>}
      {!reporte && !error && <p className="text-slate-500">Generando reporte...</p>}
      {reporte && <>
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-5 gap-4">
          {[
            ['Ventas', `$${reporte.totalVentas.toFixed(2)}`],
            ['Pedidos', reporte.totalPedidos],
            ['Ticket promedio', `$${reporte.ticketPromedio.toFixed(2)}`],
            ['Comisiones', `$${reporte.totalComisiones.toFixed(2)}`],
            ['Pendiente de pago', `$${reporte.pendientePagar.toFixed(2)}`],
          ].map(([label, value]) => <div key={label} className="bg-white border rounded-lg p-5"><p className="text-xs uppercase font-semibold text-slate-500">{label}</p><p className="text-2xl font-bold text-indigo-600 mt-2">{value}</p></div>)}
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <section className="bg-white border rounded-lg overflow-x-auto"><h2 className="font-bold p-4 border-b">Productos principales</h2><table className="w-full text-sm"><thead className="bg-slate-50"><tr><th className="p-3 text-left">Producto</th><th className="p-3 text-right">Unidades</th><th className="p-3 text-right">Ventas</th></tr></thead><tbody>{reporte.productosPrincipales.map((item) => <tr key={item.nombre} className="border-t"><td className="p-3">{item.nombre}</td><td className="p-3 text-right">{item.cantidad}</td><td className="p-3 text-right">${item.ventas.toFixed(2)}</td></tr>)}</tbody></table></section>
          <section className="bg-white border rounded-lg overflow-x-auto"><h2 className="font-bold p-4 border-b">Referidos principales</h2><table className="w-full text-sm"><thead className="bg-slate-50"><tr><th className="p-3 text-left">Referido</th><th className="p-3 text-right">Ventas</th><th className="p-3 text-right">Comisiones</th></tr></thead><tbody>{reporte.referidosPrincipales.map((item) => <tr key={item.nombre} className="border-t"><td className="p-3">{item.nombre}</td><td className="p-3 text-right">${item.ventas.toFixed(2)}</td><td className="p-3 text-right">${item.comisiones.toFixed(2)}</td></tr>)}</tbody></table></section>
        </div>
      </>}
    </div>
  );
};

export default AdminReportes;
