package logic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/ssr"
)

type sLogic struct {
	config  *genx.LogicInput
	data    map[string]*genx.LogicData
	imports map[string]string
}

func (r *sLogic) Parse(in *genx.LogicInput) (err error) {
	r.config = in
	if err = r.serviceInit(in); err != nil {
		return
	}
	if err = r.parseDir(in); err != nil {
		return
	}
	r.padding()
	for _, data := range r.data {
		if err = r.serviceInterface(data); err != nil {
			return
		}
		for _, logic := range data.Data {
			if len(logic.Funcs) == 0 {
				continue
			}
			if err = r.serviceRegInit(logic); err != nil {
				return
			}
		}
	}
	if err = r.serviceLogicInit(); err != nil {
		return err
	}
	log.Println("logic parse done")
	return err
}
func (r *sLogic) serviceRegInit(logic *genx.Logic) (err error) {
	if r.imports == nil {
		r.imports = make(map[string]string)
	}
	if logic.Main && !strings.HasSuffix(logic.Name, "_") {
		return
	}

	var data = map[string]any{
		"Base": logic.Base,
		"Name": logic.Name,
		"Path": r.config.ServicePath,
	}
	path := fmt.Sprintf("%s/%s", r.config.LogicPath, logic.Folder)
	r.imports[path] = path
	return ssr.Gen().Execute(&genx.Execute{
		Code:  registerTemplate,
		File:  fmt.Sprintf("%s/%s/%s.init.go", r.config.LogicPath, logic.Folder, logic.Base),
		Data:  data,
		Debug: strings.Contains(strings.ToLower(logic.Base), "cmd"),
	})
}
func (r *sLogic) serviceLogicInit() (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: importControllerTemplate,
		File: "internal/packed/logic.gen.go",
		Data: map[string]any{
			"Imports": r.imports,
		},
		Must: true,
	})
}
func (r *sLogic) serviceInterface(data *genx.LogicData) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: serviceTemplate,
		File: fmt.Sprintf("%s/%s.gen.go", r.config.ServicePath, data.Name),
		Must: true,
		Data: map[string]any{
			"Data":    data.Data,
			"Package": gfile.Basename(r.config.ServicePath),
			"Imports": data.Packages(),
			"Name":    data.Name,
			"Source":  data.Source,
			"Main":    data.Main,
		},
	})
}
func (r *sLogic) padding() {
	for _, data := range r.data {
		r.paddingData(data)
	}
}
func (r *sLogic) paddingData(data *genx.LogicData) {
	var (
		main   = data.GetMain()
		source = data.GetSource()
	)
	data.Main = main
	data.Source = source
	data.Data[main.Name] = main
	if source != nil {
		data.Data[source.Name] = source
	}
}
func (r *sLogic) serviceInit(in *genx.LogicInput) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: serviceInitTemplate,
		File: in.ServicePath + "/face.init.gen.go",
		Must: false,
		Data: map[string]any{
			"Package": gfile.Basename(r.config.ServicePath),
		},
	})
}
func (r *sLogic) parseDir(in *genx.LogicInput) (err error) {
	return filepath.WalkDir(in.LogicPath, r.walkDir)
}

