package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// We add a _ because the Valid function is a pointer receiver and belongs to the Form struct
func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid { // form has errors
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm) // we are not passing any form fields into the request

	has := form.Has("a") // bc we are not passing any form fields, a does not exist!

	if has {
		t.Error("Form shows 'has' field when it does not exist")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData) // re-initialize with our form with data

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	// Create a test that minlength doesn't work for non-existent field
	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows min length for a non-existent field")
	}

	// Checking for errors in our error array because validation failed
	isError := form.Errors.Get("x")

	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

	// Create a test that minlength doesn't work for values that is shorter than the given length
	postedValues := url.Values{}
	postedValues.Add("some-field", "some-value")
	form = New(postedValues)

	form.MinLength("some-field", 100)

	if form.Valid() {
		t.Error("shows min length of 100 met when data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another-field", "somevalue")
	form = New(postedValues)

	form.MinLength("another-field", 1)

	if !form.Valid() {
		t.Error("Shows minlength of 1 is not met, when the condition is met")
	}

	// Checking for 0 errors in our error array because validation passed
	isError = form.Errors.Get("another-field")

	if isError != "" {
		t.Error("Should not have an error, but did get one")
	}
}

func TestForm_Email(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")

	if form.Valid() {
		t.Error("Form shows valid email, when field is empty")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "me@here.com")
	form = New(postedValues)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("Form shows invalid email, when field is a valid email")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "invalid@email.")
	form = New(postedValues)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("Got an invalid email address")
	}
}
