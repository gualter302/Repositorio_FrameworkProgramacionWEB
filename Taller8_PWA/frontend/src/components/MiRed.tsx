import { useEffect, useMemo, useRef, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { getRed, type Referido } from '../services/api';
import { Check, Copy, Network } from 'lucide-react';
import { EmptyState, PageHeader } from './Ui';

const aLista = (nodos: Referido[]): Referido[] =>
  nodos.flatMap((nodo) => [nodo, ...aLista(nodo.hijos ?? [])]);

const MiRed = () => {
  const { userEmail } = useAuth();
  const [red, setRed] = useState<Referido[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState('');
  const [copiado, setCopiado] = useState(false);
  const enlaceRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!userEmail) return;
    getRed()
      .then(setRed)
      .catch((err: Error) => setError(err.message))
      .finally(() => setCargando(false));
  }, [userEmail]);

  const referidos = useMemo(() => red.flatMap((raiz) => aLista(raiz.hijos ?? [])), [red]);
  const codigoPropio = red[0]?.codigo;
  const enlaceReferido = codigoPropio ? `${window.location.origin}/registro?ref=${codigoPropio}` : '';

  const copiarEnlace = async () => {
    try {
      if (navigator.clipboard && window.isSecureContext) await navigator.clipboard.writeText(enlaceReferido);
      else {
        enlaceRef.current?.select();
        if (!document.execCommand('copy')) throw new Error('No se pudo copiar');
      }
      setCopiado(true);
      window.setTimeout(() => setCopiado(false), 2000);
    } catch {
      enlaceRef.current?.select();
    }
  };

  return (
    <div className="page page-stack">
      <PageHeader kicker="Comunidad" title="Mi red de referidos" description="Comparte tu enlace personal y consulta las ventas y comisiones generadas por tu red." />
      {enlaceReferido && (
        <div className="glass-panel glass-panel--accent">
          <p className="page-kicker">Invitación</p><h2 className="section-title">Tu enlace personal</h2>
          <div className="flex flex-col sm:flex-row gap-3 mt-2">
            <input ref={enlaceRef} readOnly value={enlaceReferido} className="form-input flex-1" aria-label="Enlace personal de referido" />
            <button
              type="button"
              onClick={() => void copiarEnlace()}
              className="btn-primary"
            >
              {copiado ? <Check size={17} /> : <Copy size={17} />}{copiado ? 'Copiado' : 'Copiar enlace'}
            </button>
          </div>
        </div>
      )}
      {cargando && <div className="glass-panel animate-pulse text-slate-500">Cargando red...</div>}
      {error && <div className="alert-error" role="alert">{error}</div>}
      {!cargando && !error && referidos.length === 0 && <EmptyState icon={Network} title="Tu red todavía está vacía" description="Comparte tu enlace para empezar a construirla." />}
      {!cargando && !error && referidos.length > 0 && (
        <div className="data-shell">
          <table className="data-table">
            <thead>
              <tr>
                <th>Nombre</th><th>Nivel</th><th>Ventas</th><th>Comisiones</th>
              </tr>
            </thead>
            <tbody>
              {referidos.map((referido) => (
                <tr key={referido.id}>
                  <td className="font-medium text-white">{referido.nombre}</td>
                  <td className="text-slate-400">Nivel {referido.nivel}</td>
                  <td className="mono text-cyan-300">${referido.ventas.toFixed(2)}</td>
                  <td className="mono text-emerald-300">${referido.comisiones.toFixed(2)}</td>
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
