package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userContextKey contextKey = "user"

type server struct {
	mu            sync.RWMutex
	state         appState
	store         stateStore
	jwtSecret     []byte
	mfaSecret     []byte
	allowedOrigin string
	staticDir     string
	loginAttempts map[string]loginAttempt
}
type loginAttempt struct {
	Count   int
	ResetAt time.Time
}

func newServer() *server {
	s, err := newServerWithStore(memoryStore{}, []byte("clave-de-pruebas-segura-de-mas-de-32-caracteres"), []byte("mfa-pruebas-segura"))
	if err != nil {
		panic(err)
	}
	return s
}

func newServerWithStore(store stateStore, jwtSecret, mfaSecret []byte) (*server, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	state, found, err := store.Load(ctx)
	if err != nil {
		return nil, err
	}
	if !found {
		state = defaultState()
	}
	normalizeState(&state)
	if err := store.Save(ctx, state); err != nil {
		return nil, err
	}
	return &server{state: state, store: store, jwtSecret: jwtSecret, mfaSecret: mfaSecret, allowedOrigin: "http://localhost:5173", loginAttempts: map[string]loginAttempt{}}, nil
}

func newProductionServer() (*server, error) {
	store, err := storeFromEnvironment()
	if err != nil {
		return nil, err
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 && os.Getenv("DATABASE_URL") != "" {
		_ = store.Close()
		return nil, errors.New("JWT_SECRET debe tener al menos 32 caracteres")
	}
	if jwtSecret == "" {
		jwtSecret = "desarrollo-local-jwt-cambiar-antes-de-publicar"
	}
	mfaSecret := os.Getenv("MFA_SECRET")
	if len(mfaSecret) < 16 && os.Getenv("DATABASE_URL") != "" {
		_ = store.Close()
		return nil, errors.New("MFA_SECRET debe tener al menos 16 caracteres")
	}
	if mfaSecret == "" {
		mfaSecret = "desarrollo-mfa-local"
	}
	s, err := newServerWithStore(store, []byte(jwtSecret), []byte(mfaSecret))
	if err != nil {
		_ = store.Close()
		return nil, err
	}
	s.allowedOrigin = os.Getenv("ALLOWED_ORIGIN")
	if s.allowedOrigin == "" {
		if hostname := os.Getenv("RENDER_EXTERNAL_HOSTNAME"); hostname != "" {
			s.allowedOrigin = "https://" + hostname
		} else {
			s.allowedOrigin = "http://localhost:5173"
		}
	}
	s.staticDir = os.Getenv("STATIC_DIR")
	return s, nil
}

func (s *server) close() { _ = s.store.Close() }

func (s *server) startWorkers(ctx context.Context) {
	go func() {
		s.generateReminderSweep()
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.generateReminderSweep()
			}
		}
	}()
}

