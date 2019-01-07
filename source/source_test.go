//go:generate go-bindata -pkg=templates -prefix=templates/ -o=source/testdata/templates.go ./source/testdata/templates/*
package source_test

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"go.soon.build/squeeze"
	"go.soon.build/squeeze/source"

	packr "github.com/gobuffalo/packr/v2"
	bindata "go.soon.build/squeeze/source/testdata"
)

var update = flag.Bool("update", false, "update .golden files")

func TestSource(t *testing.T) {
	testCases := []struct {
		desc   string
		source squeeze.Sourcer
	}{
		{
			desc: "packr",
			source: &source.Packr{
				Box: packr.New("sql", "./testdata/templates"),
			},
		},
		{
			desc: "directory",
			source: &source.Directory{
				Root:    "./testdata/templates",
				Pattern: "*.sql",
			},
		},
		{
			desc: "bindata",
			source: &source.Bindata{
				AssetNames: bindata.AssetNames(),
				Asset:      bindata.Asset,
				Root:       "./testdata/templates",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			sqt, err := squeeze.New(tc.source)
			if err != nil {
				t.Error(err)
			}
			tpl, err := sqt.Parse("example.byID", map[string]string{"Table": "examples"})
			if err != nil {
				t.Error(err)
			}
			golden := filepath.Join("testdata", "TestSource", "example.golden")
			if *update && err == nil {
				t.Log("update golden file")
				err := ioutil.WriteFile(golden, []byte(tpl), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}
			gb, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Log("compare template with golden file")
				t.Errorf("unexpected template; expected %s, got %s", gb, tpl)
			}
		})
	}
}
