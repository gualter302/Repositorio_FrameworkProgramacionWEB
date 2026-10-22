import { useEffect, useMemo, useState } from 'react';
import { getMisComisiones, type Comision } from '../services/api';

const Comisiones = () => {
  const [comisiones, setComisiones] = useState<Comision[]>([]);
  const [error, setError] = useState('');
  useEffect(() => {
    let activo = true;
    getMisComisiones().then((datos) => { if (activo) setComisiones(datos); }).catch((err: Error) => { if (activo) setError(err.message); });
    return () => { activo = false; };
  }, []);
  const total = useMemo(() => comisiones.reduce((suma, item) => suma + item.monto, 0), [comisiones]);
  return <div className="space-y-5"><div><h1 className="text-2xl font-bold text-slate-800">Mis comisiones</h1><p className="text-slate-500 mt-1">Consulta el origen y estado de cada comisión.</p></div>{error && <div className="bg-red-50 text-red-700 p-3 rounded-lg">{error}</div>}<div className="bg-white border rounded-lg p-5"><p className="text-sm text-slate-500">Total histórico</p><p className="text-3xl font-bold text-green-600">${total.toFixed(2)}</p></div><div className="bg-white border rounded-lg overflow-x-auto"><table className="w-full text-left text-sm"><thead className="bg-slate-50"><tr><th className="p-3">Pedido</th><th className="p-3">Nivel</th><th className="p-3">Porcentaje</th><th className="p-3">Monto</th><th className="p-3">Estado</th></tr></thead><tbody>{comisiones.map((item) => <tr key={item.id} className="border-t"><td className="p-3">{item.numeroPedido}</td><td className="p-3">{item.nivel}</td><td className="p-3">{(item.porcentaje * 100).toFixed(0)}%</td><td className="p-3 font-semibold">${item.monto.toFixed(2)}</td><td className="p-3">{item.estado}</td></tr>)}</tbody></table></div></div>;
};
export default Comisiones;
