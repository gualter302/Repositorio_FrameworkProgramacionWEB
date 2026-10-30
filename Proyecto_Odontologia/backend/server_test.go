package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func loginToken(t *testing.T, handler http.Handler, email string) string {
	t.Helper()
	body := bytes.NewBufferString(`{"email":"` + email + `","password":"Odonto2026!"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/login", body)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("login: esperado 200, recibido %d", res.Code)
	}
	var data map[string]string
	_ = json.NewDecoder(res.Body).Decode(&data)
	return data["token"]
}

func loginTokenWithOTP(t *testing.T, handler http.Handler, email string) string {
	t.Helper()
	body := bytes.NewBufferString(`{"email":"` + email + `","password":"Odonto2026!","otp":"123456"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/login", body)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("login MFA: esperado 200, recibido %d: %s", res.Code, res.Body.String())
	}
	var data map[string]any
	_ = json.NewDecoder(res.Body).Decode(&data)
	return data["token"].(string)
}
func authorized(method, path, token, body string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	return req
}
func TestProtectedRoutesRequireSession(t *testing.T) {
	res := httptest.NewRecorder()
	newServer().routes().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/patients", nil))
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("esperado 401, recibido %d", res.Code)
	}
}

func TestNormalizeStateProducesJSONArrays(t *testing.T) {
	state := appState{
		Patients: []patient{{ID: 1}},
		Notes:    []clinicalNote{{ID: 1}},
		Plans:    []treatmentPlan{{ID: 1}},
		Invoices: []invoice{{ID: 1}},
	}

	normalizeState(&state)
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("no se pudo serializar el estado: %v", err)
	}
	if bytes.Contains(raw, []byte(":null")) {
		t.Fatalf("el estado normalizado todavía contiene valores null: %s", raw)
	}
}

func TestPatientCreatesVisibleMedicalAlerts(t *testing.T) {
	handler := newServer().routes()
	token := loginToken(t, handler, "recepcion@sigodental.ec")
	body := `{"document":"0999999999","firstName":"Ana","lastName":"Paz","birthDate":"1990-01-01","phone":"0990000000","allergies":["Látex"],"conditions":["Diabetes"],"anticoagulant":true}`
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, authorized(http.MethodPost, "/api/patients", token, body))
	if res.Code != http.StatusCreated {
		t.Fatalf("esperado 201, recibido %d: %s", res.Code, res.Body.String())
	}
	var p patient
	_ = json.NewDecoder(res.Body).Decode(&p)
	if len(p.Alerts) != 3 {
		t.Fatalf("se esperaban 3 alertas, recibido %+v", p.Alerts)
	}
}

