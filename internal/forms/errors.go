package forms

// This will hold errors for our form
type errors map[string][]string

// Adds an error message for a given form field (html field)
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message) // ex. errors = {firstName: ["Some error message", "Another Error Message"]}
}

// Returns the first error message from the array of a given field
func (e errors) Get(field string) string {
	es := e[field] // ex. ["Some error message", "Another Error Message"]

	if len(es) == 0 {
		return ""
	}

	return es[0]
}
