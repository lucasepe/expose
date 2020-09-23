package slides

import (
	"html/template"
	"io"
	"os"
	"strings"

	vfstemplate "github.com/lucasepe/expose/vfs/html"
	"github.com/rakyll/statik/fs"

	// execute statik init() function
	_ "github.com/lucasepe/expose/statik"
)

const maxFileSize = 512 * 1024 // 512 MB

// Slides describes a markdown presentation.
type Slides struct {
	Title   string
	Content string
}

// FromFile load Markdown content from a file.
func FromFile(fn string) (*Slides, error) {
	fd, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return FromReader(fd)
}

// FromReader fetch Markdown content from a reader.
func FromReader(r io.Reader) (*Slides, error) {
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, io.LimitReader(r, maxFileSize)); err != nil {
		return nil, err
	}

	res := &Slides{Content: buf.String()}
	return res, nil
}

// Render fills the remark html templates with the doc info.
func (doc *Slides) Render(w io.Writer) error {
	sfs, err := fs.New()
	if err != nil {
		return err
	}
	// Parse and create the template
	tmpl := template.New("")
	if _, err := vfstemplate.ParseFiles(sfs, tmpl, "/boilerplate.html"); err != nil {
		return err
	}
	// Execute the template
	return tmpl.ExecuteTemplate(w, "boilerplate.html", doc)
}