func (s *server) generateReminderSweep() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	added := false
	for _, appointment := range s.state.Appointments {
		if appointment.Status == "Cancelado" || appointment.Status == "Atendido" || appointment.Status == "No Asistió" {
			continue
		}
		hours := appointment.Start.Sub(now).Hours()
		window := ""
		if hours >= 23 && hours <= 25 {
			window = "24h"
		} else if hours >= 47 && hours <= 49 {
			window = "48h"
		} else {
			continue
		}
		recipient := ""
		for _, p := range s.state.Patients {
			if p.ID == appointment.PatientID {
				recipient = p.Email
				break
			}
		}
		if recipient == "" {
			continue
		}
		template := "appointment_reminder_" + window + "_" + strconv.Itoa(appointment.ID)
		exists := false
		for _, message := range s.state.Outbox {
			if message.Template == template && message.Recipient == recipient {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		s.state.Outbox = append(s.state.Outbox, outboxMessage{ID: s.state.Counters.Outbox, Channel: "email", Recipient: recipient, Template: template, Status: "pending", CreatedAt: now})
		s.state.Counters.Outbox++
		added = true
	}
	if added {
		_ = s.persistLocked()
	}
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("POST /api/login", s.login)
	mux.Handle("GET /api/dashboard", s.auth(http.HandlerFunc(s.dashboard)))
	mux.Handle("GET /api/patients", s.auth(http.HandlerFunc(s.listPatients)))
	mux.Handle("POST /api/patients", s.auth(http.HandlerFunc(s.createPatient), "administrador", "recepcionista", "odontologo"))
	mux.Handle("GET /api/patients/{id}", s.auth(http.HandlerFunc(s.getPatient)))
	mux.Handle("PUT /api/patients/{id}", s.auth(http.HandlerFunc(s.updatePatient), "administrador", "recepcionista", "odontologo"))
	mux.Handle("GET /api/appointments", s.auth(http.HandlerFunc(s.listAppointments)))
	mux.Handle("POST /api/appointments", s.auth(http.HandlerFunc(s.createAppointment), "administrador", "recepcionista", "odontologo"))
	mux.Handle("PATCH /api/appointments/{id}/status", s.auth(http.HandlerFunc(s.changeAppointmentStatus), "administrador", "recepcionista", "odontologo"))
	s.registerClinicalRoutes(mux)
	s.registerFinanceRoutes(mux)
	if s.staticDir != "" {
		mux.HandleFunc("/", s.serveFrontend)
	}
	return s.cors(s.securityHeaders(mux))
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok", "service": "sigo-dental"})
}
func (s *server) persistLocked() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.store.Save(ctx, s.state)
}
func (s *server) auditLocked(r *http.Request, action, entity, entityID string) {
	u := currentUser(r)
	entry := auditEntry{ID: s.state.Counters.Audit, UserEmail: u.Email, Action: action, Entity: entity, EntityID: entityID, IP: clientIP(r), CreatedAt: time.Now()}
	s.state.Counters.Audit++
	s.state.Audit = append(s.state.Audit, entry)
}
func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	return r.RemoteAddr
}
func currentUser(r *http.Request) user {
	value, _ := r.Context().Value(userContextKey).(user)
	return value
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func fail(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiError{Message: message, Code: code})
}
func decode(r *http.Request, value any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 2<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(value)
}
func idFrom(r *http.Request) (int, error) { return strconv.Atoi(r.PathValue("id")) }

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if decode(r, &req) != nil {
		fail(w, 400, "INVALID_JSON", "La solicitud no es válida")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	attemptKey := clientIP(r) + ":" + email
	s.mu.Lock()
	attempt := s.loginAttempts[attemptKey]
	if time.Now().After(attempt.ResetAt) {
		attempt = loginAttempt{ResetAt: time.Now().Add(10 * time.Minute)}
	}
	if attempt.Count >= 5 {
		s.mu.Unlock()
		fail(w, 429, "TOO_MANY_ATTEMPTS", "Espere unos minutos antes de volver a intentar")
		return
	}
	u, ok := s.state.Users[email]
	s.mu.Unlock()
	if !ok || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		s.mu.Lock()
		attempt.Count++
		s.loginAttempts[attemptKey] = attempt
		s.mu.Unlock()
		fail(w, 401, "INVALID_CREDENTIALS", "Correo o contraseña incorrectos")
		return
	}
	if u.MFAEnabled && !validMFA(s.mfaSecret, req.OTP) {
		s.mu.Lock()
		attempt.Count++
		s.loginAttempts[attemptKey] = attempt
		s.mu.Unlock()
		writeJSON(w, 401, map[string]any{"message": "Ingrese el código de verificación vigente", "code": "MFA_REQUIRED", "mfaRequired": true})
		return
	}
	token, err := signToken(s.jwtSecret, u)
	if err != nil {
		fail(w, 500, "TOKEN_ERROR", "No se pudo crear la sesión")
		return
	}
	s.mu.Lock()
	delete(s.loginAttempts, attemptKey)
	s.mu.Unlock()
	writeJSON(w, 200, map[string]any{"token": token, "name": u.Name, "email": u.Email, "role": u.Role, "expiresIn": 900})
}

func (s *server) auth(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		claims, err := verifyToken(s.jwtSecret, token)
		if err != nil {
			fail(w, 401, "UNAUTHORIZED", "La sesión no es válida o expiró")
			return
		}
		s.mu.RLock()
		u, ok := s.state.Users[claims.Email]
		s.mu.RUnlock()
		if !ok || u.ID != claims.Sub {
			fail(w, 401, "UNAUTHORIZED", "La cuenta ya no está disponible")
			return
		}
		if len(roles) > 0 && !slices.Contains(roles, u.Role) {
			fail(w, 403, "FORBIDDEN", "Su perfil no permite esta acción")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, u)))
	})
}

