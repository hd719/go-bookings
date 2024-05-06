package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/hd719/go-bookings/internal/config"
	"github.com/hd719/go-bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)

	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		// log.Fatal("Could not get template from template cache")
		return errors.New("cant get template from cache")
	}

	buf := new(bytes.Buffer)

	td := AddDefaultData(data, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// Keeps track of the pages the user has visited (so we do not have to read them from disk every time)
	myCache := map[string]*template.Template{}

	// Get all of the files named *.page.tmpl from the ./templates dir.
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Range through all the files ending with *.page.tmpl
	for _, page := range pages {
		// Get the name of the template (about.page.tmpl, contact.page.tmpl, etc.)
		name := filepath.Base(page)

		// Parse the data in the page variable and store that parsed data in the name variable
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Get any layouts that are in the templates dir. (base and footer layout)
		// Return a slice of strings with all the layouts
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			// Whats going on here?
			// This is saying that some of the files above ending in *.page.tmpl will require the layout files, so parse those files and add them to the template set
			_, err := ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
		// fmt.Println(myCache[name])
		// fmt.Println(myCache)
		// fmt.Println(ts)
	}

	return myCache, nil
}

// myCache output:
// {
// 	"about.page.tmpl":             "0x140003c9410",
// 	"contact.page.tmpl":           "0x140003e42a0",
// 	"generals.page.tmpl":          "0x140003e50b0",
// 	"home.page.tmpl":              "0x14000400090",
// 	"majors.page.tmpl":            "0x14000400ea0",
// 	"make-reservation.page.tmpl":  "0x14000401d40",
// 	"reservation-summary.page.tmpl": "0x1400042e150",
// 	"search-availability.page.tmpl": "0x1400042f440",
// }

// ts (template-set)

// &{<nil> 0x14000268f00 0x14000000ea0 0x140001107e0}
// &{<nil> 0x14000269600 0x14000001b00 0x14000110de0}
// &{<nil> 0x14000269d00 0x140002b25a0 0x140001113e0}
// &{<nil> 0x140002d4500 0x140002b30e0 0x14000111a40}
// &{<nil> 0x140002d4c00 0x140002b3b00 0x140002e00c0}
// &{<nil> 0x140002d5340 0x140002e86c0 0x140002e06c0}
// &{<nil> 0x14000384500 0x140002e9200 0x140002e1440}
// &{<nil> 0x14000384dc0 0x140002e9d40 0x140002e1c20}
