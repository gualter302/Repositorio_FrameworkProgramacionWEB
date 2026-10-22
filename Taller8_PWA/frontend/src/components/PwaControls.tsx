import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { useRegisterSW } from 'virtual:pwa-register/react';
import { Download, Info, RefreshCw, WifiOff } from 'lucide-react';

interface InstallEvent extends Event { prompt(): Promise<void>; userChoice: Promise<{ outcome: string }> }
export default function PwaControls() {
  const [online, setOnline] = useState(navigator.onLine);
  const [install, setInstall] = useState<InstallEvent | null>(null);
  const { pathname } = useLocation();
  const { needRefresh: [refresh], updateServiceWorker } = useRegisterSW();
  useEffect(() => {
    const connectivity = () => setOnline(navigator.onLine);
    const invitation = (event: Event) => { event.preventDefault(); setInstall(event as InstallEvent); };
    const installed = () => setInstall(null);
    window.addEventListener('online', connectivity);
    window.addEventListener('offline', connectivity);
    window.addEventListener('beforeinstallprompt', invitation);
    window.addEventListener('appinstalled', installed);
    return () => {
      window.removeEventListener('online', connectivity); window.removeEventListener('offline', connectivity);
      window.removeEventListener('beforeinstallprompt', invitation); window.removeEventListener('appinstalled', installed);
    };
  }, []);
  const standalone = window.matchMedia('(display-mode: standalone)').matches;
  const canUpdate = pathname === '/catalogo' || pathname === '/carrito';
  return <aside className="pointer-events-none fixed bottom-4 right-4 z-[500] flex max-w-[calc(100vw-2rem)] flex-col items-end gap-2 text-sm" aria-live="polite">
    {!online && <div className="glass-panel flex max-w-sm items-start gap-3 border-amber-300/25 bg-slate-950/95 p-3 text-amber-100"><WifiOff className="mt-0.5 shrink-0" size={18} /><p>Sin conexión. Puedes consultar el catálogo guardado y preparar tu carrito. Para comprar debes conectarte.</p></div>}
    {refresh && <div className={`glass-panel max-w-sm bg-slate-950/95 p-3 ${canUpdate ? 'pointer-events-auto' : ''}`}><div className="flex gap-3"><RefreshCw className="mt-0.5 shrink-0 text-cyan-300" size={18} /><div><p className="font-semibold text-white">Hay una nueva versión.</p><p className="mt-1 text-slate-400">Termina tus operaciones antes de actualizar.</p><button disabled={!canUpdate} className="btn-secondary mt-3 w-full" onClick={() => updateServiceWorker(true)}>Actualizar desde catálogo o carrito</button></div></div></div>}
    {install && <button className="btn-primary pointer-events-auto shadow-2xl" onClick={async () => { await install.prompt(); await install.userChoice; setInstall(null); }}><Download size={17} />Instalar aplicación</button>}
    {!install && !standalone && <details className="group pointer-events-auto relative">
      <summary className="icon-button border-cyan-200/15 bg-slate-950/90 shadow-xl" aria-label="Cómo instalar la aplicación" title="Cómo instalar la aplicación"><Info size={19} /></summary>
      <div className="glass-panel absolute bottom-full right-0 mb-2 w-72 bg-slate-950/95 p-4 text-slate-300"><p className="font-semibold text-white">Instalar MultiCatálogo</p><p className="mt-2 leading-relaxed">En Android o computadora abre el menú del navegador y elige Instalar aplicación. En iPhone usa Safari, Compartir y Añadir a pantalla de inicio.</p></div>
    </details>}
  </aside>;
}
