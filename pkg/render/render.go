package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/krzysztofkaptur/book-and-go/pkg/config"
	"github.com/krzysztofkaptur/book-and-go/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache
		tc = app.TemplateCache
	} else {
		tc, _ = CreateCacheTemplate()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("no template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateCacheTemplate() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all files named *.page.tmpl from './templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./layouts/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./layouts/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
