package main

import "time"

type user struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PasswordHash string `json:"passwordHash"`
	MFAEnabled   bool   `json:"mfaEnabled"`
}

type patient struct {
	ID               int       `json:"id"`
	Document         string    `json:"document"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	BirthDate        string    `json:"birthDate"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	Address          string    `json:"address"`
	EmergencyContact string    `json:"emergencyContact"`
	Guardian         string    `json:"guardian"`
	Insurance        string    `json:"insurance"`
	Conditions       []string  `json:"conditions"`
	Allergies        []string  `json:"allergies"`
	Surgeries        string    `json:"surgeries"`
	Anticoagulant    bool      `json:"anticoagulant"`
	Alerts           []string  `json:"alerts"`
	CreatedAt        time.Time `json:"createdAt"`
}

type appointment struct {
	ID          int       `json:"id"`
	PatientID   int       `json:"patientId"`
	PatientName string    `json:"patientName"`
	Dentist     string    `json:"dentist"`
	Chair       string    `json:"chair"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Status      string    `json:"status"`
	Reason      string    `json:"reason"`
}

type odontogramEntry struct {
	ID           int       `json:"id"`
	PatientID    int       `json:"patientId"`
	Tooth        int       `json:"tooth"`
	Surface      string    `json:"surface"`
	Condition    string    `json:"condition"`
	Nomenclature string    `json:"nomenclature"`
	Notes        string    `json:"notes"`
	RecordedBy   string    `json:"recordedBy"`
	RecordedAt   time.Time `json:"recordedAt"`
}

type periodontalRecord struct {
	ID         int       `json:"id"`
	PatientID  int       `json:"patientId"`
	Tooth      int       `json:"tooth"`
	Depths     [6]int    `json:"depths"`
	Margin     int       `json:"margin"`
	Furcation  int       `json:"furcation"`
	Mobility   int       `json:"mobility"`
	Bleeding   bool      `json:"bleeding"`
	RecordedAt time.Time `json:"recordedAt"`
}

type amendment struct {
	ID        int       `json:"id"`
	Reason    string    `json:"reason"`
	Content   string    `json:"content"`
	SignedBy  string    `json:"signedBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type clinicalNote struct {
	ID         int         `json:"id"`
	PatientID  int         `json:"patientId"`
	Subjective string      `json:"subjective"`
	Objective  string      `json:"objective"`
	Assessment string      `json:"assessment"`
	Plan       string      `json:"plan"`
	SignedBy   string      `json:"signedBy"`
	SignedAt   time.Time   `json:"signedAt"`
	LockedAt   time.Time   `json:"lockedAt"`
	Amendments []amendment `json:"amendments"`
}

type consent struct {
	ID            int       `json:"id"`
	PatientID     int       `json:"patientId"`
	Procedure     string    `json:"procedure"`
	Risks         string    `json:"risks"`
	PatientName   string    `json:"patientName"`
	SignatureData string    `json:"signatureData"`
	OTPValidated  bool      `json:"otpValidated"`
	SignedAt      time.Time `json:"signedAt"`
}

type treatmentItem struct {
	ID              int     `json:"id"`
	Procedure       string  `json:"procedure"`
	Phase           int     `json:"phase"`
	Sessions        int     `json:"sessions"`
	ProfessionalFee float64 `json:"professionalFee"`
	LaboratoryCost  float64 `json:"laboratoryCost"`
	Discount        float64 `json:"discount"`
	Status          string  `json:"status"`
}

type treatmentPlan struct {
	ID          int             `json:"id"`
	PatientID   int             `json:"patientId"`
	PatientName string          `json:"patientName"`
	Status      string          `json:"status"`
	Items       []treatmentItem `json:"items"`
	Total       float64         `json:"total"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type inventoryItem struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	Stock     float64   `json:"stock"`
	Minimum   float64   `json:"minimum"`
	UnitCost  float64   `json:"unitCost"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type invoice struct {
	ID              int            `json:"id"`
	Number          string         `json:"number"`
	PatientID       int            `json:"patientId"`
	PatientName     string         `json:"patientName"`
	TreatmentPlanID int            `json:"treatmentPlanId"`
	Subtotal        float64        `json:"subtotal"`
	Tax             float64        `json:"tax"`
	Total           float64        `json:"total"`
	PaymentMethod   string         `json:"paymentMethod"`
	Payments        []paymentEntry `json:"payments"`
	Status          string         `json:"status"`
	FiscalStatus    string         `json:"fiscalStatus"`
	CreatedAt       time.Time      `json:"createdAt"`
}

type paymentEntry struct {
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
}

type cashClosure struct {
	ID           int       `json:"id"`
	Date         string    `json:"date"`
	DeclaredCash float64   `json:"declaredCash"`
	ExpectedCash float64   `json:"expectedCash"`
	Difference   float64   `json:"difference"`
	ClosedBy     string    `json:"closedBy"`
	CreatedAt    time.Time `json:"createdAt"`
}

type professionalSettlement struct {
	ID           int       `json:"id"`
	Professional string    `json:"professional"`
	Period       string    `json:"period"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

type outboxMessage struct {
	ID        int       `json:"id"`
	Channel   string    `json:"channel"`
	Recipient string    `json:"recipient"`
	Template  string    `json:"template"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type auditEntry struct {
	ID        int       `json:"id"`
	UserEmail string    `json:"userEmail"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entityId"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
}

type counters struct {
	Patient, Appointment, Odontogram, Periodontal, Note, Amendment   int
	Consent, Plan, PlanItem, Inventory, Invoice, Closure, Settlement int
	Outbox, Audit                                                    int
}

type appState struct {
	Users        map[string]user          `json:"users"`
	Patients     []patient                `json:"patients"`
	Appointments []appointment            `json:"appointments"`
	Odontogram   []odontogramEntry        `json:"odontogram"`
	Periodontal  []periodontalRecord      `json:"periodontal"`
	Notes        []clinicalNote           `json:"notes"`
	Consents     []consent                `json:"consents"`
	Plans        []treatmentPlan          `json:"plans"`
	Inventory    []inventoryItem          `json:"inventory"`
	Invoices     []invoice                `json:"invoices"`
	Closures     []cashClosure            `json:"closures"`
	Settlements  []professionalSettlement `json:"settlements"`
	Outbox       []outboxMessage          `json:"outbox"`
	Audit        []auditEntry             `json:"audit"`
	Counters     counters                 `json:"counters"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
}
type statusRequest struct {
	Status string `json:"status"`
}
type apiError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
