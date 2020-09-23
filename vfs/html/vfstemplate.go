package html

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/lucasepe/expose/vfs"
)

// ParseFiles creates a new Template if t is nil and parses the template definitions from
// the named files. The returned template's name will have the (base) name and
// (parsed) contents of the first file. There must be at least one file.
// If an error occurs, parsing stops and the returned *Template is nil.
func ParseFiles(fs http.FileSystem, t *template.Template, filenames ...string) (*template.Template, error) {
	return parseFiles(fs, t, filenames...)
}

// parseFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func parseFiles(fs http.FileSystem, t *template.Template, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("vfs/html/vfstemplate: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := vfs.ReadFile(fs, filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := path.Base(filename)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
