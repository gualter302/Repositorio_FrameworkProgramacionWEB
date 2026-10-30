export type Session = { token: string; name: string; email: string; role: string; expiresIn: number }
export type Patient = { id: number; document: string; firstName: string; lastName: string; birthDate: string; phone: string; email: string; address: string; emergencyContact: string; guardian: string; insurance: string; conditions: string[]; allergies: string[]; surgeries: string; anticoagulant: boolean; alerts: string[]; createdAt: string }
export type Appointment = { id: number; patientId: number; patientName: string; dentist: string; chair: string; start: string; end: string; status: string; reason: string }
export type Dashboard = { patients: number; todayAppointments: number; waiting: number; criticalAlerts: number; estimatedRevenue: number; lowStock: number; pendingInvoices: number; nextAppointments?: Appointment[] }
export type OdontogramEntry = { id: number; patientId: number; tooth: number; surface: string; condition: string; nomenclature: string; notes: string; recordedBy: string; recordedAt: string }
export type PeriodontalRecord = { id: number; patientId: number; tooth: number; depths: number[]; margin: number; furcation: number; mobility: number; bleeding: boolean; recordedAt: string }
export type Amendment = { id: number; reason: string; content: string; signedBy: string; createdAt: string }
export type ClinicalNote = { id: number; patientId: number; subjective: string; objective: string; assessment: string; plan: string; signedBy: string; signedAt: string; lockedAt: string; amendments: Amendment[] }
export type Consent = { id: number; patientId: number; procedure: string; risks: string; patientName: string; signatureData: string; otpValidated: boolean; signedAt: string }
export type TreatmentItem = { id?: number; procedure: string; phase: number; sessions: number; professionalFee: number; laboratoryCost: number; discount: number; status?: string }
export type TreatmentPlan = { id: number; patientId: number; patientName: string; status: string; items: TreatmentItem[]; total: number; createdAt: string }
export type ClinicalRecord = { patient: Patient; odontogram: OdontogramEntry[]; periodontal: PeriodontalRecord[]; notes: ClinicalNote[]; consents: Consent[]; plans: TreatmentPlan[] }
export type InventoryItem = { id: number; name: string; unit: string; stock: number; minimum: number; unitCost: number; updatedAt: string }
export type Invoice = { id: number; number: string; patientId: number; patientName: string; treatmentPlanId: number; subtotal: number; tax: number; total: number; paymentMethod: string; status: string; fiscalStatus: string; createdAt: string }
export type CashClosure = { id: number; date: string; declaredCash: number; expectedCash: number; difference: number; closedBy: string; createdAt: string }
export type Report = { revenue: number; invoices: number; patients: number; appointments: number; attended: number; cancelled: number; noShow: number; inventoryValue: number; lowStock: number; auditEvents: number; pendingMessages: number }
export type AuditEntry = { id: number; userEmail: string; action: string; entity: string; entityId: string; ip: string; createdAt: string }
export type Settlement = { id: number; professional: string; period: string; amount: number; status: string; createdAt: string }
export type OutboxMessage = { id: number; channel: string; recipient: string; template: string; status: string; createdAt: string }

const asList = <T,>(value: T[] | null | undefined): T[] => value ?? []
const normalizeClinicalRecord = (record: ClinicalRecord): ClinicalRecord => ({
  ...record,
  odontogram: asList(record.odontogram),
  periodontal: asList(record.periodontal),
  notes: asList(record.notes).map(note => ({ ...note, amendments: asList(note.amendments) })),
  consents: asList(record.consents),
  plans: asList(record.plans).map(plan => ({ ...plan, items: asList(plan.items) })),
})

const base = '/api'
let token = localStorage.getItem('sigo_token') ?? ''
export const setToken = (value: string) => { token = value; if (value) localStorage.setItem('sigo_token', value); else localStorage.removeItem('sigo_token') }

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(base + path, { ...options, headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}), ...options.headers } })
  const body = response.status === 204 ? null : await response.json()
  if (!response.ok) { const error = new Error(body?.message ?? 'No se pudo completar la operación') as Error & { code?: string }; error.code = body?.code; throw error }
  return body as T
}

