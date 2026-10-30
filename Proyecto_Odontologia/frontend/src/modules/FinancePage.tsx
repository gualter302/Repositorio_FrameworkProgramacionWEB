import { useEffect, useState } from 'react'
import { Banknote, CheckCircle2, FileText, LockKeyhole, Receipt, XCircle } from 'lucide-react'
import { api, type CashClosure, type Invoice, type Settlement, type TreatmentPlan } from '../api'

const taxRate = 0.15
const withTax = (subtotal: number) => subtotal * (1 + taxRate)

export default function FinancePage({ role }: { role: string }) {
  const [plans, setPlans] = useState<TreatmentPlan[]>([])
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [closures, setClosures] = useState<CashClosure[]>([])
  const [settlements, setSettlements] = useState<Settlement[]>([])
  const [settlementForm, setSettlementForm] = useState({ professional: '', period: new Date().toISOString().slice(0, 7), percentage: 40 })
  const [methods, setMethods] = useState<Record<number, string>>({})
  const [declared, setDeclared] = useState(0)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const canManage = role === 'administrador' || role === 'recepcionista'
  const canSettle = role === 'administrador' || role === 'auditor'
  const load = () => Promise.all([api.plans(), api.invoices(), api.cashClosures(), canSettle ? api.settlements() : Promise.resolve<Settlement[]>([])]).then(([planList, invoiceList, closureList, settlementList]) => { setPlans(planList); setInvoices(invoiceList); setClosures(closureList); setSettlements(settlementList) }).catch(reason => setError(reason.message))
  useEffect(() => { void load() }, [])
  const act = async (action: () => Promise<unknown>, successMessage: string) => { setError(''); setMessage(''); try { await action(); setMessage(successMessage); await load() } catch (reason) { setError(reason instanceof Error ? reason.message : 'No se pudo completar') } }
  const today = new Date().toISOString().slice(0, 10)
  const todayClosure = closures.find(item => item.date === today)

  return <>
    <header className="page-head"><div><span className="eyebrow">Finanzas y autorizaciones</span><h1>Caja y facturación</h1><p>Presupuestos aprobados, cobros y cierre ciego de caja.</p></div><span className="date-chip">{invoices.length} comprobantes</span></header>
    {!canManage && <p className="info-message">Perfil de consulta: puede revisar la información financiera, sin modificarla.</p>}
    {message && <p className="success-message" role="status">{message}</p>}{error && <p className="form-error" role="alert">{error}</p>}
    <section className="split-grid finance-grid">
      <article className="panel"><div className="panel-head"><h2>Planes de tratamiento</h2><FileText/></div><div className="cards-list">{plans.length === 0 ? <p className="empty">Aún no existen presupuestos.</p> : plans.map(plan => {
        const invoiced = invoices.some(invoice => invoice.treatmentPlanId === plan.id)
        const method = methods[plan.id] ?? 'Efectivo'
        return <div className="data-card" key={plan.id}><div><span className="eyebrow">{invoiced ? 'Facturado' : plan.status}</span><b>{plan.patientName}</b><small>{plan.items.length} procedimiento(s) · Total con IVA ${withTax(plan.total).toFixed(2)}</small></div><strong>${plan.total.toFixed(2)}</strong>{canManage && !invoiced && <div className="row-actions">{plan.status === 'Borrador' && <><button type="button" className="secondary" onClick={() => void act(() => api.changePlanStatus(plan.id, 'Aprobado'), 'Presupuesto aprobado')}>Aprobar</button><button type="button" className="secondary" onClick={() => void act(() => api.changePlanStatus(plan.id, 'Rechazado'), 'Presupuesto rechazado')}><XCircle/>Rechazar</button></>}{plan.status.startsWith('Aprobado') && <><label>Forma de pago<select value={method} onChange={event => setMethods({ ...methods, [plan.id]: event.target.value })}><option>Efectivo</option><option>Tarjeta</option><option>Transferencia</option></select></label><button type="button" className="primary" onClick={() => void act(() => api.createInvoice(plan.id, method), 'Factura emitida')}>Facturar</button></>}</div>}</div>
      })}</div></article>
      <article className="panel cash-panel"><div className="panel-head"><h2>Cierre ciego</h2><LockKeyhole/></div>{todayClosure ? <><p>La caja de hoy ya fue cerrada.</p><div className="data-card"><span>Declarado: ${todayClosure.declaredCash.toFixed(2)}</span><span>Esperado: ${todayClosure.expectedCash.toFixed(2)}</span><strong>Diferencia: ${todayClosure.difference.toFixed(2)}</strong></div></> : canManage ? <><p>Ingrese el conteo físico antes de conocer el total esperado. El resultado aparece solo después de cerrar.</p><label>Efectivo contado<input type="number" min="0" max="1000000" step="0.01" value={declared} onChange={event => setDeclared(Number(event.target.value))}/></label><button type="button" className="primary" onClick={() => void act(() => api.closeCash(declared), 'Caja cerrada; revise la diferencia registrada')}><Banknote/>Cerrar caja de hoy</button></> : <p>El cierre de caja está disponible para administración y recepción.</p>}</article>
    </section>
    <section className="split-grid finance-grid">
      <article className="panel"><div className="panel-head"><h2>Liquidaciones profesionales</h2><Banknote/></div>{canSettle ? <form className="stack-form" onSubmit={event => { event.preventDefault(); void act(() => api.createSettlement(settlementForm.professional, settlementForm.period, settlementForm.percentage), 'Liquidación calculada sobre las citas atendidas del período') }}><label>Profesional<input required minLength={3} maxLength={80} value={settlementForm.professional} onChange={event => setSettlementForm({ ...settlementForm, professional: event.target.value })}/></label><div className="form-row"><label>Período<input required type="month" value={settlementForm.period} onChange={event => setSettlementForm({ ...settlementForm, period: event.target.value })}/></label><label>Comisión (%)<input required type="number" min="1" max="100" step="0.5" value={settlementForm.percentage} onChange={event => setSettlementForm({ ...settlementForm, percentage: Number(event.target.value) })}/></label></div><button className="secondary">Calcular liquidación</button></form> : <p className="info-message">Las liquidaciones se calculan desde administración o auditoría.</p>}<div className="record-list">{settlements.length === 0 ? <p className="empty">Sin liquidaciones registradas.</p> : settlements.map(item => <span key={item.id}><b>{item.professional} · {item.period}</b>${item.amount.toFixed(2)} · {item.status}</span>)}</div></article>
      <article className="panel"><div className="panel-head"><h2>Resumen de cobros</h2><Receipt/></div><div className="record-list">{['Efectivo', 'Tarjeta', 'Transferencia', 'Mixto'].map(method => <span key={method}><b>{method}</b>${invoices.filter(item => item.paymentMethod === method).reduce((sum, item) => sum + item.total, 0).toFixed(2)}</span>)}<span><b>Total facturado</b>${invoices.reduce((sum, item) => sum + item.total, 0).toFixed(2)}</span></div></article>
    </section>
    <section className="panel invoices"><div className="panel-head"><h2>Comprobantes</h2><Receipt/></div><div className="table-card inner"><table><thead><tr><th>Número</th><th>Paciente</th><th>Pago</th><th>Total</th><th>Estado fiscal</th></tr></thead><tbody>{invoices.map(item => <tr key={item.id}><td><b>{item.number}</b><small>{new Date(item.createdAt).toLocaleDateString('es-EC')}</small></td><td>{item.patientName}</td><td>{item.paymentMethod}</td><td>${item.total.toFixed(2)}</td><td><span className="status waiting">{item.fiscalStatus}</span></td></tr>)}</tbody></table></div>{invoices.length === 0 && <div className="empty-state"><CheckCircle2/><p>Las facturas emitidas aparecerán aquí.</p></div>}</section>
  </>
}
