package views

import (
	"html/template"
)

// NewView is for assigning to the view type
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
	}
}

// View struct
type View struct {
	Template *template.Template
}
