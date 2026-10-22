import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { useRegisterSW } from 'virtual:pwa-register/react';

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
  return <aside className="bg-indigo-950 text-white p-3 text-sm" aria-live="polite">
    {!online && <p>Sin conexión. Puedes consultar el catálogo guardado y preparar tu carrito. Para comprar debes conectarte.</p>}
    {install && <button className="min-h-11 px-4 underline" onClick={async () => { await install.prompt(); await install.userChoice; setInstall(null); }}>Instalar MultiCatálogo</button>}
    {!install && !window.matchMedia('(display-mode: standalone)').matches && <details><summary className="cursor-pointer min-h-11 py-3">Cómo instalar la aplicación</summary><p>En Android o computadora: menú del navegador → Instalar aplicación. En iPhone: Safari → Compartir → Añadir a pantalla de inicio.</p></details>}
    {refresh && <div><p>Hay una nueva versión. Termina tus operaciones antes de actualizar.</p>
      <button disabled={pathname !== '/catalogo' && pathname !== '/carrito'} className="min-h-11 underline disabled:opacity-50" onClick={() => updateServiceWorker(true)}>Actualizar desde catálogo o carrito</button></div>}
  </aside>;
}
