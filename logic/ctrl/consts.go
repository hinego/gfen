package ctrl

const structInterface = `package {{.VersionName}}

import (
	"context"
)


type I{{.FileName}} interface {
	{{range .Functions}}{{.Name}}(ctx context.Context, req *{{.Name}}Req) (res *{{.Name}}Res, err error)
	{{end}}
}
`
const structTemplate = `
package {{.VersionName}}

import (
	"github.com/gogf/gf/v2/frame/g" {{range .Packages}}
	"{{.}}"{{end}}
)

type {{.Name}}Req struct {
	g.Meta ` + "`path:\"{{.Path}}\" tags:\"{{.Tags}}\" method:\"{{.Method}}\" summary:\"{{.Summary}}\"`" + ` {{range .Request}}
	{{.}} {{end}}
}
type {{.Name}}Res struct {
	g.Meta{{range .Response}}
	{{.}} {{end}}
}
`
const controllerTemplate = `package {{.ApiName}}

import (
	"context"
	{{range .Packages}}
	"{{.}}"{{end}}
	{{.FileName}}{{.VersionName | title}} "github.com/sucold/starter/api/{{.ApiName}}/{{.VersionName}}/{{.FileName}}"
)

func (r *controller{{title .VersionName}}) {{title .FunctionName}}(ctx context.Context, req *{{.FileName}}{{.VersionName | title}}.{{title .FunctionName}}Req) (res *{{.FileName}}{{.VersionName | title}}.{{title .FunctionName}}Res, err error) { {{if .Code}} 
	#code# {{else}}
	return nil, gerror.NewCode(gcode.CodeNotImplemented) {{end}}
}
`
const initTemplate = `package {{.ApiName}}

import (
	"github.com/sucold/starter/internal/controller"
)

func init() {
	controller.Register("{{.ApiName}}{{title .FileName}}{{title .VersionName}}", &controller{{title .VersionName}}{})
}

type controller{{title .VersionName}} struct{}
`
const controllerPackageTemplate = `package controller

import (
{{- range $api := .APIs}}
	{{$api.Name}}Face "github.com/sucold/starter/api/{{$api.Name}}"
{{- range $version := .Data}}{{- range $file := .Data}}
	{{$api.Name}}{{title $file.Name}}{{title $version.Name}}Api "github.com/sucold/starter/api/{{$api.Name}}/{{$version.Name}}/{{$file.Name}}"
{{- end}}{{- end}}{{- end}}
)

var (
{{- range $api := .APIs}}{{- range $version := .Data}}{{- range $file := .Data}}
	_{{$api.Name}}{{title $file.Name}}{{title $version.Name}} {{$api.Name}}{{title $file.Name}}{{title $version.Name}}Api.I{{title $file.Name}}
{{- end}}{{- end}}{{- end}}
)

func Register(name string, data any) {
	var ok bool
	switch name {
{{- range $api := .APIs}}{{- range $version := .Data}}{{- range $file := .Data}}
	case "{{$api.Name}}{{title $file.Name}}{{title $version.Name}}":
		if _{{$api.Name}}{{title $file.Name}}{{title $version.Name}},ok = data.({{$api.Name}}{{title $file.Name}}{{title $version.Name}}Api.I{{title $file.Name}});!ok {
			panic("{{$api.Name}}{{title $file.Name}}{{title $version.Name}} register error")
		}
		break
{{- end}}{{- end}}{{- end}}
	}
}

{{- range $api := .APIs}}
type {{$api.Name}} struct{}
func {{$api.Name | title}}() {{$api.Name}}Face.Interface {
	return &{{$api.Name}}{}
}
{{- end}}
{{- range $api := .APIs}}{{- range $version := .Data}}
func (r *{{$api.Name}}) {{title $version.Name}}() {{$api.Name}}Face.{{title $version.Name}}Interface {
	return &{{$api.Name}}{{title $version.Name}}{}
}
{{- end}}{{- end}}
{{- range $api := .APIs}}{{- range $version := .Data}}
type {{$api.Name}}{{title $version.Name}} struct{}
{{- range $file := .Data}}
func (d *{{$api.Name}}{{title $version.Name}}) {{title $file.Name}}() {{$api.Name}}{{title $file.Name}}{{title $version.Name}}Api.I{{title $file.Name}} {
	return _{{$api.Name}}{{title $file.Name}}{{title $version.Name}}
}
{{- end}}{{- end}}{{- end}}
`
const importControllerTemplate = `package packed

import (
{{- range $api := .APIs}}{{- range $version := .Data}}{{- range $file := .Data}}
	_ "github.com/sucold/starter/internal/controller/{{$api.Name}}/{{$file.Name}}/{{$version.Name}}"
{{- end}}{{- end}}{{- end}}
)
`
const interfaceTemplate = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package {{.ApiName}}

import (
{{range $version := .Versions}}
{{range .Data}}
	{{.Name}}{{$version.Name | title}} "github.com/sucold/starter/api/{{$.ApiName}}/{{$version.Name}}/{{.Name}}"
{{- end}}
{{- end}}
)

type Interface interface { {{range .Versions}}
	{{title .Name}}() {{title .Name}}Interface
{{- end}}
}

{{range $version := .Versions}}
type {{$version.Name | title}}Interface interface { {{range .Data}}
	{{.Name | title}}() {{.Name}}{{$version.Name | title}}.I{{.Name | title}}
    {{- end}}
}
{{- end}}
`
