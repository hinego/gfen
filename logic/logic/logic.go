package logic

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/ssr"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type sLogic struct {
	config *genx.LogicInput
	data   map[string]*genx.LogicData
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
			if err = r.serviceRegInit(logic); err != nil {
				return
			}
		}
	}
	log.Println("logic parse done")
	//gfile.PutContents("d.json", gjson.MustEncodeString(r.data))
	return err
}
func (r *sLogic) serviceRegInit(logic *genx.Logic) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: registerTemplate,
		File: fmt.Sprintf("%s/%s/%s.init.go", r.config.LogicPath, logic.Folder, logic.Base),
		Data: map[string]any{
			"Base": logic.Base,
			"Name": logic.Name,
			"Path": r.config.ServicePath,
		},
		Must: true,
	})
}
func (r *sLogic) serviceInterface(data *genx.LogicData) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: serviceTemplate,
		File: fmt.Sprintf("%s/%s.go", r.config.ServicePath, data.Name),
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
		File: in.ServicePath + "/face.init.go",
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
							for _, name := range field.Names {
								t, pkg := r.exprToString(field.Type, imports)
								returns = append(returns, name.Name+" "+t)
								if pkg != "" {
									pkgs = append(pkgs, pkg)
								}
							}
						}
					}

					logic.Funcs = append(logic.Funcs, &genx.FuncInfo{
						Name:       funcDecl.Name.Name,
						Parameters: params,
						Returns:    returns,
						Packages:   pkgs,
					})
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
