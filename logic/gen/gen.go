package gen

import (
	"bytes"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/hinego/gfen/genx"
	"go/format"
	"golang.org/x/mod/modfile"
	"log"
	"strings"
	"text/template"
)

type sGen struct {
	module string
}

func (r *sGen) Execute(in *genx.Execute) (err error) {
	if gfile.Exists(in.File) && !in.Must {
		log.Println("skipfile", in.File)
		return nil
	}
	//asd
	var (
		code    *template.Template
		buffer  bytes.Buffer
		dataMap = map[string]any{}
		funcMap = template.FuncMap{
			"title":    strings.Title,
			"lower":    strings.ToLower,
			"basename": gfile.Basename,
		}
		data []byte
	)
	if code, err = template.New("code").Funcs(funcMap).Parse(in.Code); err != nil {
		return
	}
	if err = gconv.Scan(in.Data, &dataMap); err != nil {
		return err
	}
	dataMap["Module"] = r.GetModule()
	dataMap["SymbolQuota"] = "`"
	if err = code.Execute(&buffer, dataMap); err != nil {
		return
	}
	log.Println("generate", in.File)
	if data, err = format.Source(buffer.Bytes()); err != nil {
		return err
	}
	if in.Debug {
		log.Println(string(data))
	}
	return gfile.PutContents(in.File, string(data))
}

func (r *sGen) GetModule() string {
	if r.module != "" {
		return r.module
	}
	file, err := modfile.Parse("go.mod", gfile.GetBytes("go.mod"), nil)
	if err != nil {
		log.Fatal(err)
	}
	r.module = file.Module.Mod.Path
	return file.Module.Mod.Path
}
