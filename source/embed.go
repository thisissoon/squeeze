package source

import (
	"os"
	"strings"
	"text/template"

	"go.soon.build/squeeze"
)

type Embedded struct {
	Name string
	Data []byte
}

// Source implements the Sourcer interface, parsing files from a []byte
func (e *Embedded) Source(store *squeeze.Store) error {
	tmpl := e.Name

	// replace path separators with namespaces
	parts := strings.Split(tmpl, string(os.PathSeparator))
	tmpl = strings.Join(parts[:len(parts)-1], squeeze.NamespaceSeperator)

	// parse template
	t := template.New(tmpl)
	t, err := t.Parse(string(e.Data))
	if err != nil {
		return err
	}
	err = store.Add(tmpl, t)
	if err != nil {
		return err
	}
	return nil
}
