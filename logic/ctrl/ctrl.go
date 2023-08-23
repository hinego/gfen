package ctrl

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/ssr"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

type sCtrl struct {
	controller string
	api        string
	init       string
}

func (r *sCtrl) Generate(data *genx.ApiInput) (err error) {
	var (
		apis = data.Data
	)
	r.controller = data.Controller
	r.api = data.API
	r.init = data.Init
	if r.controller == "" {
		r.controller = "internal/controller"
	}
	if r.api == "" {
		r.api = "api"
	}
	if r.init == "" {
		r.init = "internal/packed"
	}
	if err = r.filter(apis); err != nil {
		return
	}
	if err = r.controllerInitImport(apis); err != nil {
		return
	}
	if err = r.controllerInitInterface(apis); err != nil {
		return
	}
	for _, api := range apis {
		if err = r.apiInterfaceInit(api); err != nil {
			return
		}
		for _, version := range api.Data {
			for _, file := range version.Data {
				var data = &genx.Fun{
					API:     api,
					Version: version,
					File:    file,
				}
				if err = r.apiFuncInterface(data); err != nil {
					return
				}
				if err = r.controllerInit(data); err != nil {
					return
				}
				for _, fun := range file.Data {
					var input = &genx.Fun{
						API:      api,
						Version:  version,
						File:     file,
						Function: fun,
					}
					if err = r.controllerFunc(input); err != nil {
						return
					}
					if err = r.apiFunStruct(input); err != nil {
						return
					}
				}
			}
		}
	}
	return
}
func (r *sCtrl) filter(apis []*genx.API) (err error) {
	for _, api := range apis {
		for _, version := range api.Data {
			for _, file := range version.Data {
				for _, fun := range file.Data {
					if fun.Path == "" {
						fun.Path = fmt.Sprintf("/%s/%s/%s/%s", version.Name, api.Name, file.Name, fun.Name)
					}
					if fun.Tags == "" {
						fun.Tags = strings.Title(file.Name)
					}
					if fun.Method == "" {
						fun.Method = "post"
					}
					if fun.Mime == "" {
						fun.Mime = "application/json"
					}
					if fun.Summary == "" {
						fun.Summary = strings.Title(fun.Name)
					}
					if fun.Description == "" {
						fun.Description = strings.Title(fun.Name)
					}
					fun.Name = cases.Title(language.English).String(fun.Name)
				}
			}
		}
	}
	return
}
func (r *sCtrl) controllerInitImport(apis []*genx.API) (err error) {
	var input = &genx.Execute{
		Code: importControllerTemplate,
		Data: map[string]any{
			"APIs": apis,
		},
		File: r.init + "/init_controller.go",
		Must: true,
	}
	return ssr.Gen().Execute(input)
}
func (r *sCtrl) controllerInitInterface(apis []*genx.API) (err error) {
	var input = &genx.Execute{
		Code: controllerPackageTemplate,
		Data: map[string]any{
			"APIs": apis,
		},
		File: r.controller + "/controller.go",
		Must: true,
	}
	return ssr.Gen().Execute(input)
}
func (r *sCtrl) apiInterfaceInit(api *genx.API) (err error) {
	var input = &genx.Execute{
		Code: interfaceTemplate,
		Data: map[string]interface{}{
			"ApiName":  api.Name,
			"Versions": api.Data,
		},
		File: fmt.Sprintf(r.api+"/%s/%s.go", api.Name, api.Name),
		Must: true,
	}
	return ssr.Gen().Execute(input)
}
func (r *sCtrl) apiFuncInterface(data *genx.Fun) (err error) {
	var (
		api     = data.API
		version = data.Version
		file    = data.File
		input   = &genx.Execute{
			Code: structInterface,
			Data: g.Map{
				"Functions":   file.Data,
				"VersionName": version.Name,
				"FileName":    strings.Title(file.Name),
			},
			File: fmt.Sprintf(r.api+"/%s/%s/%s/%s_%s_%s.go", api.Name, version.Name, file.Name, api.Name, version.Name, file.Name),
			Must: true,
		}
	)
	return ssr.Gen().Execute(input)
}
func (r *sCtrl) controllerInit(data *genx.Fun) (err error) {
	var (
		api     = data.API
		version = data.Version
		file    = data.File
		input   = &genx.Execute{
			Code: initTemplate,
			Data: g.Map{
				"ApiName":     api.Name,
				"FileName":    file.Name,
				"VersionName": version.Name,
			},
			File: fmt.Sprintf(r.controller+"/%s/%s/%s/%s_%s_%s.go", api.Name, file.Name, version.Name, api.Name, file.Name, version.Name),
		}
	)
	return ssr.Gen().Execute(input)
}
func (r *sCtrl) controllerFunc(data *genx.Fun) (err error) {
	var (
		api      = data.API
		version  = data.Version
		file     = data.File
		function = data.Function
		input    = &genx.Execute{
			Code: controllerTemplate,
			Data: map[string]any{},
			Map: map[string]any{
				"ApiName":      api.Name,
				"VersionName":  version.Name,
				"FileName":     file.Name,
				"FunctionName": function.Name,
				"Code":         function.Code,
				"Packages": []string{
					"github.com/gogf/gf/v2/errors/gcode",
					"github.com/gogf/gf/v2/errors/gerror",
				},
			},
			Replace: map[string]string{},
			Must:    function.Must,
			File:    fmt.Sprintf(r.controller+"/%s/%s/%s/%s_%s_%s_%s.go", api.Name, file.Name, version.Name, api.Name, file.Name, version.Name, function.Name),
		}
	)
	if function.Code != nil {
		input.Replace["#code#"] = function.Code.Code
		input.Map["Packages"] = function.Code.Import()
		for k, v := range function.Code.Data {
			input.Map[k] = v
		}
	}
	return ssr.Gen().Execute(input)
}
func (r *sCtrl) apiFunStruct(data *genx.Fun) (err error) {
	type StructInput struct {
		*genx.Function
		VersionName string
	}
	var (
		api     = data.API
		version = data.Version
		file    = data.File
		fun     = data.Function
		input   = &genx.Execute{
			Code: structTemplate,
			Data: StructInput{
				Function:    fun,
				VersionName: version.Name,
			},
			Map:  map[string]any{},
			File: fmt.Sprintf(r.api+"/%s/%s/%s/%s_%s_%s_%s.go", api.Name, version.Name, file.Name, api.Name, version.Name, file.Name, fun.Name),
			Must: fun.Must,
		}
	)

	if fun.Code != nil {
		Map := map[string]string{
			"{{FileName}}": strings.Title(file.Name),
		}
		input.Map = map[string]any{
			"Packages": fun.Code.ApiImport(),
			"Request":  fun.Code.Request(Map),
			"Response": fun.Code.Response(Map),
		}
	} else {
		input.Map = map[string]any{
			"Packages": make([]string, 0),
			"Request":  make([]string, 0),
			"Response": make([]string, 0),
		}
	}
	return ssr.Gen().Execute(input)
}
