import { useEffect, useMemo, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { getRed, type Referido } from '../services/api';

const aLista = (nodos: Referido[]): Referido[] =>
  nodos.flatMap((nodo) => [nodo, ...aLista(nodo.hijos ?? [])]);

const Dashboard = () => {
  const { userEmail } = useAuth();
  const [red, setRed] = useState<Referido[]>([]);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!userEmail) return;
    getRed().then(setRed).catch((err: Error) => setError(err.message));
  }, [userEmail]);

  const resumen = useMemo(() => {
    const personas = aLista(red);
    return {
      ventas: personas.reduce((total, persona) => total + persona.ventas, 0),
      comisiones: personas.reduce((total, persona) => total + persona.comisiones, 0),
      activos: personas.length,
    };
  }, [red]);

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Resumen General</h1>
      {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg mb-4">{error}</div>}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow-sm border border-slate-100">
          <p className="text-sm text-slate-500 uppercase font-semibold">Ventas de la red</p>
          <p className="text-3xl font-bold text-indigo-600">${resumen.ventas.toFixed(2)}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-sm border border-slate-100">
          <p className="text-sm text-slate-500 uppercase font-semibold">Comisiones</p>
          <p className="text-3xl font-bold text-indigo-600">${resumen.comisiones.toFixed(2)}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-sm border border-slate-100">
          <p className="text-sm text-slate-500 uppercase font-semibold">Referidos activos</p>
          <p className="text-3xl font-bold text-indigo-600">{resumen.activos}</p>
        </div>
      </div>
      <p className="text-sm text-slate-500 mt-5">Después de una compra, las comisiones se calculan en segundo plano y este resumen muestra el resultado actualizado.</p>
    </div>
  );
};

export default Dashboard;
