package readme

import (
	"io"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe"
)

type Section func(md *typirecipe.Markdown) error

// Readme detail
type Readme struct {
	Title       string
	Description string
	Sections    map[string]Section
	Titles      []string
}

// NewReadme return new instance of Readme
func NewReadme() *Readme {
	return &Readme{}
}

// SetTitle to set the title
func (r *Readme) SetTitle(title string) *Readme {
	r.Title = title
	return r
}

// SetDescription to set the description
func (r *Readme) SetDescription(description string) *Readme {
	r.Description = description
	return r
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

// Output to write the output
func (r *Readme) Output(w io.Writer) (err error) {
	md := typirecipe.NewMarkdown(w)
	md.Comment(typirecipe.WaterMark)
	md.Heading1(r.Title)
	md.Writeln(r.Description)

	for _, title := range r.Titles {
		md.Heading2(title)

		err = r.Sections[title](md)
		if err != nil {
			return
		}
	}

	return
}

// OutputToFile to write the recipe to file
func (r *Readme) OutputToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	return r.Output(file)
}
