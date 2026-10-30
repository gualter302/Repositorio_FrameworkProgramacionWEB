package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type invoiceRequest struct {
	TreatmentPlanID int            `json:"treatmentPlanId"`
	PaymentMethod   string         `json:"paymentMethod"`
	Payments        []paymentEntry `json:"payments"`
}
type inventoryAdjustment struct {
	Quantity float64 `json:"quantity"`
	Reason   string  `json:"reason"`
}
type closureRequest struct {
	DeclaredCash float64 `json:"declaredCash"`
}
type settlementRequest struct {
	Professional string  `json:"professional"`
	Period       string  `json:"period"`
	Percentage   float64 `json:"percentage"`
}

func (s *server) registerFinanceRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/inventory", s.auth(http.HandlerFunc(s.listInventory), "administrador", "odontologo", "recepcionista", "auditor"))
	mux.Handle("POST /api/inventory", s.auth(http.HandlerFunc(s.createInventoryItem), "administrador"))
	mux.Handle("POST /api/inventory/{id}/adjust", s.auth(http.HandlerFunc(s.adjustInventory), "administrador"))
	mux.Handle("GET /api/invoices", s.auth(http.HandlerFunc(s.listInvoices), "administrador", "recepcionista", "auditor"))
	mux.Handle("POST /api/invoices", s.auth(http.HandlerFunc(s.createInvoice), "administrador", "recepcionista"))
	mux.Handle("POST /api/cash-closures", s.auth(http.HandlerFunc(s.createCashClosure), "administrador", "recepcionista"))
	mux.Handle("GET /api/cash-closures", s.auth(http.HandlerFunc(s.listCashClosures), "administrador", "recepcionista", "auditor"))
	mux.Handle("POST /api/settlements", s.auth(http.HandlerFunc(s.createSettlement), "administrador", "auditor"))
	mux.Handle("GET /api/settlements", s.auth(http.HandlerFunc(s.listSettlements), "administrador", "auditor"))
	mux.Handle("GET /api/reports/general", s.auth(http.HandlerFunc(s.generalReport), "administrador", "recepcionista", "auditor"))
	mux.Handle("GET /api/audit", s.auth(http.HandlerFunc(s.listAudit), "administrador", "auditor"))
	mux.Handle("GET /api/outbox", s.auth(http.HandlerFunc(s.listOutbox), "administrador", "auditor"))
}

