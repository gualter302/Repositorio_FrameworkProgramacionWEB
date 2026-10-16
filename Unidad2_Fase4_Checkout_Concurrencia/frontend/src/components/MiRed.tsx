import { useEffect, useMemo, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { getRed, type Referido } from '../services/api';

const aLista = (nodos: Referido[]): Referido[] =>
  nodos.flatMap((nodo) => [nodo, ...aLista(nodo.hijos ?? [])]);

const MiRed = () => {
  const { userEmail } = useAuth();
  const [red, setRed] = useState<Referido[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!userEmail) return;
    getRed(userEmail)
      .then(setRed)
      .catch((err: Error) => setError(err.message))
      .finally(() => setCargando(false));
  }, [userEmail]);

  const referidos = useMemo(() => aLista(red), [red]);

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Mi Red de Referidos</h1>
      {cargando && <p className="text-slate-500">Cargando red...</p>}
      {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg mb-4">{error}</div>}
      {!cargando && !error && (
        <div className="bg-white rounded-lg border border-slate-200 shadow-sm overflow-x-auto">
          <table className="w-full text-left">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="p-4 font-semibold text-slate-600">Nombre</th>
                <th className="p-4 font-semibold text-slate-600">Nivel</th>
                <th className="p-4 font-semibold text-slate-600">Ventas</th>
                <th className="p-4 font-semibold text-slate-600">Comisiones</th>
              </tr>
            </thead>
            <tbody>
              {referidos.map((referido) => (
                <tr key={referido.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="p-4 text-slate-700">{referido.nombre}</td>
                  <td className="p-4 text-slate-500">Nivel {referido.nivel}</td>
                  <td className="p-4 text-indigo-600 font-medium">${referido.ventas.toFixed(2)}</td>
                  <td className="p-4 text-green-600 font-medium">${referido.comisiones.toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default MiRed;