func (s *server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.allowedOrigin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (s *server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; connect-src 'self'; img-src 'self' data:")
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) serveFrontend(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		fail(w, 404, "NOT_FOUND", "Ruta no encontrada")
		return
	}
	clean := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if clean == "." || strings.HasPrefix(clean, "..") {
		clean = "index.html"
	}
	candidate := filepath.Join(s.staticDir, clean)
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		http.ServeFile(w, r, candidate)
		return
	}
	http.ServeFile(w, r, filepath.Join(s.staticDir, "index.html"))
}

func (s *server) dashboard(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	today := time.Now().Format("2006-01-02")
	appointments, waiting, alerts := 0, 0, 0
	next := []appointment{}
	for _, a := range s.state.Appointments {
		if a.Start.Format("2006-01-02") == today {
			appointments++
			if a.Status == "En Sala de Espera" {
				waiting++
			}
			if a.Status != "Cancelado" && a.Status != "No Asistió" && a.Status != "Atendido" {
				next = append(next, a)
			}
		}
	}
	slices.SortFunc(next, func(a, b appointment) int { return a.Start.Compare(b.Start) })
	if len(next) > 3 {
		next = next[:3]
	}
	for _, p := range s.state.Patients {
		if len(p.Alerts) > 0 {
			alerts++
		}
	}
	revenue := 0.0
	for _, invoice := range s.state.Invoices {
		if invoice.CreatedAt.Format("2006-01-02") == today && invoice.Status != "Anulada" {
			revenue += invoice.Total
		}
	}
	writeJSON(w, 200, map[string]any{"patients": len(s.state.Patients), "todayAppointments": appointments, "waiting": waiting, "criticalAlerts": alerts, "estimatedRevenue": revenue, "lowStock": s.lowStockCountLocked(), "pendingInvoices": s.pendingInvoicesLocked(), "nextAppointments": next})
}

