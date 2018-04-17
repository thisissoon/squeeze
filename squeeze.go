package squeeze

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// A ErrTemplateDefined is returned when a template of the same name already
// exists in the Store
type ErrTemplateDefined struct {
	Name string
}

// Error implements the error interface
func (e ErrTemplateDefined) Error() string {
	return fmt.Sprintf("%s already defined, pls choose a different name", e.Name)
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

// A DirectorySource adds SQL files to a store from a directory
type DirectorySource struct {
	Root    string
	Pattern string
}

// Source implements the Sourcer interface parsing files that
// match the Pattern in the Root and adding them to the Store
func (d *DirectorySource) Source(store *Store) error {
	rp := strings.Split(d.Root, string(os.PathSeparator))
	for file := range walk(d.Root, d.Pattern) {
		tpl, err := template.ParseFiles(file)
		if err != nil {
			return err
		}
		fp := strings.Split(file, string(os.PathSeparator))
		name := strings.Join(fp[len(rp):len(fp)-1], ".")
		if err := store.Add(name, tpl); err != nil {
			return err
		}
	}
	return nil
}

// Directory returns a DirectorySource to load SQL templates
// from a directory
func Directory(root string) *DirectorySource {
	return &DirectorySource{
		Root:    root,
		Pattern: ".*sql",
	}
}

// walk walks the path for files matching the pattern, returning
// a channel of the matched file names
func walk(path, pattern string) <-chan string {
	c := make(chan string)
	walker := func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() {
			files, _ := filepath.Glob(filepath.Join(path, "*.sql"))
			for _, file := range files {
				c <- file
			}
		}
		return nil
	}
	go func() {
		defer close(c)
		filepath.Walk(path, walker)
	}()
	return (<-chan string)(c)
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