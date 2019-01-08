# squeeze
SQL template manager

Writing SQL is nice, embedding SQL in application code is not. Squeeze
provides an easy way to keep your SQL query templates separate.

```
go get go.soon.build/squeeze
```

1. Write templates
```sql
-- ./templates/user/user.sql
{{define "byID"}}
SELECT *
FROM "{{.Table}}"
WHERE id = $1
{{end}}
```
2. Load templates into squeeze store
```golang
import (
    "go.soon.build/squeeze"
    "go.soon.build/squeeze/source"
)

sqt, err := squeeze.New(source.NewDirectory("./templates"))
if err != nil {
    // handle err
}
```
3. Build a query from a template
```golang
qry := sqt.Parse("user.byID", struct{
    Table string
}, {
    Table: "page",
})
```

## Sourcing Templates

The directory source, as used above, reads templates from a directory tree.
To source templates from alternative locations (eg. bundled static files)
implement the `squeeze.Sourcer` interface to add templates to the store or
use one of the existing source implementations:

 - [FS Directory](https://github.com/thisissoon/squeeze/blob/master/source/directory.go)
 - [Go-Bindata](https://github.com/thisissoon/squeeze/blob/master/source/bindata.go)
 - [Packr](https://github.com/thisissoon/squeeze/blob/master/source/packr.go)
