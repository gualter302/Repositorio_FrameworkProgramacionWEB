package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type stateStore interface {
	Load(context.Context) (appState, bool, error)
	Save(context.Context, appState) error
	Close() error
}

type memoryStore struct{}

func (memoryStore) Load(context.Context) (appState, bool, error) { return appState{}, false, nil }
func (memoryStore) Save(context.Context, appState) error         { return nil }
func (memoryStore) Close() error                                 { return nil }

type postgresStore struct{ db *sql.DB }

func openPostgresStore(databaseURL string) (*postgresStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS sigo_app_state (
		id SMALLINT PRIMARY KEY CHECK (id = 1), data JSONB NOT NULL, updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return &postgresStore{db: db}, nil
}

func (p *postgresStore) Load(ctx context.Context) (appState, bool, error) {
	var raw []byte
	err := p.db.QueryRowContext(ctx, `SELECT data FROM sigo_app_state WHERE id = 1`).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return appState{}, false, nil
	}
	if err != nil {
		return appState{}, false, err
	}
	var state appState
	if err := json.Unmarshal(raw, &state); err != nil {
		return appState{}, false, fmt.Errorf("estado guardado inválido: %w", err)
	}
	return state, true, nil
}

func (p *postgresStore) Save(ctx context.Context, state appState) error {
	raw, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = p.db.ExecContext(ctx, `INSERT INTO sigo_app_state(id, data, updated_at) VALUES (1, $1, NOW())
		ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data, updated_at = NOW()`, raw)
	return err
}
func (p *postgresStore) Close() error { return p.db.Close() }

