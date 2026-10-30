package main

import (
	"errors"
	"math"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

var (
	personNamePattern = regexp.MustCompile(`^[\p{L}]+(?:[ '\-][\p{L}]+)*$`)
	documentPattern   = regexp.MustCompile(`^[0-9]{10}$`)
	phonePattern      = regexp.MustCompile(`^[0-9]{10}$`)
)

func validPersonName(value string) bool {
	value = strings.TrimSpace(value)
	return len([]rune(value)) >= 2 && len([]rune(value)) <= 60 && personNamePattern.MatchString(value)
}

func validEmail(value string) bool {
	if value == "" {
		return true
	}
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && len(value) <= 120
}

func finiteNonNegative(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0
}

func normalizeTextList(values []string, maxItems int) ([]string, error) {
	if len(values) > maxItems {
		return nil, errors.New("hay demasiados elementos en la lista")
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if len([]rune(value)) > 80 {
			return nil, errors.New("cada condición o alergia debe tener máximo 80 caracteres")
		}
		result = append(result, value)
	}
	return result, nil
}

func validatePatient(p *patient) error {
	p.Document = strings.TrimSpace(p.Document)
	p.FirstName = strings.TrimSpace(p.FirstName)
	p.LastName = strings.TrimSpace(p.LastName)
	p.BirthDate = strings.TrimSpace(p.BirthDate)
	p.Phone = strings.TrimSpace(p.Phone)
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	p.Address = strings.TrimSpace(p.Address)
	p.EmergencyContact = strings.TrimSpace(p.EmergencyContact)
	p.Guardian = strings.TrimSpace(p.Guardian)
	p.Insurance = strings.TrimSpace(p.Insurance)
	p.Surgeries = strings.TrimSpace(p.Surgeries)

	if !validPersonName(p.FirstName) || !validPersonName(p.LastName) {
		return errors.New("nombres y apellidos solo pueden contener letras, espacios, apóstrofes o guiones")
	}
	if !documentPattern.MatchString(p.Document) {
		return errors.New("la identificación debe contener exactamente 10 dígitos")
	}
	if !phonePattern.MatchString(p.Phone) {
		return errors.New("el teléfono debe contener exactamente 10 dígitos")
	}
	birthDate, err := time.Parse("2006-01-02", p.BirthDate)
	if err != nil || birthDate.After(time.Now()) || birthDate.Before(time.Now().AddDate(-120, 0, 0)) {
		return errors.New("la fecha de nacimiento no es válida")
	}
	if !validEmail(p.Email) {
		return errors.New("el correo electrónico no es válido")
	}
	if len([]rune(p.Address)) > 160 || len([]rune(p.EmergencyContact)) > 100 || len([]rune(p.Guardian)) > 100 || len([]rune(p.Insurance)) > 100 || len([]rune(p.Surgeries)) > 500 {
		return errors.New("uno de los textos supera la longitud permitida")
	}
	p.Allergies, err = normalizeTextList(p.Allergies, 20)
	if err != nil {
		return err
	}
	p.Conditions, err = normalizeTextList(p.Conditions, 20)
	if err != nil {
		return err
	}
	p.Alerts = make([]string, 0, len(p.Allergies)+len(p.Conditions)+1)
	for _, allergy := range p.Allergies {
		p.Alerts = append(p.Alerts, "Alergia: "+allergy)
	}
	for _, condition := range p.Conditions {
		p.Alerts = append(p.Alerts, "Condición: "+condition)
	}
	if p.Anticoagulant {
		p.Alerts = append(p.Alerts, "Paciente anticoagulado")
	}
	return nil
}

func validFDITooth(tooth int) bool {
	quadrant, position := tooth/10, tooth%10
	return quadrant >= 1 && quadrant <= 4 && position >= 1 && position <= 8
}
