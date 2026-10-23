import { useEffect, useMemo, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { getRed, type Referido } from '../services/api';
import { ArrowUpRight, BadgeDollarSign, DollarSign, Network } from 'lucide-react';
import { Link } from 'react-router-dom';
import { PageHeader } from './Ui';

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
    const referidos = red.flatMap((raiz) => aLista(raiz.hijos ?? []));
    return {
      ventas: red.reduce((total, raiz) => total + raiz.ventas, 0),
      comisiones: red.reduce((total, raiz) => total + raiz.comisiones, 0),
      activos: referidos.length,
    };
  }, [red]);

  return (
    <div className="page page-stack">
      <PageHeader kicker="Actividad" title="Resumen general" description="Una lectura rápida de tus ventas, comisiones y personas referidas." action={<Link to="/catalogo" className="btn-secondary">Ver catálogo <ArrowUpRight size={17} /></Link>} />
      {error && <div className="alert-error" role="alert">{error}</div>}
      <div className="metric-grid">
        <div className="metric-card glass-panel--accent">
          <DollarSign className="mb-5 text-cyan-300" size={22} aria-hidden="true" />
          <p className="metric-label">Ventas de la red</p>
          <p className="metric-value">${resumen.ventas.toFixed(2)}</p>
        </div>
        <div className="metric-card">
          <BadgeDollarSign className="mb-5 text-cyan-300" size={22} aria-hidden="true" />
          <p className="metric-label">Comisiones</p>
          <p className="metric-value">${resumen.comisiones.toFixed(2)}</p>
        </div>
        <div className="metric-card">
          <Network className="mb-5 text-cyan-300" size={22} aria-hidden="true" />
          <p className="metric-label">Referidos activos</p>
          <p className="metric-value">{resumen.activos}</p>
        </div>
      </div>
      <div className="glass-panel flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div><h2 className="section-title">Datos que se actualizan contigo</h2><p className="section-copy mt-1 text-sm">Después de una compra, las comisiones se calculan en segundo plano y aparecen aquí.</p></div>
        <Link to="/mi-red" className="btn-secondary shrink-0">Explorar mi red <ArrowUpRight size={17} /></Link>
      </div>
    </div>
  );
};

export default Dashboard;
