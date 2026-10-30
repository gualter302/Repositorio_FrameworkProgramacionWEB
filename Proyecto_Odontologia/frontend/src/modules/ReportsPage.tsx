import { useEffect, useState } from 'react'
import { Activity, CalendarX2, CircleDollarSign, ClipboardCheck, FileClock, Mail, PackageSearch } from 'lucide-react'
import { api, type AuditEntry, type OutboxMessage, type Report } from '../api'

export default function ReportsPage({ role }: { role: string }) {
  const [report, setReport] = useState<Report | null>(null)
  const [audit, setAudit] = useState<AuditEntry[]>([])
  const [outbox, setOutbox] = useState<OutboxMessage[]>([])
  const [error, setError] = useState('')
  const canReviewTools = role === 'administrador' || role === 'auditor'
  useEffect(() => {
    api.report().then(setReport).catch(reason => setError(reason.message))
    if (canReviewTools) {
      api.audit().then(setAudit).catch(reason => setError(reason.message))
      api.outbox().then(setOutbox).catch(reason => setError(reason.message))
    }
  }, [canReviewTools])
  const cards = [['Ingresos', `$${(report?.revenue ?? 0).toFixed(2)}`, CircleDollarSign], ['Citas atendidas', report?.attended ?? 0, ClipboardCheck], ['Canceladas / No-show', `${report?.cancelled ?? 0} / ${report?.noShow ?? 0}`, CalendarX2], ['Valor inventario', `$${(report?.inventoryValue ?? 0).toFixed(2)}`, PackageSearch], ['Mensajes pendientes', report?.pendingMessages ?? 0, FileClock], ['Eventos auditados', report?.auditEvents ?? 0, Activity]] as const
  return <>
    <header className="page-head"><div><span className="eyebrow">Indicadores y trazabilidad</span><h1>Reportes</h1><p>Resumen financiero, operativo y bitácora de seguridad.</p></div></header>
    {error && <p className="form-error" role="alert">{error}</p>}
    <section className="metric-grid report-grid">{cards.map(([label, value, Icon]) => <article className="metric" key={label}><span className="metric-icon"><Icon/></span><small>{label}</small><strong>{value}</strong></article>)}</section>
    {canReviewTools && <>
      <section className="panel audit-panel"><div className="panel-head"><h2>Auditoría inmutable</h2><span>Últimos {Math.min(audit.length, 20)} de {audit.length} eventos</span></div><div className="audit-list">{[...audit].reverse().slice(0, 20).map(item => <div key={item.id}><span className="audit-dot"/><div><b>{item.action} · {item.entity}</b><small>{item.userEmail} · {new Date(item.createdAt).toLocaleString('es-EC')} · {item.ip}</small></div><code>#{item.entityId}</code></div>)}</div>{audit.length === 0 && <p className="empty">Las acciones nuevas quedarán registradas aquí.</p>}</section>
      <section className="panel audit-panel"><div className="panel-head"><h2>Bandeja de notificaciones</h2><Mail/></div><p className="info-message">Los mensajes se generan automáticamente. El envío real requiere conectar WhatsApp, SMS o correo.</p><div className="audit-list">{outbox.length === 0 ? <p className="empty">No hay mensajes en cola.</p> : [...outbox].reverse().slice(0, 15).map(item => <div key={item.id}><span className="audit-dot"/><div><b>{item.template}</b><small>{item.recipient} · {item.channel} · {new Date(item.createdAt).toLocaleString('es-EC')}</small></div><code>{item.status}</code></div>)}</div></section>
    </>}
  </>
}