export const api = {
  login: (email: string, password: string, otp = '') => request<Session>('/login', { method: 'POST', body: JSON.stringify({ email, password, otp }) }),
  dashboard: () => request<Dashboard>('/dashboard'),
  patients: () => request<Patient[] | null>('/patients').then(asList),
  createPatient: (patient: Partial<Patient>) => request<Patient>('/patients', { method: 'POST', body: JSON.stringify(patient) }),
  updatePatient: (id: number, patient: Partial<Patient>) => request<Patient>(`/patients/${id}`, { method: 'PUT', body: JSON.stringify(patient) }),
  appointments: () => request<Appointment[] | null>('/appointments').then(asList),
  createAppointment: (appointment: Partial<Appointment>) => request<Appointment>('/appointments', { method: 'POST', body: JSON.stringify(appointment) }),
  changeAppointmentStatus: (id: number, status: string) => request<Appointment>(`/appointments/${id}/status`, { method: 'PATCH', body: JSON.stringify({ status }) }),
  clinicalRecord: (patientId: number) => request<ClinicalRecord>(`/clinical/${patientId}`).then(normalizeClinicalRecord),
  addOdontogram: (patientId: number, entry: Partial<OdontogramEntry>) => request<OdontogramEntry>(`/clinical/${patientId}/odontogram`, { method: 'POST', body: JSON.stringify(entry) }),
  addPeriodontal: (patientId: number, entry: Partial<PeriodontalRecord>) => request<PeriodontalRecord>(`/clinical/${patientId}/periodontogram`, { method: 'POST', body: JSON.stringify(entry) }),
  addNote: (patientId: number, note: Partial<ClinicalNote>) => request<ClinicalNote>(`/clinical/${patientId}/notes`, { method: 'POST', body: JSON.stringify(note) }),
  addAmendment: (noteId: number, reason: string, content: string) => request<Amendment>(`/clinical/notes/${noteId}/amendments`, { method: 'POST', body: JSON.stringify({ reason, content }) }),
  addConsent: (patientId: number, consent: Record<string, string>) => request<Consent>(`/clinical/${patientId}/consents`, { method: 'POST', body: JSON.stringify(consent) }),
  plans: () => request<TreatmentPlan[] | null>('/treatment-plans').then(asList),
  createPlan: (plan: { patientId: number; items: TreatmentItem[] }) => request<TreatmentPlan>('/treatment-plans', { method: 'POST', body: JSON.stringify(plan) }),
  changePlanStatus: (id: number, status: string) => request<TreatmentPlan>(`/treatment-plans/${id}/status`, { method: 'PATCH', body: JSON.stringify({ status }) }),
  inventory: () => request<InventoryItem[] | null>('/inventory').then(asList),
  createInventoryItem: (item: Partial<InventoryItem>) => request<InventoryItem>('/inventory', { method: 'POST', body: JSON.stringify(item) }),
  adjustInventory: (id: number, quantity: number, reason: string) => request<InventoryItem>(`/inventory/${id}/adjust`, { method: 'POST', body: JSON.stringify({ quantity, reason }) }),
  invoices: () => request<Invoice[] | null>('/invoices').then(asList),
  createInvoice: (treatmentPlanId: number, paymentMethod: string) => request<Invoice>('/invoices', { method: 'POST', body: JSON.stringify({ treatmentPlanId, paymentMethod }) }),
  closeCash: (declaredCash: number) => request<CashClosure>('/cash-closures', { method: 'POST', body: JSON.stringify({ declaredCash }) }),
  cashClosures: () => request<CashClosure[] | null>('/cash-closures').then(asList),
  report: () => request<Report>('/reports/general'),
  audit: () => request<AuditEntry[] | null>('/audit').then(asList),
  settlements: () => request<Settlement[] | null>('/settlements').then(asList),
  createSettlement: (professional: string, period: string, percentage: number) => request<Settlement>('/settlements', { method: 'POST', body: JSON.stringify({ professional, period, percentage }) }),
  outbox: () => request<OutboxMessage[] | null>('/outbox').then(asList),
}
