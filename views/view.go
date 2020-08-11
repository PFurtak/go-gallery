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
	// TemplateExtension is the extension of our templates
	TemplateExtension string = ".gohtml"
)

// NewView is for assigning to the view type
func NewView(layout string, files ...string) *View {
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

// Render is used to render view with predefined layout
func (v *View) Render(rw http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(rw, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)

	if err != nil {
		panic(err)
	}
	return files
}
