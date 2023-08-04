// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dao

const TemplateGenDaoDoContent = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. {TplCreatedAtDatetimeStr}
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"{{range .Imports}}
	"{{.}}"{{end}}
)

type {{.Name}} struct {
	g.Meta {{.SymbolQuota}}orm:"table:{{.TableName}}, do:true"{{.SymbolQuota}} {{range .Data}}
	{{.Name}} {{.DoType}} {{end}}
}
`
