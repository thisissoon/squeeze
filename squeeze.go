package squeeze

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// NamespaceSeperator seperates templates in different directories
var NamespaceSeperator = "."

// A ErrTemplateDefined is returned when a template of the same name already
// exists in the Store
type ErrTemplateDefined struct {
	Name string
}

// Error implements the error interface
func (e ErrTemplateDefined) Error() string {
	return fmt.Sprintf("%s already defined", e.Name)
}

// A ErrTemplateNotFound is returned when a template is not found be Parse
type ErrTemplateNotFound struct {
	Name string
}

// Error implements the error interface
func (e ErrTemplateNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Name)
}

// A Sourcer can Source templates and add them to the Store
type Sourcer interface {
	Source(s *Store) error
}

// The SourceFunc is an adapter allowing regular functions to act as Sourcers
type SourceFunc func(*Store) error

// Source implements the Sourcer interface calling the wrapped function
func (fn SourceFunc) Source(store *Store) error {
	return fn(store)
}

// A Store stores parsed SQL templates that can be executed by their name
type Store struct {
	templates map[string]*template.Template
}

// New constructs a New store
func New(sources ...Sourcer) (*Store, error) {
	s := &Store{
		templates: make(map[string]*template.Template),
	}
	if err := s.From(sources...); err != nil {
		return nil, err
	}
	return s, nil
}

// From loads templates in from any number of sources that
// implement the Sourcer interface
func (s *Store) From(sources ...Sourcer) error {
	for _, source := range sources {
		if err := source.Source(s); err != nil {
			return err
		}
	}
	return nil
}

// Add adds a template to the store, if a template already exists
// an ErrTemplateDefined will be returned
func (s *Store) Add(name string, tpl *template.Template) error {
	if _, ok := s.templates[name]; ok {
		return ErrTemplateDefined{name}
	}
	s.templates[name] = tpl
	return nil
}

// Parse parses a template
// foo.bar.list
func (s *Store) Parse(path string, v interface{}) (string, error) {
	var block string
	parts := strings.Split(path, NamespaceSeperator)
	block, parts = parts[len(parts)-1], parts[:len(parts)-1]
	t, ok := s.templates[strings.Join(parts, NamespaceSeperator)]
	if !ok {
		return "", ErrTemplateNotFound{path}
	}
	var w = new(bytes.Buffer)
	if err := t.ExecuteTemplate(w, block, v); err != nil {
		return "", err
	}
	return w.String(), nil
}

// String adds a template to the store from a raw SQL string
func String(name, sql string) Sourcer {
	return SourceFunc(func(store *Store) error {
		tpl := template.New(name)
		t, err := tpl.Parse(sql)
		if err != nil {
			return err
		}
		return store.Add(name, t)
	})
}
