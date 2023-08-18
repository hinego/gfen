package logic

const serviceTemplate = `package {{.Package}}

import ({{range .Imports}}
	"{{.}}"{{end}}
)


type (
{{range .Data}}	I{{.Name}} interface { {{range .Funcs}}
		{{.Name}}({{range $index, $param := .Parameters}}{{if $index}}, {{end}}{{$param}}{{end}}) ({{range $index, $param := .Returns}}{{if $index}}, {{end}}{{$param}}{{end}}) {{end}}
	}
{{end}}	s{{.Name | title}} struct { 
		I{{.Source.Name}} {{range .Data}}{{if .Sub}}
		{{.Base}} I{{.Name}} {{end}}{{end}}
	}
)


var (
	local{{.Name | title}} = &s{{.Name | title}}{}
)
{{range .Data}}{{if .Sub}}
func (r *s{{$.Name | title}}) {{.Base | title}}() I{{.Name}} {
	if r.{{.Base}} == nil {
		if v, ok := faceMap["I{{.Name}}"]; ok {
			if v1, ok1 := v.(I{{.Name}}); ok1 {
				r.{{.Base}} = v1
			}
		}
	}
	if r.{{.Base}} == nil {
		panic("implement not found for interface I{{.Name}}, forgot register?")
	}
	return r.{{.Base}}
}
{{end}}{{end}}
func {{.Name | title}}() I{{.Name | title}} { 
	if local{{.Name | title}}.I{{.Source.Name}} == nil {
		if v, ok := faceMap["I{{.Source.Name}}"]; ok {
			if v1, ok1 := v.(I{{.Source.Name}}); ok1 {
				local{{.Name | title}}.I{{.Source.Name}} = v1
			}
		}
	}
	if local{{.Name | title}}.I{{.Source.Name}} == nil {
		panic("implement not found for interface I{{.Source.Name}}, forgot register?")
	} 
	return local{{$.Name | title}}
}
`
const registerTemplate = `package {{.Base}}

import "{{.Module}}/{{.Path}}"

func init() {
	service.Register("I{{.Name | title}}", &s{{.Base | title}}{})
}
`
const serviceInitTemplate = `package {{.Package}}

var faceMap = map[string]any{}

func Register(name string, face any) {
	faceMap[name] = face
}`
