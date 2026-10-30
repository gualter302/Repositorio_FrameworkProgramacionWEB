package main

import (
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type consentRequest struct {
	Procedure     string `json:"procedure"`
	Risks         string `json:"risks"`
	PatientName   string `json:"patientName"`
	SignatureData string `json:"signatureData"`
	OTP           string `json:"otp"`
}
type amendmentRequest struct {
	Reason  string `json:"reason"`
	Content string `json:"content"`
}
type planStatusRequest struct {
	Status          string `json:"status"`
	ApprovedItemIDs []int  `json:"approvedItemIds"`
}
type periodontalRequest struct {
	Tooth     int   `json:"tooth"`
	Depths    []int `json:"depths"`
	Margin    int   `json:"margin"`
	Furcation int   `json:"furcation"`
	Mobility  int   `json:"mobility"`
	Bleeding  bool  `json:"bleeding"`
}

func (s *server) registerClinicalRoutes(mux *http.ServeMux) {
	clinical := []string{"administrador", "odontologo"}
	mux.Handle("GET /api/clinical/{id}", s.auth(http.HandlerFunc(s.getClinicalRecord), append(clinical, "asistente")...))
	mux.Handle("POST /api/clinical/{id}/odontogram", s.auth(http.HandlerFunc(s.createOdontogramEntry), clinical...))
	mux.Handle("POST /api/clinical/{id}/periodontogram", s.auth(http.HandlerFunc(s.createPeriodontalRecord), clinical...))
	mux.Handle("POST /api/clinical/{id}/notes", s.auth(http.HandlerFunc(s.createClinicalNote), clinical...))
	mux.Handle("POST /api/clinical/notes/{id}/amendments", s.auth(http.HandlerFunc(s.createAmendment), clinical...))
	mux.Handle("POST /api/clinical/{id}/consents", s.auth(http.HandlerFunc(s.createConsent), append(clinical, "recepcionista")...))
	mux.Handle("GET /api/treatment-plans", s.auth(http.HandlerFunc(s.listTreatmentPlans)))
	mux.Handle("POST /api/treatment-plans", s.auth(http.HandlerFunc(s.createTreatmentPlan), clinical...))
	mux.Handle("PATCH /api/treatment-plans/{id}/status", s.auth(http.HandlerFunc(s.changeTreatmentPlanStatus), "administrador", "odontologo", "recepcionista"))
}

func (s *server) patientExistsLocked(id int) (patient, bool) {
	for _, p := range s.state.Patients {
		if p.ID == id {
			return p, true
		}
	}
	return patient{}, false
}

func (s *server) getClinicalRecord(w http.ResponseWriter, r *http.Request) {
	patientID, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Paciente inválido")
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.patientExistsLocked(patientID)
	if !ok {
		fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
		return
	}
	odontogram := []odontogramEntry{}
	periodontal := []periodontalRecord{}
	notes := []clinicalNote{}
	consents := []consent{}
	plans := []treatmentPlan{}
	for _, item := range s.state.Odontogram {
		if item.PatientID == patientID {
			odontogram = append(odontogram, item)
		}
	}
	for _, item := range s.state.Periodontal {
		if item.PatientID == patientID {
			periodontal = append(periodontal, item)
		}
	}
	for _, item := range s.state.Notes {
		if item.PatientID == patientID {
			notes = append(notes, item)
		}
	}
	for _, item := range s.state.Consents {
		if item.PatientID == patientID {
			consents = append(consents, item)
		}
	}
	for _, item := range s.state.Plans {
		if item.PatientID == patientID {
			plans = append(plans, item)
		}
	}
	writeJSON(w, 200, map[string]any{"patient": p, "odontogram": odontogram, "periodontal": periodontal, "notes": notes, "consents": consents, "plans": plans})
}

func (s *server) createOdontogramEntry(w http.ResponseWriter, r *http.Request) {
	patientID, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Paciente inválido")
		return
	}
	var entry odontogramEntry
	if decode(r, &entry) != nil {
		fail(w, 400, "INVALID_JSON", "Registro inválido")
		return
	}
	surfaces := map[string]bool{"mesial": true, "distal": true, "oclusal": true, "incisal": true, "vestibular": true, "lingual": true, "palatino": true, "completo": true}
	entry.Condition = strings.TrimSpace(entry.Condition)
	entry.Notes = strings.TrimSpace(entry.Notes)
	if !validFDITooth(entry.Tooth) || !surfaces[strings.ToLower(entry.Surface)] || entry.Condition == "" || len([]rune(entry.Condition)) > 100 || len([]rune(entry.Notes)) > 500 {
		fail(w, 422, "VALIDATION_ERROR", "Diente, superficie o condición inválidos")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.patientExistsLocked(patientID); !ok {
		fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
		return
	}
	entry.ID = s.state.Counters.Odontogram
	s.state.Counters.Odontogram++
	entry.PatientID = patientID
	entry.Surface = strings.ToLower(entry.Surface)
	if entry.Nomenclature == "" {
		entry.Nomenclature = "FDI"
	}
	entry.RecordedBy = currentUser(r).Email
	entry.RecordedAt = time.Now()
	s.state.Odontogram = append(s.state.Odontogram, entry)
	s.auditLocked(r, "CREATE", "odontogram", strconv.Itoa(entry.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar el odontograma")
		return
	}
	writeJSON(w, 201, entry)
}

func (s *server) createPeriodontalRecord(w http.ResponseWriter, r *http.Request) {
	patientID, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Paciente inválido")
		return
	}
	var req periodontalRequest
	if decode(r, &req) != nil {
		fail(w, 400, "INVALID_JSON", "Registro inválido")
		return
	}
	if !validFDITooth(req.Tooth) {
		fail(w, 422, "VALIDATION_ERROR", "Diente inválido")
		return
	}
	if len(req.Depths) != 6 {
		fail(w, 422, "VALIDATION_ERROR", "Ingrese exactamente seis profundidades")
		return
	}
	for _, depth := range req.Depths {
		if depth < 0 || depth > 15 {
			fail(w, 422, "VALIDATION_ERROR", "Las profundidades deben estar entre 0 y 15 mm")
			return
		}
	}
	if req.Margin < -15 || req.Margin > 15 || req.Furcation < 0 || req.Furcation > 3 || req.Mobility < 0 || req.Mobility > 3 {
		fail(w, 422, "VALIDATION_ERROR", "Los valores periodontales están fuera del rango permitido")
		return
	}
	record := periodontalRecord{Tooth: req.Tooth, Margin: req.Margin, Furcation: req.Furcation, Mobility: req.Mobility, Bleeding: req.Bleeding}
	copy(record.Depths[:], req.Depths)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.patientExistsLocked(patientID); !ok {
		fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
		return
	}
	record.ID = s.state.Counters.Periodontal
	s.state.Counters.Periodontal++
	record.PatientID = patientID
	record.RecordedAt = time.Now()
	s.state.Periodontal = append(s.state.Periodontal, record)
	s.auditLocked(r, "CREATE", "periodontal_record", strconv.Itoa(record.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar el periodontograma")
		return
	}
	writeJSON(w, 201, record)
}

func (s *server) createClinicalNote(w http.ResponseWriter, r *http.Request) {
	patientID, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Paciente inválido")
		return
	}
	var note clinicalNote
	if decode(r, &note) != nil {
		fail(w, 400, "INVALID_JSON", "Nota inválida")
		return
	}
	note.Subjective = strings.TrimSpace(note.Subjective)
	note.Objective = strings.TrimSpace(note.Objective)
	note.Assessment = strings.TrimSpace(note.Assessment)
	note.Plan = strings.TrimSpace(note.Plan)
	if note.Subjective == "" || note.Objective == "" || note.Assessment == "" || note.Plan == "" || len([]rune(note.Subjective)) > 2000 || len([]rune(note.Objective)) > 2000 || len([]rune(note.Assessment)) > 2000 || len([]rune(note.Plan)) > 2000 {
		fail(w, 422, "VALIDATION_ERROR", "Complete los cuatro apartados SOAP")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.patientExistsLocked(patientID); !ok {
		fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
		return
	}
	now := time.Now()
	note.ID = s.state.Counters.Note
	s.state.Counters.Note++
	note.PatientID = patientID
	note.SignedBy = currentUser(r).Email
	note.SignedAt = now
	note.LockedAt = now.Add(24 * time.Hour)
	note.Amendments = []amendment{}
	s.state.Notes = append(s.state.Notes, note)
	s.auditLocked(r, "CREATE", "clinical_note", strconv.Itoa(note.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar la nota")
		return
	}
	writeJSON(w, 201, note)
}

func (s *server) createAmendment(w http.ResponseWriter, r *http.Request) {
	noteID, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Nota inválida")
		return
	}
	var req amendmentRequest
	if decode(r, &req) != nil || strings.TrimSpace(req.Reason) == "" || strings.TrimSpace(req.Content) == "" {
		fail(w, 422, "VALIDATION_ERROR", "Indique el motivo y la corrección")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, note := range s.state.Notes {
		if note.ID == noteID {
			item := amendment{ID: s.state.Counters.Amendment, Reason: req.Reason, Content: req.Content, SignedBy: currentUser(r).Email, CreatedAt: time.Now()}
			s.state.Counters.Amendment++
			s.state.Notes[i].Amendments = append(s.state.Notes[i].Amendments, item)
			s.auditLocked(r, "APPEND", "clinical_note", strconv.Itoa(noteID))
			if err := s.persistLocked(); err != nil {
				fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar la adenda")
				return
			}
			writeJSON(w, 201, item)
			return
		}
	}
	fail(w, 404, "NOT_FOUND", "Nota clínica no encontrada")
}

func (s *server) createConsent(w http.ResponseWriter, r *http.Request) {
	patientID, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Paciente inválido")
		return
	}
	var req consentRequest
	if decode(r, &req) != nil {
		fail(w, 400, "INVALID_JSON", "Consentimiento inválido")
		return
	}
	req.Procedure = strings.TrimSpace(req.Procedure)
	req.Risks = strings.TrimSpace(req.Risks)
	req.PatientName = strings.TrimSpace(req.PatientName)
	req.SignatureData = strings.TrimSpace(req.SignatureData)
	req.OTP = strings.TrimSpace(req.OTP)
	expected := os.Getenv("CONSENT_OTP")
	if expected == "" {
		expected = "654321"
	}
	if len([]rune(req.Procedure)) < 3 || len([]rune(req.Procedure)) > 160 || len([]rune(req.Risks)) < 3 || len([]rune(req.Risks)) > 2000 || !validPersonName(req.PatientName) || len([]rune(req.SignatureData)) < 3 || len([]rune(req.SignatureData)) > 120 || req.OTP != expected {
		fail(w, 422, "VALIDATION_ERROR", "Procedimiento, firma o código de validación inválidos")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.patientExistsLocked(patientID); !ok {
		fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
		return
	}
	item := consent{ID: s.state.Counters.Consent, PatientID: patientID, Procedure: req.Procedure, Risks: req.Risks, PatientName: req.PatientName, SignatureData: req.SignatureData, OTPValidated: true, SignedAt: time.Now()}
	s.state.Counters.Consent++
	s.state.Consents = append(s.state.Consents, item)
	s.auditLocked(r, "SIGN", "consent", strconv.Itoa(item.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar el consentimiento")
		return
	}
	writeJSON(w, 201, item)
}

func (s *server) listTreatmentPlans(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Plans)
}
func (s *server) createTreatmentPlan(w http.ResponseWriter, r *http.Request) {
	var plan treatmentPlan
	if decode(r, &plan) != nil {
		fail(w, 400, "INVALID_JSON", "Presupuesto inválido")
		return
	}
	if plan.PatientID < 1 || len(plan.Items) == 0 {
		fail(w, 422, "VALIDATION_ERROR", "Seleccione un paciente y al menos un procedimiento")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.patientExistsLocked(plan.PatientID)
	if !ok {
		fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
		return
	}
	plan.ID = s.state.Counters.Plan
	s.state.Counters.Plan++
	plan.PatientName = p.FirstName + " " + p.LastName
	plan.Status = "Borrador"
	plan.CreatedAt = time.Now()
	plan.Total = 0
	for i := range plan.Items {
		plan.Items[i].Procedure = strings.TrimSpace(plan.Items[i].Procedure)
		if len([]rune(plan.Items[i].Procedure)) < 3 || len([]rune(plan.Items[i].Procedure)) > 160 || plan.Items[i].Phase < 1 || plan.Items[i].Phase > 20 || plan.Items[i].Sessions < 1 || plan.Items[i].Sessions > 100 || !finiteNonNegative(plan.Items[i].ProfessionalFee) || !finiteNonNegative(plan.Items[i].LaboratoryCost) || !finiteNonNegative(plan.Items[i].Discount) || plan.Items[i].Discount > plan.Items[i].ProfessionalFee+plan.Items[i].LaboratoryCost || math.IsNaN(plan.Total) {
			fail(w, 422, "VALIDATION_ERROR", "Los datos del procedimiento no son válidos")
			return
		}
		plan.Items[i].ID = s.state.Counters.PlanItem
		s.state.Counters.PlanItem++
		plan.Items[i].Status = "Pendiente"
		line := plan.Items[i].ProfessionalFee + plan.Items[i].LaboratoryCost - plan.Items[i].Discount
		if line < 0 {
			line = 0
		}
		plan.Total += line
	}
	s.state.Plans = append(s.state.Plans, plan)
	s.auditLocked(r, "CREATE", "treatment_plan", strconv.Itoa(plan.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar el presupuesto")
		return
	}
	writeJSON(w, 201, plan)
}
func (s *server) changeTreatmentPlanStatus(w http.ResponseWriter, r *http.Request) {
	id, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "Presupuesto inválido")
		return
	}
	var req planStatusRequest
	if decode(r, &req) != nil {
		fail(w, 400, "INVALID_JSON", "Estado inválido")
		return
	}
	allowed := map[string]bool{"Borrador": true, "Aprobado Parcial": true, "Aprobado": true, "Rechazado": true}
	if !allowed[req.Status] {
		fail(w, 422, "INVALID_STATUS", "Estado no permitido")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, plan := range s.state.Plans {
		if plan.ID == id {
			if req.Status == "Aprobado Parcial" && len(req.ApprovedItemIDs) == 0 {
				fail(w, 422, "VALIDATION_ERROR", "Seleccione los procedimientos aprobados")
				return
			}
			s.state.Plans[i].Status = req.Status
			if strings.HasPrefix(req.Status, "Aprobado") {
				for j := range s.state.Plans[i].Items {
					if req.Status == "Aprobado" || slices.Contains(req.ApprovedItemIDs, s.state.Plans[i].Items[j].ID) {
						s.state.Plans[i].Items[j].Status = "Agendable"
					}
				}
			}
			s.auditLocked(r, "APPROVE", "treatment_plan", strconv.Itoa(id))
			if err := s.persistLocked(); err != nil {
				fail(w, 500, "PERSISTENCE_ERROR", "No se pudo actualizar el presupuesto")
				return
			}
			writeJSON(w, 200, s.state.Plans[i])
			return
		}
	}
	fail(w, 404, "NOT_FOUND", "Presupuesto no encontrado")
}
