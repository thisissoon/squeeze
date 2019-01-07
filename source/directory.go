package source

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go.soon.build/squeeze"
)

// Directory adds SQL files to a store from a filesystem directory
type Directory struct {
	Root    string
	Pattern string
}

// Source implements the Sourcer interface parsing files that
// match the Pattern in the Root and adding them to the Store
func (d *Directory) Source(store *squeeze.Store) error {
	root := strings.Split(d.Root, string(os.PathSeparator))
	for dir := range walk(d.Root, d.Pattern) {
		tpl, err := template.ParseGlob(filepath.Join(dir, d.Pattern))
		if err != nil {
			return err
		}
		path := strings.Split(dir, string(os.PathSeparator))
		name := strings.Join(path[len(root)-1:], squeeze.NamespaceSeperator)
		if err := store.Add(name, tpl); err != nil {
			return err
		}
	}
	return nil
}

// NewDirectory constructs a source to load SQL templates from a
// directory
func NewDirectory(root string) *Directory {
	return &Directory{
		Root:    root,
		Pattern: "*.*",
	}
}

// walks the path for files matching the pattern, returning
// a channel of the matched file names
func walk(path, pattern string) <-chan string {
	c := make(chan string)
	walker := func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() {
			log.Print(filepath.Join(path, pattern))
			files, _ := filepath.Glob(filepath.Join(path, pattern))
			if len(files) > 0 {
				c <- path
			}
		}
		return nil
	}
	go func() {
		defer close(c)
		_ = filepath.Walk(path, walker)
	}()
	return (<-chan string)(c)
}
