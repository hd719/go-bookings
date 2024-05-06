package models

import "github.com/hd719/go-bookings/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string // Cross Site Request Forgery token
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
