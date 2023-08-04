// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dao

const TemplateGenDaoEntityContent = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package entity

import ( {{range .Imports}}
	"{{.}}"{{end}}
)
// {{.Name}} is the golang structure for table {{.TableName}}.
type {{.Name}} struct { {{range .Data}}
	{{.Name}} {{.FieldType}} {{end}}
}
{{range .Relations}} 
func (r *{{$.Name}}) Get{{.Name}}(force ...bool) (err error) {
	if r.{{.Name}} != nil && len(force) == 0 {
		return nil
	}
	var where = g.Map{
		"{{.Field}}": r.{{.Value}},
	}
	return g.DB().Model("{{.Table}}").Where(where).Scan(&r.{{.Name}})
}
{{end}}
`
const ModelContent = `package {{.Package}}

var Objects = []any{ {{range .Data}}
	&{{.}}{}, {{end}}
}
`
const MainContent = `package main

import (
	"github.com/hinego/gfen"
	"github.com/hinego/gfen/genx"
	"{{.Module}}/{{.TypePath}}"
	"log"
)

func main() {
	var err error
	err = gfen.Dao(&genx.DaoInput{
		DaoPath:   "{{.DaoPath}}",
		ModelPath: "{{.ModelPath}}",
		TypePath:  "{{.TypePath}}",
		Data:      {{.TypePath | basename}}.Objects,
	})
	log.Println(err)
}
`
