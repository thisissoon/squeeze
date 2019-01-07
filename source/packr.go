package source

import (
	"os"
	"strings"
	"text/template"

	"go.soon.build/squeeze"

	packr "github.com/gobuffalo/packr/v2"
)

// Packr adds SQL files to a squeeze store from a packr.Box
type Packr struct {
	Box *packr.Box
}

// Source implements the Sourcer interface, parsing files from a packr Box
// and adding them to the template Store
func (s *Packr) Source(store *squeeze.Store) error {
	for _, tmpl := range s.Box.List() {
		data, err := s.Box.FindString(tmpl)
		if err != nil {
			return err
		}
		// replace path separators with namespaces
		parts := strings.Split(tmpl, string(os.PathSeparator))
		tmpl = strings.Join(parts[:len(parts)-1], squeeze.NamespaceSeperator)
		// parse template
		t := template.New(tmpl)
		t, err = t.Parse(string(data))
		if err != nil {
			return err
		}
		err = store.Add(tmpl, t)
		if err != nil {
			return err
		}
	}
	return nil
}
