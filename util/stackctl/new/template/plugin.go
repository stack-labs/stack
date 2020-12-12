package template

var (
	Plugin = `package main
{{if .Plugins}}
import ({{range .Plugins}}
	_ "github.com/stack-labs/go-plugins/{{.}}"{{end}}
){{end}}
`
)