func TestPatientValidationRejectsInvalidIdentityData(t *testing.T) {
	handler := newServer().routes()
	token := loginToken(t, handler, "recepcion@sigodental.ec")
	tests := []struct {
		name string
		body string
	}{
		{"números en nombres", `{"document":"0999999998","firstName":"Ana2","lastName":"Paz","birthDate":"1990-01-01","phone":"0990000000"}`},
		{"identificación con letras", `{"document":"09ABC99998","firstName":"Ana","lastName":"Paz","birthDate":"1990-01-01","phone":"0990000000"}`},
		{"teléfono corto", `{"document":"0999999998","firstName":"Ana","lastName":"Paz","birthDate":"1990-01-01","phone":"09900"}`},
		{"fecha futura", `{"document":"0999999998","firstName":"Ana","lastName":"Paz","birthDate":"2990-01-01","phone":"0990000000"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, authorized(http.MethodPost, "/api/patients", token, test.body))
			if response.Code != http.StatusUnprocessableEntity {
				t.Fatalf("esperado 422, recibido %d: %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestOdontogramRejectsInvalidFDITooth(t *testing.T) {
	handler := newServer().routes()
	token := loginTokenWithOTP(t, handler, "odontologo@sigodental.ec")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, authorized(http.MethodPost, "/api/clinical/1/odontogram", token, `{"tooth":19,"surface":"oclusal","condition":"Caries"}`))
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("esperado 422 para diente FDI inválido, recibido %d: %s", response.Code, response.Body.String())
	}
}

func TestAppointmentRejectsChairOverlap(t *testing.T) {
	handler := newServer().routes()
	token := loginToken(t, handler, "recepcion@sigodental.ec")
	day := time.Now()
	start := time.Date(day.Year(), day.Month(), day.Day(), 9, 15, 0, 0, time.Local)
	end := start.Add(30 * time.Minute)
	body, _ := json.Marshal(appointment{PatientID: 2, Dentist: "Dr. Test", Chair: "Sillón 1", Start: start, End: end, Reason: "Control"})
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, authorized(http.MethodPost, "/api/appointments", token, string(body)))
	if res.Code != http.StatusConflict {
		t.Fatalf("esperado 409 por sobreposición, recibido %d: %s", res.Code, res.Body.String())
	}
}
func TestAppointmentAllowsDifferentChair(t *testing.T) {
	handler := newServer().routes()
	token := loginToken(t, handler, "recepcion@sigodental.ec")
	day := time.Now()
	start := time.Date(day.Year(), day.Month(), day.Day(), 9, 15, 0, 0, time.Local)
	end := start.Add(30 * time.Minute)
	body, _ := json.Marshal(appointment{PatientID: 2, Dentist: "Dr. Test", Chair: "Sillón 3", Start: start, End: end, Reason: "Control"})
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, authorized(http.MethodPost, "/api/appointments", token, string(body)))
	if res.Code != http.StatusCreated {
		t.Fatalf("esperado 201, recibido %d: %s", res.Code, res.Body.String())
	}
}

func TestClinicalAndFinancialFlow(t *testing.T) {
	handler := newServer().routes()
	clinicalToken := loginTokenWithOTP(t, handler, "odontologo@sigodental.ec")
	receptionToken := loginToken(t, handler, "recepcion@sigodental.ec")

	noteBody := `{"subjective":"Dolor al masticar","objective":"Lesión en diente 16","assessment":"Caries profunda","plan":"Restauración"}`
	noteResponse := httptest.NewRecorder()
	handler.ServeHTTP(noteResponse, authorized(http.MethodPost, "/api/clinical/1/notes", clinicalToken, noteBody))
	if noteResponse.Code != http.StatusCreated {
		t.Fatalf("nota SOAP: esperado 201, recibido %d: %s", noteResponse.Code, noteResponse.Body.String())
	}
	var note clinicalNote
	_ = json.NewDecoder(noteResponse.Body).Decode(&note)

	amendmentResponse := httptest.NewRecorder()
	handler.ServeHTTP(amendmentResponse, authorized(http.MethodPost, "/api/clinical/notes/1/amendments", clinicalToken, `{"reason":"Aclaración","content":"La lesión corresponde a la superficie oclusal"}`))
	if amendmentResponse.Code != http.StatusCreated {
		t.Fatalf("adenda: esperado 201, recibido %d: %s", amendmentResponse.Code, amendmentResponse.Body.String())
	}
	if note.Subjective != "Dolor al masticar" {
		t.Fatal("la nota original fue alterada")
	}

	planBody := `{"patientId":1,"items":[{"procedure":"Restauración","phase":1,"sessions":1,"professionalFee":80,"laboratoryCost":10,"discount":5}]}`
	planResponse := httptest.NewRecorder()
	handler.ServeHTTP(planResponse, authorized(http.MethodPost, "/api/treatment-plans", clinicalToken, planBody))
	if planResponse.Code != http.StatusCreated {
		t.Fatalf("plan: esperado 201, recibido %d: %s", planResponse.Code, planResponse.Body.String())
	}

	approveResponse := httptest.NewRecorder()
	handler.ServeHTTP(approveResponse, authorized(http.MethodPatch, "/api/treatment-plans/1/status", receptionToken, `{"status":"Aprobado"}`))
	if approveResponse.Code != http.StatusOK {
		t.Fatalf("aprobación: esperado 200, recibido %d: %s", approveResponse.Code, approveResponse.Body.String())
	}

	invoiceResponse := httptest.NewRecorder()
	handler.ServeHTTP(invoiceResponse, authorized(http.MethodPost, "/api/invoices", receptionToken, `{"treatmentPlanId":1,"paymentMethod":"Efectivo"}`))
	if invoiceResponse.Code != http.StatusCreated {
		t.Fatalf("factura: esperado 201, recibido %d: %s", invoiceResponse.Code, invoiceResponse.Body.String())
	}
	var invoice invoice
	_ = json.NewDecoder(invoiceResponse.Body).Decode(&invoice)
	if invoice.Total != 97.75 || invoice.FiscalStatus != "Pendiente de integración" {
		t.Fatalf("totales o estado fiscal inesperados: %+v", invoice)
	}

	duplicateInvoice := httptest.NewRecorder()
	handler.ServeHTTP(duplicateInvoice, authorized(http.MethodPost, "/api/invoices", receptionToken, `{"treatmentPlanId":1,"paymentMethod":"Efectivo"}`))
	if duplicateInvoice.Code != http.StatusConflict {
		t.Fatalf("factura duplicada: esperado 409, recibido %d: %s", duplicateInvoice.Code, duplicateInvoice.Body.String())
	}

	closuresResponse := httptest.NewRecorder()
	handler.ServeHTTP(closuresResponse, authorized(http.MethodGet, "/api/cash-closures", receptionToken, ""))
	if closuresResponse.Code != http.StatusOK {
		t.Fatalf("consulta de cierres para recepción: esperado 200, recibido %d", closuresResponse.Code)
	}
}

func TestCompletedAppointmentConsumesInventory(t *testing.T) {
	app := newServer()
	handler := app.routes()
	token := loginToken(t, handler, "recepcion@sigodental.ec")
	before := app.state.Inventory[0].Stock
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, authorized(http.MethodPatch, "/api/appointments/1/status", token, `{"status":"Atendido"}`))
	if response.Code != http.StatusOK {
		t.Fatalf("esperado 200, recibido %d: %s", response.Code, response.Body.String())
	}
	if app.state.Inventory[0].Stock >= before {
		t.Fatal("el procedimiento atendido no descontó insumos")
	}
}

func TestAdministratorRequiresSecondFactor(t *testing.T) {
	handler := newServer().routes()
	withoutOTP := httptest.NewRecorder()
	handler.ServeHTTP(withoutOTP, httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"email":"admin@sigodental.ec","password":"Odonto2026!"}`)))
	if withoutOTP.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401 sin MFA, recibido %d", withoutOTP.Code)
	}
	if token := loginTokenWithOTP(t, handler, "admin@sigodental.ec"); token == "" {
		t.Fatal("no se emitió JWT después del segundo factor")
	}
}

