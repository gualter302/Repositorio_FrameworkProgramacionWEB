import { useEffect, useState } from 'react';
import { getComisionesAdmin, liquidarComisiones, type Comision } from '../services/api';

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
  return <div className="space-y-5"><div><h1 className="text-2xl font-bold text-slate-800">Libro de comisiones</h1><p className="text-slate-500 mt-1">Selecciona comisiones pendientes del mismo referido para liquidarlas.</p></div><div className="flex gap-3 items-center"><select value={filtro} onChange={(e) => { setFiltro(e.target.value); setSeleccionadas([]); }} className="border rounded-lg p-2"><option value="pendiente">Pendientes</option><option value="pagada">Pagadas</option><option value="">Todas</option></select><button disabled={!seleccionadas.length || filtro === 'pagada'} onClick={() => void liquidar()} className="bg-indigo-600 text-white px-4 py-2 rounded-lg disabled:opacity-50">Crear liquidación</button></div>{error && <div className="bg-red-50 text-red-700 p-3 rounded-lg">{error}</div>}<div className="bg-white border rounded-lg overflow-x-auto"><table className="w-full text-left text-sm"><thead className="bg-slate-50"><tr><th className="p-3">Sel.</th><th className="p-3">Pedido</th><th className="p-3">Referido</th><th className="p-3">Monto</th><th className="p-3">Estado</th></tr></thead><tbody>{comisiones.map((item) => <tr key={item.id} className="border-t"><td className="p-3"><input type="checkbox" disabled={item.estado !== 'pendiente'} checked={seleccionadas.includes(item.id)} onChange={() => alternar(item.id)} /></td><td className="p-3">{item.numeroPedido}</td><td className="p-3">{item.referido}</td><td className="p-3">${item.monto.toFixed(2)}</td><td className="p-3">{item.estado}</td></tr>)}</tbody></table></div></div>;
};
export default AdminComisiones;
