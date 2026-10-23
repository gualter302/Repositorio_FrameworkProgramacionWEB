import { useEffect, useRef, type ReactNode } from 'react';
import { createPortal } from 'react-dom';
import { Boxes, Database, Orbit, ShieldCheck, TriangleAlert, Wifi, X, type LucideIcon } from 'lucide-react';

type IconType = LucideIcon;

export const BrandMark = ({ compact = false }: { compact?: boolean }) => (
  <div className="brand-mark" aria-label="MultiCatálogo">
    <span className="brand-symbol"><Boxes size={21} aria-hidden="true" /></span>
    {!compact && (
      <span className="brand-copy">
        <span>MultiCatálogo</span>
        <small>Oscar Gonzabay</small>
      </span>
    )}
  </div>
);

interface PageHeaderProps {
  kicker: string;
  title: string;
  description?: string;
  action?: ReactNode;
}

export const PageHeader = ({ kicker, title, description, action }: PageHeaderProps) => (
  <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
    <div className="page-header">
      <p className="page-kicker">{kicker}</p>
      <h1 className="page-title">{title}</h1>
      {description && <p className="page-copy">{description}</p>}
    </div>
    {action && <div className="shrink-0">{action}</div>}
  </header>
);

interface EmptyStateProps {
  icon: IconType;
  title: string;
  description: string;
  action?: ReactNode;
}

export const EmptyState = ({ icon: Icon, title, description, action }: EmptyStateProps) => (
  <div className="empty-state">
    <span className="empty-icon"><Icon size={26} aria-hidden="true" /></span>
    <h2 className="section-title">{title}</h2>
    <p className="section-copy max-w-md">{description}</p>
    {action}
  </div>
);

const toneForStatus = (status: string) => {
  const normalized = status.toLowerCase();
  if (['entregado', 'pagada', 'activo'].includes(normalized)) return 'success';
  if (['cancelado', 'inactivo'].includes(normalized)) return 'danger';
  return 'info';
};

export const StatusBadge = ({ children }: { children: string }) => (
  <span className="status-pill" data-tone={toneForStatus(children)}>{children}</span>
);

interface ConfirmDialogProps {
  open: boolean;
  title: string;
  description: string;
  confirmLabel: string;
  onConfirm: () => void;
  onCancel: () => void;
  busy?: boolean;
  tone?: 'danger' | 'primary';
}

export const ConfirmDialog = ({
  open,
  title,
  description,
  confirmLabel,
  onConfirm,
  onCancel,
  busy = false,
  tone = 'danger',
}: ConfirmDialogProps) => {
  const cancelButton = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (!open) return;
    const overflowAnterior = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    cancelButton.current?.focus();
    const cerrarConEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && !busy) onCancel();
    };
    document.addEventListener('keydown', cerrarConEscape);
    return () => {
      document.body.style.overflow = overflowAnterior;
      document.removeEventListener('keydown', cerrarConEscape);
    };
  }, [busy, onCancel, open]);

  if (!open) return null;

  return createPortal(
    <div className="modal-backdrop" onMouseDown={() => { if (!busy) onCancel(); }}>
      <section
        role="alertdialog"
        aria-modal="true"
        aria-labelledby="confirm-dialog-title"
        aria-describedby="confirm-dialog-description"
        className="confirm-dialog"
        data-tone={tone}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <button type="button" className="icon-button modal-close" onClick={onCancel} disabled={busy} aria-label="Cerrar confirmación"><X size={18} /></button>
        <span className="confirm-dialog__icon"><TriangleAlert size={26} aria-hidden="true" /></span>
        <p className="page-kicker">Confirmación</p>
        <h2 id="confirm-dialog-title" className="confirm-dialog__title">{title}</h2>
        <p id="confirm-dialog-description" className="confirm-dialog__description">{description}</p>
        <div className="confirm-dialog__actions">
          <button ref={cancelButton} type="button" className="btn-secondary" onClick={onCancel} disabled={busy}>Volver</button>
          <button type="button" className={tone === 'danger' ? 'btn-danger' : 'btn-primary'} onClick={onConfirm} disabled={busy}>{busy ? 'Procesando...' : confirmLabel}</button>
        </div>
      </section>
    </div>,
    document.body,
  );
};

interface AuthShellProps {
  kicker: string;
  title: string;
  description: string;
  children: ReactNode;
  footer: ReactNode;
}

export const AuthShell = ({ kicker, title, description, children, footer }: AuthShellProps) => (
  <main className="auth-page">
    <section className="auth-visual" aria-label="Beneficios de MultiCatálogo">
      <BrandMark />
      <div className="holo-stage" aria-hidden="true">
        <span className="holo-ring" />
        <span className="holo-ring" />
        <span className="holo-ring" />
        <span className="holo-core"><Orbit size={54} strokeWidth={1.25} /></span>
      </div>
      <div className="auth-feature-row">
        <div className="auth-feature"><Wifi size={17} aria-hidden="true" /><strong>PWA</strong>Catálogo disponible sin conexión</div>
        <div className="auth-feature"><Database size={17} aria-hidden="true" /><strong>24H</strong>Copia pública controlada</div>
        <div className="auth-feature"><ShieldCheck size={17} aria-hidden="true" /><strong>ID</strong>Compras sin duplicados</div>
      </div>
    </section>
    <section className="auth-form-column">
      <div className="auth-card">
        <div className="lg:hidden"><BrandMark /></div>
        <p className="page-kicker mt-7 lg:mt-0">{kicker}</p>
        <h1 className="auth-title">{title}</h1>
        <p className="page-copy mb-7">{description}</p>
        {children}
        <div className="mt-6 text-center text-sm text-slate-400">{footer}</div>
      </div>
    </section>
  </main>
);
