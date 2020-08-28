package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/Users/patrickfurtak/desktop/go-gallery/context"
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
	v.Render(rw, r, nil)
}

// Render is used to render view with predefined layout
func (v *View) Render(rw http.ResponseWriter, r *http.Request, data interface{}) error {
	rw.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		return (err)
	}
	io.Copy(rw, &buf)
	return nil
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
