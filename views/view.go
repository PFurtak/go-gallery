package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Variables for ease of updating layout template paths
var (

	// LayoutDir is file path that holds templates
	LayoutDir string = "views/layouts/"
	// TemplateDir is assigned the directory path views are stored in
	TemplateDir string = "views/"
	// TemplateExtension is the extension of our templates
	TemplateExtension string = ".gohtml"
)

// NewView is for assigning to the view type
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExtension(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View struct
type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := v.Render(rw, nil); err != nil {
		panic(err)
	}
}

// Render is used to render view with predefined layout
func (v *View) Render(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(rw, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)

	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath prepends template directory
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExtension appends template extension
func addTemplateExtension(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExtension
	}
}
