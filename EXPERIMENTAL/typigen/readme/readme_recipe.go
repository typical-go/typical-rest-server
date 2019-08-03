package readme

import (
	"io"
	"text/template"
)

// ReadmeRecipe is recipe to generate readme
type ReadmeRecipe struct {
	Title       string
	Description string
	Sections    []SectionPogo
}

func (r ReadmeRecipe) Write(w io.Writer) (err error) {
	write(w, "<!-- Autogenerate by Typical-Go. DO NOT EDIT.--> \n\n")
	write(w, "# "+r.Title+"\n\n")
	write(w, r.Description+"\n\n")

	for _, section := range r.Sections {
		write(w, "## "+section.Title+"\n\n")

		if section.Data != nil {
			var tmpl *template.Template
			tmpl, err = template.New(section.Title).Parse(section.Content)
			if err != nil {
				return
			}

			tmpl.Execute(w, section.Data)
			write(w, "\n\n")
		} else {
			write(w, section.Content+"\n\n")
		}
	}

	return
}

func write(w io.Writer, s string) {
	w.Write([]byte(s))
}
