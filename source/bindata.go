package source

import (
	"os"
	"strings"
	"text/template"

	"go.soon.build/squeeze"
)

// Bindata adds SQL files to a store from go-bindata
type Bindata struct {
	AssetNames []string
	Asset      func(string) ([]byte, error)
	Root       string
}

// Source implements the Sourcer interface, parsing files that
// match the Pattern in the Root and adding them to the Store
func (s *Bindata) Source(store *squeeze.Store) error {
	root := strings.Split(s.Root, string(os.PathSeparator))
	for _, asset := range s.AssetNames {
		data, err := s.Asset(asset)
		if err != nil {
			return err
		}
		parts := strings.Split(asset, string(os.PathSeparator))
		name := strings.Join(parts[len(root):len(parts)-1], squeeze.NamespaceSeperator)
		t := template.New(name)
		t, err = t.Parse(string(data))
		if err != nil {
			return err
		}
		err = store.Add(name, t)
		if err != nil {
			return err
		}
	}
	return nil
}
