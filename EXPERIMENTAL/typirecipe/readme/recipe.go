package readme

import (
	"io"
	"text/template"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe"
)

// Recipe is recipe to generate readme
type Recipe struct {
	Title       string
	Description string
	Sections    []Section
}

func (r Recipe) Write(w io.Writer) (err error) {
	write(w, "<!-- "+typirecipe.WaterMark+" -->\n\n")
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
