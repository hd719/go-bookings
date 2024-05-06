package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Creates a custom form struct and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has check if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		// f.Errors.Add(field, "This field cannot be blank")
		return false
	}

	return true
}

// Check if the form is valid and we do this by checking the errors obj
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		v := f.Get(field)

		// if the form field has no value
		if strings.TrimSpace(v) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Checks for Min. length of string
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

// Checks validation for email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