func TestExpectedCashDoesNotDoubleCountInvoice(t *testing.T) {
	app := newServer()
	handler := app.routes()
	admin := loginTokenWithOTP(t, handler, "admin@sigodental.ec")
	clinicalToken := loginTokenWithOTP(t, handler, "odontologo@sigodental.ec")

	planBody := `{"patientId":1,"items":[{"procedure":"Restauración","phase":1,"sessions":1,"professionalFee":100,"laboratoryCost":0,"discount":0}]}`
	handler.ServeHTTP(httptest.NewRecorder(), authorized(http.MethodPost, "/api/treatment-plans", clinicalToken, planBody))
	handler.ServeHTTP(httptest.NewRecorder(), authorized(http.MethodPatch, "/api/treatment-plans/1/status", admin, `{"status":"Aprobado"}`))
	invoiceResponse := httptest.NewRecorder()
	handler.ServeHTTP(invoiceResponse, authorized(http.MethodPost, "/api/invoices", admin, `{"treatmentPlanId":1,"paymentMethod":"Efectivo"}`))
	if invoiceResponse.Code != http.StatusCreated {
		t.Fatalf("factura: esperado 201, recibido %d: %s", invoiceResponse.Code, invoiceResponse.Body.String())
	}

	closureResponse := httptest.NewRecorder()
	handler.ServeHTTP(closureResponse, authorized(http.MethodPost, "/api/cash-closures", admin, `{"declaredCash":115}`))
	if closureResponse.Code != http.StatusCreated {
		t.Fatalf("cierre: esperado 201, recibido %d: %s", closureResponse.Code, closureResponse.Body.String())
	}
	var closure cashClosure
	_ = json.NewDecoder(closureResponse.Body).Decode(&closure)
	if closure.ExpectedCash != 115 || closure.Difference != 0 {
		t.Fatalf("el efectivo esperado debe coincidir con la factura sin duplicarse: %+v", closure)
	}
}

func TestUpdateMissingPatientReturnsNotFound(t *testing.T) {
	handler := newServer().routes()
	token := loginToken(t, handler, "recepcion@sigodental.ec")
	body := `{"document":"0999999998","firstName":"Ana","lastName":"Paz","birthDate":"1990-01-01","phone":"0990000000"}`
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, authorized(http.MethodPut, "/api/patients/999", token, body))
	if response.Code != http.StatusNotFound {
		t.Fatalf("paciente inexistente: esperado 404, recibido %d: %s", response.Code, response.Body.String())
	}
}
