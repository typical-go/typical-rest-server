package readme

import (
	"io"
	"os"
	"text/template"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe"
)

// Readme detail
type Readme struct {
	Title       string
	Description string
	Sections    map[string]Section
	Titles      []string
}

// SetSection to set section
func (r *Readme) SetSection(title string, section Section) *Readme {
	if r.Sections == nil {
		r.Sections = make(map[string]Section)
	}
	_, ok := r.Sections[title]
	if !ok {
		r.Titles = append(r.Titles, title)
	}
	r.Sections[title] = section
	return r
}

func (r Readme) Write(w io.Writer) (err error) {
	write(w, "<!-- "+typirecipe.WaterMark+" -->\n\n")
	write(w, "# "+r.Title+"\n\n")
	write(w, r.Description+"\n\n")

	for _, title := range r.Titles {
		write(w, "## "+title+"\n\n")

		section := r.Sections[title]
		if section.Data != nil {
			var tmpl *template.Template
			tmpl, err = template.New(title).Parse(section.Content)
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

// WriteToFile to write the recipe to file
func (r Readme) WriteToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	return r.Write(file)
}

func write(w io.Writer, s string) {
	w.Write([]byte(s))
}
