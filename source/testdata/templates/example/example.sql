-- ./templates/example/example.sql
{{- define "byID" -}}
SELECT *
FROM "{{.Table}}"
WHERE id = $1
{{- end }}