func (s *server) listPatients(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Patients)
}
func (s *server) createPatient(w http.ResponseWriter, r *http.Request) {
	var p patient
	if decode(r, &p) != nil {
		fail(w, 400, "INVALID_JSON", "La ficha no es válida")
		return
	}
	if err := validatePatient(&p); err != nil {
		fail(w, 422, "VALIDATION_ERROR", err.Error())
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.state.Patients {
		if item.Document == p.Document {
			fail(w, 409, "DOCUMENT_EXISTS", "Ya existe un paciente con esa identificación")
			return
		}
	}
	p.ID = s.state.Counters.Patient
	s.state.Counters.Patient++
	p.CreatedAt = time.Now()
	s.state.Patients = append(s.state.Patients, p)
	s.auditLocked(r, "CREATE", "patient", strconv.Itoa(p.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar la ficha")
		return
	}
	writeJSON(w, 201, p)
}
func (s *server) getPatient(w http.ResponseWriter, r *http.Request) {
	id, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "El identificador no es válido")
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.state.Patients {
		if p.ID == id {
			writeJSON(w, 200, p)
			return
		}
	}
	fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
}
func (s *server) updatePatient(w http.ResponseWriter, r *http.Request) {
	id, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "El identificador no es válido")
		return
	}
	var p patient
	if decode(r, &p) != nil {
		fail(w, 400, "INVALID_JSON", "La ficha no es válida")
		return
	}
	if err := validatePatient(&p); err != nil {
		fail(w, 422, "VALIDATION_ERROR", err.Error())
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.patientExistsLocked(id); !ok {
		fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
		return
	}
	for _, patient := range s.state.Patients {
		if patient.ID != id && patient.Document == p.Document {
			fail(w, 409, "DOCUMENT_EXISTS", "Ya existe un paciente con esa identificación")
			return
		}
	}
	for i, item := range s.state.Patients {
		if item.ID == id {
			p.ID = id
			p.CreatedAt = item.CreatedAt
			s.state.Patients[i] = p
			s.auditLocked(r, "UPDATE", "patient", strconv.Itoa(id))
			if err := s.persistLocked(); err != nil {
				fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar la ficha")
				return
			}
			writeJSON(w, 200, p)
			return
		}
	}
	fail(w, 404, "NOT_FOUND", "Paciente no encontrado")
}
func (s *server) listAppointments(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, 200, s.state.Appointments)
}
func (s *server) createAppointment(w http.ResponseWriter, r *http.Request) {
	var a appointment
	if decode(r, &a) != nil {
		fail(w, 400, "INVALID_JSON", "La cita no es válida")
		return
	}
	a.Dentist = strings.TrimSpace(a.Dentist)
	a.Chair = strings.TrimSpace(a.Chair)
	a.Reason = strings.TrimSpace(a.Reason)
	allowedChairs := []string{"Sillón 1", "Sillón 2", "Sillón 3"}
	if a.PatientID < 1 || len([]rune(a.Dentist)) < 3 || len([]rune(a.Dentist)) > 80 || !slices.Contains(allowedChairs, a.Chair) || len([]rune(a.Reason)) < 3 || len([]rune(a.Reason)) > 120 || a.Start.IsZero() || a.End.IsZero() || !a.End.After(a.Start) || a.End.Sub(a.Start) > 12*time.Hour {
		fail(w, 422, "VALIDATION_ERROR", "Complete los datos y use un horario válido")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	found := false
	var patientEmail string
	for _, p := range s.state.Patients {
		if p.ID == a.PatientID {
			a.PatientName = p.FirstName + " " + p.LastName
			patientEmail = p.Email
			found = true
			break
		}
	}
	if !found {
		fail(w, 404, "PATIENT_NOT_FOUND", "El paciente no existe")
		return
	}
	for _, item := range s.state.Appointments {
		if item.Chair == a.Chair && item.Status != "Cancelado" && a.Start.Before(item.End) && a.End.After(item.Start) {
			fail(w, 409, "CHAIR_CONFLICT", fmt.Sprintf("%s ya está ocupado en ese horario", a.Chair))
			return
		}
	}
	a.ID = s.state.Counters.Appointment
	s.state.Counters.Appointment++
	if a.Status == "" {
		a.Status = "Programado"
	}
	s.state.Appointments = append(s.state.Appointments, a)
	if patientEmail != "" {
		s.state.Outbox = append(s.state.Outbox, outboxMessage{ID: s.state.Counters.Outbox, Channel: "email", Recipient: patientEmail, Template: "appointment_confirmation", Status: "pending", CreatedAt: time.Now()})
		s.state.Counters.Outbox++
	}
	s.auditLocked(r, "CREATE", "appointment", strconv.Itoa(a.ID))
	if err := s.persistLocked(); err != nil {
		fail(w, 500, "PERSISTENCE_ERROR", "No se pudo guardar la cita")
		return
	}
	writeJSON(w, 201, a)
}
func (s *server) changeAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	id, err := idFrom(r)
	if err != nil {
		fail(w, 400, "INVALID_ID", "El identificador no es válido")
		return
	}
	var req statusRequest
	if decode(r, &req) != nil {
		fail(w, 400, "INVALID_JSON", "El estado no es válido")
		return
	}
	allowed := []string{"Programado", "Confirmado", "En Sala de Espera", "En Atención", "Atendido", "Cancelado", "No Asistió"}
	if !slices.Contains(allowed, req.Status) {
		fail(w, 422, "INVALID_STATUS", "Estado de cita no permitido")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, a := range s.state.Appointments {
		if a.ID == id {
			s.state.Appointments[i].Status = req.Status
			if req.Status == "Atendido" {
				s.consumeRecipeLocked(a.Reason)
			}
			s.auditLocked(r, "UPDATE", "appointment", strconv.Itoa(id))
			if err := s.persistLocked(); err != nil {
				fail(w, 500, "PERSISTENCE_ERROR", "No se pudo actualizar la cita")
				return
			}
			writeJSON(w, 200, s.state.Appointments[i])
			return
		}
	}
	fail(w, 404, "NOT_FOUND", "Cita no encontrada")
}