func defaultState() appState {
	hash, err := bcrypt.GenerateFromPassword([]byte("Odonto2026!"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	day := func(hour, minute int) time.Time {
		return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.Local)
	}
	p1 := patient{ID: 1, Document: "0912345678", FirstName: "María", LastName: "González", BirthDate: "1988-04-12", Phone: "0991234567", Email: "maria@example.com", Insurance: "Salud Dental Plus", Allergies: []string{"Penicilina"}, Alerts: []string{"Alergia: Penicilina"}, CreatedAt: now}
	p2 := patient{ID: 2, Document: "0923456789", FirstName: "Carlos", LastName: "Mendoza", BirthDate: "1975-09-23", Phone: "0987654321", Email: "carlos@example.com", Conditions: []string{"Hipertensión"}, Anticoagulant: true, Alerts: []string{"Condición: Hipertensión", "Paciente anticoagulado"}, CreatedAt: now}
	p3 := patient{ID: 3, Document: "0934567890", FirstName: "Lucía", LastName: "Vera", BirthDate: "2001-01-08", Phone: "0971112233", Email: "lucia@example.com", CreatedAt: now}
	return appState{
		Users: map[string]user{
			"admin@sigodental.ec":      {ID: 1, Name: "Elena Torres", Email: "admin@sigodental.ec", Role: "administrador", PasswordHash: string(hash), MFAEnabled: true},
			"recepcion@sigodental.ec":  {ID: 2, Name: "Sofía León", Email: "recepcion@sigodental.ec", Role: "recepcionista", PasswordHash: string(hash)},
			"odontologo@sigodental.ec": {ID: 3, Name: "Andrea Ruiz", Email: "odontologo@sigodental.ec", Role: "odontologo", PasswordHash: string(hash), MFAEnabled: true},
			"asistente@sigodental.ec":  {ID: 4, Name: "Camila Reyes", Email: "asistente@sigodental.ec", Role: "asistente", PasswordHash: string(hash), MFAEnabled: true},
			"auditor@sigodental.ec":    {ID: 5, Name: "Mateo Silva", Email: "auditor@sigodental.ec", Role: "auditor", PasswordHash: string(hash), MFAEnabled: true},
		},
		Patients: []patient{p1, p2, p3},
		Appointments: []appointment{
			{ID: 1, PatientID: 1, PatientName: "María González", Dentist: "Dra. Andrea Ruiz", Chair: "Sillón 1", Start: day(9, 0), End: day(9, 45), Status: "Confirmado", Reason: "Profilaxis"},
			{ID: 2, PatientID: 2, PatientName: "Carlos Mendoza", Dentist: "Dr. Pablo Paz", Chair: "Sillón 2", Start: day(10, 30), End: day(11, 30), Status: "En Sala de Espera", Reason: "Endodoncia"},
			{ID: 3, PatientID: 3, PatientName: "Lucía Vera", Dentist: "Dra. Andrea Ruiz", Chair: "Sillón 1", Start: day(12, 0), End: day(12, 30), Status: "Programado", Reason: "Control"},
		},
		Odontogram: []odontogramEntry{
			{ID: 1, PatientID: 1, Tooth: 16, Surface: "oclusal", Condition: "Restauración existente", Nomenclature: "FDI", Notes: "Control de restauración en próxima cita", RecordedBy: "odontologo@sigodental.ec", RecordedAt: yesterday},
		},
		Periodontal: []periodontalRecord{
			{ID: 1, PatientID: 1, Tooth: 16, Depths: [6]int{2, 2, 3, 2, 2, 3}, RecordedAt: yesterday},
		},
		Notes: []clinicalNote{
			{ID: 1, PatientID: 1, Subjective: "Sensibilidad ocasional al frío", Objective: "Restauración estable en diente 16", Assessment: "Control preventivo sin signos de urgencia", Plan: "Profilaxis y control en seis meses", SignedBy: "odontologo@sigodental.ec", SignedAt: yesterday, LockedAt: yesterday.Add(24 * time.Hour), Amendments: []amendment{}},
		},
		Consents: []consent{
			{ID: 1, PatientID: 1, Procedure: "Profilaxis dental", Risks: "Sensibilidad temporal", PatientName: "María González", SignatureData: "María González", OTPValidated: true, SignedAt: yesterday},
		},
		Plans: []treatmentPlan{
			{ID: 1, PatientID: 1, PatientName: "María González", Status: "Aprobado", Items: []treatmentItem{{ID: 1, Procedure: "Profilaxis dental", Phase: 1, Sessions: 1, ProfessionalFee: 45, Status: "Agendable"}}, Total: 45, CreatedAt: yesterday},
			{ID: 2, PatientID: 2, PatientName: "Carlos Mendoza", Status: "Borrador", Items: []treatmentItem{{ID: 2, Procedure: "Endodoncia unirradicular", Phase: 1, Sessions: 2, ProfessionalFee: 180, LaboratoryCost: 20, Status: "Pendiente"}}, Total: 200, CreatedAt: now},
		},
		Inventory: []inventoryItem{
			{ID: 1, Name: "Guantes de nitrilo", Unit: "caja", Stock: 18, Minimum: 10, UnitCost: 7.5, UpdatedAt: now},
			{ID: 2, Name: "Resina fotocurable", Unit: "jeringa", Stock: 8, Minimum: 5, UnitCost: 24, UpdatedAt: now},
			{ID: 3, Name: "Anestesia local", Unit: "cartucho", Stock: 42, Minimum: 20, UnitCost: 1.8, UpdatedAt: now},
		},
		Invoices: []invoice{
			{ID: 1, Number: "SIGO-DEMO-0001", PatientID: 1, PatientName: "María González", TreatmentPlanID: 1, Subtotal: 45, Tax: 6.75, Total: 51.75, PaymentMethod: "Efectivo", Payments: []paymentEntry{{Method: "Efectivo", Amount: 51.75}}, Status: "Pagada", FiscalStatus: "Pendiente de integración", CreatedAt: yesterday},
		},
		Closures: []cashClosure{
			{ID: 1, Date: yesterday.Format("2006-01-02"), DeclaredCash: 51.75, ExpectedCash: 51.75, Difference: 0, ClosedBy: "recepcion@sigodental.ec", CreatedAt: yesterday},
		},
		Settlements: []professionalSettlement{
			{ID: 1, Professional: "Dra. Andrea Ruiz", Period: now.Format("2006-01"), Amount: 18, Status: "Pendiente", CreatedAt: now},
		},
		Outbox: []outboxMessage{
			{ID: 1, Channel: "email", Recipient: "lucia@example.com", Template: "appointment_confirmation", Status: "pending", CreatedAt: now},
		},
		Audit: []auditEntry{
			{ID: 1, UserEmail: "odontologo@sigodental.ec", Action: "CREATE", Entity: "clinical_note", EntityID: "1", IP: "demo", CreatedAt: yesterday},
			{ID: 2, UserEmail: "recepcion@sigodental.ec", Action: "CREATE", Entity: "invoice", EntityID: "1", IP: "demo", CreatedAt: yesterday},
			{ID: 3, UserEmail: "recepcion@sigodental.ec", Action: "CLOSE", Entity: "cash_register", EntityID: "1", IP: "demo", CreatedAt: yesterday},
		},
		Counters: counters{Patient: 4, Appointment: 4, Odontogram: 2, Periodontal: 2, Note: 2, Amendment: 1, Consent: 2, Plan: 3, PlanItem: 3, Inventory: 4, Invoice: 2, Closure: 2, Settlement: 2, Outbox: 2, Audit: 4},
	}
}

func normalizeState(state *appState) {
	if state.Users == nil {
		state.Users = map[string]user{}
	}
	if _, exists := state.Users["asistente@sigodental.ec"]; !exists {
		passwordHash := ""
		nextUserID := 1
		for _, current := range state.Users {
			if current.ID >= nextUserID {
				nextUserID = current.ID + 1
			}
		}
		if source, ok := state.Users["odontologo@sigodental.ec"]; ok {
			passwordHash = source.PasswordHash
		} else if source, ok := state.Users["recepcion@sigodental.ec"]; ok {
			passwordHash = source.PasswordHash
		}
		if passwordHash != "" {
			state.Users["asistente@sigodental.ec"] = user{ID: nextUserID, Name: "Camila Reyes", Email: "asistente@sigodental.ec", Role: "asistente", PasswordHash: passwordHash, MFAEnabled: true}
		}
	}
	if state.Patients == nil {
		state.Patients = []patient{}
	}
	if state.Appointments == nil {
		state.Appointments = []appointment{}
	}
	if state.Odontogram == nil {
		state.Odontogram = []odontogramEntry{}
	}
	if state.Periodontal == nil {
		state.Periodontal = []periodontalRecord{}
	}
	if state.Notes == nil {
		state.Notes = []clinicalNote{}
	}
	if state.Consents == nil {
		state.Consents = []consent{}
	}
	if state.Plans == nil {
		state.Plans = []treatmentPlan{}
	}
	if state.Inventory == nil {
		state.Inventory = []inventoryItem{}
	}
	if state.Invoices == nil {
		state.Invoices = []invoice{}
	}
	if state.Closures == nil {
		state.Closures = []cashClosure{}
	}
	if state.Settlements == nil {
		state.Settlements = []professionalSettlement{}
	}
	if state.Outbox == nil {
		state.Outbox = []outboxMessage{}
	}
	if state.Audit == nil {
		state.Audit = []auditEntry{}
	}
	for index := range state.Patients {
		if state.Patients[index].Conditions == nil {
			state.Patients[index].Conditions = []string{}
		}
		if state.Patients[index].Allergies == nil {
			state.Patients[index].Allergies = []string{}
		}
		if state.Patients[index].Alerts == nil {
			state.Patients[index].Alerts = []string{}
		}
	}
	for index := range state.Notes {
		if state.Notes[index].Amendments == nil {
			state.Notes[index].Amendments = []amendment{}
		}
	}
	for index := range state.Plans {
		if state.Plans[index].Items == nil {
			state.Plans[index].Items = []treatmentItem{}
		}
	}
	for index := range state.Invoices {
		if state.Invoices[index].Payments == nil {
			state.Invoices[index].Payments = []paymentEntry{}
		}
	}
}

func storeFromEnvironment() (stateStore, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return memoryStore{}, nil
	}
	return openPostgresStore(databaseURL)
}