func (r *sLogic) walkDir(path string, info os.DirEntry, e error) (err error) {
	if e != nil {
		return e
	}
	if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
		if err = r.parseFile(path); err != nil {
			return err
		}
	}
	return
}
func (r *sLogic) parseFile(file string) (err error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return
	}
	serv := ssr.Gen().Path(r.config.ServicePath)
	imports := make(map[string]string)
	for _, i := range node.Imports {
		// Trim quotes from paths
		path := i.Path.Value[1 : len(i.Path.Value)-1]
		name := filepath.Base(path)
		if i.Name != nil {
			name = i.Name.Name
		}
		imports[name] = path
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			_, ok = typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			var only = fmt.Sprintf("s%v", strings.Title(gfile.Basename(filepath.Dir(file))))
			if typeSpec.Name.Name != only {
				continue
			}
			if !strings.HasPrefix(typeSpec.Name.Name, "s") {
				continue
			}
			logic := &genx.Logic{
				Folder: gstr.ReplaceByArray(filepath.Dir(file), []string{
					"\\", "/",
					r.config.LogicPath + "/", "",
				}),
				Name: typeSpec.Name.Name,
			}
			for _, f := range node.Decls {
				funcDecl, ok := f.(*ast.FuncDecl)
				if !ok || funcDecl.Recv == nil {
					continue
				}
				for _, field := range funcDecl.Recv.List {
					starExpr, ok := field.Type.(*ast.StarExpr)
					if !ok {
						continue
					}

					ident, ok := starExpr.X.(*ast.Ident)
					if !ok || ident.Name != logic.Name {
						continue
					}

					params := []string{}
					pkgs := []string{}
					for _, field := range funcDecl.Type.Params.List {
						for _, name := range field.Names {
							t, pkg := r.exprToString(field.Type, imports)
							params = append(params, name.Name+" "+t)
							if pkg != "" {
								pkgs = append(pkgs, pkg)
							}
						}
					}

					returns := []string{}
					if funcDecl.Type.Results != nil {
						for _, field := range funcDecl.Type.Results.List {
							t, pkg := r.exprToString(field.Type, imports)
							if serv == pkg {
								pkg = ""
							}
							ss := gfile.Basename(r.config.ServicePath)
							t = strings.ReplaceAll(t, ss+".", "")
							if len(field.Names) > 0 {
								for _, name := range field.Names {
									returns = append(returns, name.Name+" "+t)
									if pkg != "" {
										pkgs = append(pkgs, pkg)
									}
								}
							} else {
								returns = append(returns, t)
								if pkg != "" {
									pkgs = append(pkgs, pkg)
								}
							}
						}
					}
					if funcDecl.Name.IsExported() {
						logic.Funcs = append(logic.Funcs, &genx.FuncInfo{
							Name:       funcDecl.Name.Name,
							Parameters: params,
							Returns:    returns,
							Packages:   pkgs,
						})
					}
				}
			}
			//data = append(data, logic)
			r.onLogic(logic)
		}
	}
	return
}
func (r *sLogic) onLogic(data *genx.Logic) {
	var (
		paths  = strings.Split(data.Folder, "/")
		base   = paths[0]
		latest = paths[len(paths)-1]
	)
	data.Base = latest
	data.Main = len(paths) == 1
	data.Sub = len(paths) > 1
	data.Name = r.formatPath(paths)
	r.addLogic(base, data)
	return
}
func (r *sLogic) addLogic(base string, data *genx.Logic) {
	if r.data == nil {
		r.data = make(map[string]*genx.LogicData)
	}
	if _, ok := r.data[base]; !ok {
		r.data[base] = &genx.LogicData{
			Name: base,
			Data: map[string]*genx.Logic{
				data.Name: data,
			},
		}
	} else {
		r.data[base].Data[data.Name] = data
	}
}
func (r *sLogic) exprToString(expr ast.Expr, imports map[string]string) (string, string) {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name, ""
	case *ast.Ellipsis:
		t, pkg := r.exprToString(e.Elt, imports)
		return "..." + t, pkg
	case *ast.SelectorExpr:
		ident, ok := e.X.(*ast.Ident)
		if !ok {
			return "unknown", ""
		}
		return ident.Name + "." + e.Sel.Name, imports[ident.Name]
	case *ast.StarExpr:
		t, pkg := r.exprToString(e.X, imports)
		return "*" + t, pkg
	case *ast.ArrayType:
		t, pkg := r.exprToString(e.Elt, imports)
		return "[]" + t, pkg
	case *ast.FuncType:
		// 处理函数类型
		params := []string{}
		for _, field := range e.Params.List {
			t, _ := r.exprToString(field.Type, imports)
			// 如果有多个名称与同一类型关联，例如: x, y int
			if len(field.Names) == 0 {
				params = append(params, t)
			} else {
				for _ = range field.Names {
					params = append(params, t)
				}
			}
		}
		results := []string{}
		if e.Results != nil {
			for _, field := range e.Results.List {
				t, _ := r.exprToString(field.Type, imports)
				if len(field.Names) == 0 {
					results = append(results, t)
				} else {
					for _ = range field.Names {
						results = append(results, t)
					}
				}
			}
		}
		return fmt.Sprintf("func(%s) %s", strings.Join(params, ", "), strings.Join(results, ", ")), ""
	case *ast.MapType:
		keyType, keyPkg := r.exprToString(e.Key, imports)
		valueType, valuePkg := r.exprToString(e.Value, imports)

		// 优先返回value的pkg，如果value没有包名再返回key的pkg
		if valuePkg != "" {
			return fmt.Sprintf("map[%s]%s", keyType, valueType), valuePkg
		}
		return fmt.Sprintf("map[%s]%s", keyType, valueType), keyPkg
	case *ast.ChanType:
		valueType, _ := r.exprToString(e.Value, imports)
		return "chan " + valueType, "" // 注意，这仅处理简单的channel类型，可能需要扩展以处理发送和接收的方向。
	case *ast.ParenExpr:
		return r.exprToString(e.X, imports)
	default:
		return "unknown", ""
	}
}
func (r *sLogic) formatPath(paths []string) string {
	var names = make([]string, 0)
	for _, v := range paths {
		names = append(names, strings.Title(v))
	}
	return strings.Join(names, "")
}
