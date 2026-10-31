import { useEffect, useState } from 'react'
import { AlertTriangle, Boxes, PackagePlus, RefreshCw } from 'lucide-react'
import { api, type InventoryItem } from '../api'

export default function InventoryPage({ role }: { role: string }) {
  const [items, setItems] = useState<InventoryItem[]>([])
  const [form, setForm] = useState({ name: '', unit: 'unidad', stock: 0, minimum: 0, unitCost: 0 })
  const [adjust, setAdjust] = useState<Record<number, number>>({})
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const canEdit = role === 'administrador'
  const load = () => api.inventory().then(setItems).catch(reason => setError(reason.message))
  useEffect(() => { void load() }, [])
  const create = async (event: React.FormEvent) => { event.preventDefault(); setError(''); setMessage(''); try { await api.createInventoryItem(form); setForm({ name: '', unit: 'unidad', stock: 0, minimum: 0, unitCost: 0 }); setMessage('Insumo creado'); await load() } catch (reason) { setError(reason instanceof Error ? reason.message : 'No se pudo guardar') } }
  const update = async (id: number) => { const quantity = adjust[id] ?? 0; if (!quantity) return; setError(''); setMessage(''); try { await api.adjustInventory(id, quantity, 'Ajuste manual de inventario'); setAdjust({ ...adjust, [id]: 0 }); setMessage('Existencia actualizada'); await load() } catch (reason) { setError(reason instanceof Error ? reason.message : 'No se pudo ajustar') } }

  return <>
    <header className="page-head"><div><span className="eyebrow">Insumos por consultorio</span><h1>Inventario clínico</h1><p>Control de existencias y alertas de reposición.</p></div><span className="date-chip">{items.filter(item => item.stock <= item.minimum).length} alertas</span></header>
    {!canEdit && <p className="info-message">Perfil de consulta: los ajustes y nuevos insumos corresponden al administrador.</p>}
    {message && <p className="success-message" role="status">{message}</p>}{error && <p className="form-error" role="alert">{error}</p>}
    <section className={canEdit ? 'inventory-layout' : ''}><article className="panel"><div className="panel-head"><h2>Existencias</h2><Boxes/></div><div className="stock-list">{items.map(item => { const low = item.stock <= item.minimum; const quantity = adjust[item.id] ?? 0; return <div className={`stock-card ${low ? 'low' : ''}`} key={item.id}><span className="stock-icon">{low ? <AlertTriangle/> : <Boxes/>}</span><div><b>{item.name}</b><small>Mínimo {item.minimum} {item.unit}</small></div><strong>{item.stock.toFixed(2)} <small>{item.unit}</small></strong>{canEdit && <div className="stock-adjust"><input aria-label={`Cantidad de ajuste para ${item.name}`} type="number" min={-item.stock} max="1000000" step="0.01" value={quantity} onChange={event => setAdjust({ ...adjust, [item.id]: Number(event.target.value) })}/><button type="button" aria-label={`Aplicar ajuste de ${item.name}`} title="Aplicar ajuste" className="secondary" disabled={quantity === 0} onClick={() => void update(item.id)}><RefreshCw/></button></div>}</div> })}</div></article>
      {canEdit && <form className="panel stack-form" onSubmit={create}><div className="panel-head"><h2>Nuevo insumo</h2><PackagePlus/></div><label>Nombre<input required minLength={2} maxLength={100} value={form.name} onChange={event => setForm({ ...form, name: event.target.value })}/></label><label>Unidad<input required minLength={1} maxLength={30} value={form.unit} onChange={event => setForm({ ...form, unit: event.target.value.replace(/[^\p{L}\s]/gu, '') })}/></label><div className="form-row"><label>Stock inicial<input required type="number" min="0" max="1000000" step="0.01" value={form.stock} onChange={event => setForm({ ...form, stock: Number(event.target.value) })}/></label><label>Mínimo<input required type="number" min="0" max="1000000" step="0.01" value={form.minimum} onChange={event => setForm({ ...form, minimum: Number(event.target.value) })}/></label></div><label>Costo unitario<input required type="number" min="0" max="1000000" step="0.01" value={form.unitCost} onChange={event => setForm({ ...form, unitCost: Number(event.target.value) })}/></label><button className="primary">Agregar al inventario</button></form>}
    </section>
  </>
}