func (s *server) lowStockCountLocked() int {
	count := 0
	for _, item := range s.state.Inventory {
		if item.Stock <= item.Minimum {
			count++
		}
	}
	return count
}
func (s *server) pendingInvoicesLocked() int {
	count := 0
	for _, item := range s.state.Invoices {
		if item.FiscalStatus == "Pendiente de integración" {
			count++
		}
	}
	return count
}
func (s *server) listInventory(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Inventory)
}
func (s *server) createInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item inventoryItem
	if decode(r, &item) != nil {
		fail(w, 400, "INVALID_JSON", "Artículo de inventario inválido")
		return
	}
	item.Name = strings.TrimSpace(item.Name)
	item.Unit = strings.TrimSpace(item.Unit)
	if len([]rune(item.Name)) < 2 || len([]rune(item.Name)) > 100 || len([]rune(item.Unit)) < 1 || len([]rune(item.Unit)) > 30 || !finiteNonNegative(item.Stock) || !finiteNonNegative(item.Minimum) || !finiteNonNegative(item.UnitCost) {
		fail(w, 422, "VALIDATION_ERROR", "Artículo de inventario inválido")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, current := range s.state.Inventory {
		if strings.EqualFold(current.Name, item.Name) {
			fail(w, 409, "ITEM_EXISTS", "Ya existe un insumo con ese nombre")
			return
		}
	}
	item.ID = s.state.Counters.Inventory
	s.state.Counters.Inventory++
	item.UpdatedAt = time.Now()
	s.state.Inventory = append(s.state.Inventory, item)
	s.auditLocked(r, "CREATE", "inventory_item", strconv.Itoa(item.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar el artículo")
		return
	}
	writeJSON(w, 201, item)
}
func (s *server) adjustInventory(w http.ResponseWriter, r *http.Request) {
	id, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Artículo inválido")
		return
	}
	var req inventoryAdjustment
	if decode(r, &req) != nil || req.Quantity == 0 || math.IsNaN(req.Quantity) || math.IsInf(req.Quantity, 0) || len([]rune(strings.TrimSpace(req.Reason))) < 3 || len([]rune(strings.TrimSpace(req.Reason))) > 160 {
		fail(w, 422, "VALIDATION_ERROR", "Indique cantidad y motivo")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, item := range s.state.Inventory {
		if item.ID == id {
			if item.Stock+req.Quantity < 0 {
				fail(w, 409, "INSUFFICIENT_STOCK", "La salida supera el stock disponible")
				return
			}
			s.state.Inventory[i].Stock += req.Quantity
			s.state.Inventory[i].UpdatedAt = time.Now()
			s.auditLocked(r, "ADJUST", "inventory_item", strconv.Itoa(id))
			if err := s.persistLocked(); err != nil {
				fail(w, 500, "PERSISTENCE_ERROR", "No se pudo ajustar el inventario")
				return
			}
			writeJSON(w, 200, s.state.Inventory[i])
			return
		}
	}
	fail(w, 404, "NOT_FOUND", "Artículo no encontrado")
}

func (s *server) consumeRecipeLocked(reason string) {
	reason = strings.ToLower(reason)
	usage := map[int]float64{1: 0.02}
	if strings.Contains(reason, "endodoncia") {
		usage[2] = 0.25
		usage[3] = 2
	}
	if strings.Contains(reason, "profilaxis") {
		usage[3] = 1
	}
	for id, quantity := range usage {
		for i, item := range s.state.Inventory {
			if item.ID == id && item.Stock >= quantity {
				s.state.Inventory[i].Stock -= quantity
				s.state.Inventory[i].UpdatedAt = time.Now()
			}
		}
	}
}

func (s *server) listInvoices(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Invoices)
}
func (s *server) createInvoice(w http.ResponseWriter, r *http.Request) {
	var req invoiceRequest
	if decode(r, &req) != nil || req.TreatmentPlanID < 1 || (req.PaymentMethod == "" && len(req.Payments) == 0) {
		fail(w, 422, "VALIDATION_ERROR", "Plan y forma de pago son obligatorios")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var plan treatmentPlan
	found := false
	for _, item := range s.state.Plans {
		if item.ID == req.TreatmentPlanID {
			plan = item
			found = true
			break
		}
	}
	if !found {
		fail(w, 404, "NOT_FOUND", "Presupuesto no encontrado")
		return
	}
	for _, current := range s.state.Invoices {
		if current.TreatmentPlanID == plan.ID && current.Status != "Anulada" {
			fail(w, 409, "PLAN_ALREADY_INVOICED", "El presupuesto ya tiene una factura emitida")
			return
		}
	}
	if !strings.HasPrefix(plan.Status, "Aprobado") {
		fail(w, 409, "PLAN_NOT_APPROVED", "El presupuesto aún no está aprobado")
		return
	}
	tax := plan.Total * 0.15
	total := plan.Total + tax
	applyPayments(&req, total)
	paid := 0.0
	allowedPayments := map[string]bool{"Efectivo": true, "Tarjeta": true, "Transferencia": true}
	for _, payment := range req.Payments {
		if !allowedPayments[payment.Method] || payment.Amount <= 0 || math.IsNaN(payment.Amount) || math.IsInf(payment.Amount, 0) {
			fail(w, 422, "VALIDATION_ERROR", "Las formas de pago no son válidas")
			return
		}
		paid += payment.Amount
	}
	if math.Abs(paid-total) > 0.01 {
		fail(w, 422, "PAYMENT_MISMATCH", "La suma de pagos debe coincidir con el total")
		return
	}
	paymentMethod := req.Payments[0].Method
	if len(req.Payments) > 1 {
		paymentMethod = "Mixto"
	}
	item := invoice{ID: s.state.Counters.Invoice, Number: "SIGO-" + time.Now().Format("20060102") + "-" + strconv.Itoa(s.state.Counters.Invoice), PatientID: plan.PatientID, PatientName: plan.PatientName, TreatmentPlanID: plan.ID, Subtotal: plan.Total, Tax: tax, Total: total, PaymentMethod: paymentMethod, Payments: req.Payments, Status: "Pagada", FiscalStatus: "Pendiente de integración", CreatedAt: time.Now()}
	s.state.Counters.Invoice++
	s.state.Invoices = append(s.state.Invoices, item)
	s.auditLocked(r, "CREATE", "invoice", strconv.Itoa(item.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar la factura")
		return
	}
	writeJSON(w, 201, item)
}

// applyPayments completa la forma de pago. Si el cliente no detalla el desglose,
// el total se registra como un único pago entero para no sumarlo dos veces al
// calcular el efectivo esperado en el cierre de caja.
func applyPayments(req *invoiceRequest, total float64) {
	if len(req.Payments) == 0 {
		req.Payments = []paymentEntry{{Method: req.PaymentMethod, Amount: total}}
	}
}

func (s *server) expectedCashLocked(date string) float64 {
	total := 0.0
	for _, item := range s.state.Invoices {
		if item.CreatedAt.Format("2006-01-02") != date || item.Status == "Anulada" {
			continue
		}
		for _, payment := range item.Payments {
			if strings.EqualFold(payment.Method, "Efectivo") {
				total += payment.Amount
			}
		}
	}
	return total
}
func (s *server) createCashClosure(w http.ResponseWriter, r *http.Request) {
	var req closureRequest
	if decode(r, &req) != nil || !finiteNonNegative(req.DeclaredCash) {
		fail(w, 422, "VALIDATION_ERROR", "El conteo físico no es válido")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	date := time.Now().Format("2006-01-02")
	for _, item := range s.state.Closures {
		if item.Date == date {
			fail(w, 409, "ALREADY_CLOSED", "La caja de hoy ya fue cerrada")
			return
		}
	}
	expected := s.expectedCashLocked(date)
	item := cashClosure{ID: s.state.Counters.Closure, Date: date, DeclaredCash: req.DeclaredCash, ExpectedCash: expected, Difference: req.DeclaredCash - expected, ClosedBy: currentUser(r).Email, CreatedAt: time.Now()}
	s.state.Counters.Closure++
	s.state.Closures = append(s.state.Closures, item)
	s.auditLocked(r, "CLOSE", "cash_register", strconv.Itoa(item.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo cerrar la caja")
		return
	}
	writeJSON(w, 201, item)
}
func (s *server) listCashClosures(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Closures)
}

func (s *server) createSettlement(w http.ResponseWriter, r *http.Request) {
	var req settlementRequest
	if decode(r, &req) != nil || strings.TrimSpace(req.Professional) == "" || strings.TrimSpace(req.Period) == "" || req.Percentage <= 0 || req.Percentage > 100 || math.IsNaN(req.Percentage) || math.IsInf(req.Percentage, 0) {
		fail(w, 422, "VALIDATION_ERROR", "Liquidación inválida")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	base := 0.0
	for _, a := range s.state.Appointments {
		if a.Dentist == req.Professional && a.Status == "Atendido" && strings.HasPrefix(a.Start.Format("2006-01"), req.Period) {
			base += 45
		}
	}
	item := professionalSettlement{ID: s.state.Counters.Settlement, Professional: req.Professional, Period: req.Period, Amount: base * req.Percentage / 100, Status: "Pendiente", CreatedAt: time.Now()}
	s.state.Counters.Settlement++
	s.state.Settlements = append(s.state.Settlements, item)
	s.auditLocked(r, "CREATE", "professional_settlement", strconv.Itoa(item.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo crear la liquidación")
		return
	}
	writeJSON(w, 201, item)
}
func (s *server) listSettlements(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Settlements)
}

func (s *server) generalReport(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	revenue := 0.0
	for _, item := range s.state.Invoices {
		if item.Status != "Anulada" {
			revenue += item.Total
		}
	}
	attended, cancelled, noShow := 0, 0, 0
	for _, a := range s.state.Appointments {
		switch a.Status {
		case "Atendido":
			attended++
		case "Cancelado":
			cancelled++
		case "No Asistió":
			noShow++
		}
	}
	inventoryValue := 0.0
	for _, item := range s.state.Inventory {
		inventoryValue += item.Stock * item.UnitCost
	}
	writeJSON(w, 200, map[string]any{"revenue": revenue, "invoices": len(s.state.Invoices), "patients": len(s.state.Patients), "appointments": len(s.state.Appointments), "attended": attended, "cancelled": cancelled, "noShow": noShow, "inventoryValue": inventoryValue, "lowStock": s.lowStockCountLocked(), "auditEvents": len(s.state.Audit), "pendingMessages": len(s.state.Outbox)})
}
func (s *server) listAudit(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Audit)
}
func (s *server) listOutbox(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Outbox)
}
