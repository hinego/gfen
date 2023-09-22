package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/horm"
	"golang.org/x/mod/modfile"
)

type sGen struct {
	skip   bool
	module string
	Files  map[string]struct{}
}

func (r *sGen) Execute(in *genx.Execute) (err error) {
	if r.Files == nil {
		r.Files = map[string]struct{}{}
	}
	name := gfile.Abs(in.File)
	r.Files[name] = struct{}{}
	if gfile.Exists(in.File) && !in.Must {
		if !r.skip {
			log.Println("skipfile", in.File)
		}
		return nil
	}
	var (
		code    *template.Template
		buffer  bytes.Buffer
		dataMap = map[string]any{}
		funcMap = template.FuncMap{
			"title":    horm.ToName,
			"title2":   strings.Title,
			"lower":    strings.ToLower,
			"basename": gfile.Basename,
			"ToName":   horm.ToName,
			"jsonName": func(str string) string {
				matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
				matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
				snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
				snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
				return strings.ToLower(snake)
			},
		}
		data []byte
	)
	if !r.skip {
		log.Println("generate", in.File)
	}
	if len(in.Replace) > 0 {
		in.Code = gstr.ReplaceByMap(in.Code, in.Replace)
	}
	if code, err = template.New("code").Funcs(funcMap).Parse(in.Code); err != nil {
		return
	}
	if err = gconv.Scan(in.Data, &dataMap); err != nil {
		return err
	}
	dataMap["Module"] = r.GetModule()
	dataMap["SymbolQuota"] = "`"
	for k, v := range in.Map {
		dataMap[k] = v
	}
	if err = code.Execute(&buffer, dataMap); err != nil {
		return
	}

	if in.SkipFormat || strings.HasSuffix(in.File, ".go") {
		if data, err = format.Source(buffer.Bytes()); err != nil {
			if in.Debug {
				log.Println("格式化失败", in.File)
			}
			gfile.PutContents(in.File, buffer.String())
			return err
		}
		if in.Debug {
			log.Println(string(data))
		}
	} else {
		data = buffer.Bytes()
	}

	return gfile.PutContents(in.File, string(data))
}
func (r *sGen) ClearPath(keep string, paths ...string) {
	r.Clear(keep, "*.go", paths...)
}
func (r *sGen) Clear(keep string, pattern string, paths ...string) {
	if pattern == "" {
		pattern = "*.go"
	}
	for _, path := range paths {
		if files, err := gfile.ScanDirFile(path, pattern, true); err != nil {
			continue
		} else {
			for _, file := range files {
				if !strings.HasSuffix(file, keep) {
					if _, ok := r.Files[file]; !ok {
						gfile.Remove(file)
					}
				}
			}
		}
	}
	for _, path := range paths {
		if files, err := gfile.ScanDir(path, pattern, false); err != nil {
			continue
		} else {
			for _, file := range files {
				if !strings.HasSuffix(file, keep) {
					if _, ok := r.Files[file]; !ok {
						gfile.Remove(file)
					}
				}
			}
		}
	}
	for _, path := range paths {
		removeEmptyDirs(path)
	}
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
func (r *sGen) Path(paths ...string) string {
	paths = append([]string{r.GetModule()}, paths...)
	code := strings.Join(paths, "/")
	return strings.ReplaceAll(code, "//", "/")
}
func (r *sGen) Skip(skip bool) {
	r.skip = skip
}
func removeEmptyDirs(dir string) error {
	isEmpty, err := isDirEmpty(dir)
	if err != nil {
		return err
	}

	// 如果当前目录为空，删除它
	if isEmpty {
		fmt.Println("Removing:", dir)
		return os.Remove(dir)
	}

	// 否则，递归地检查其子目录
	subdirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, d := range subdirs {
		if d.IsDir() {
			subDirPath := filepath.Join(dir, d.Name())
			if err := removeEmptyDirs(subDirPath); err != nil {
				return err
			}
		}
	}

	// 再次检查当前目录是否为空，因为子目录可能已经被删除
	isEmpty, err = isDirEmpty(dir)
	if err != nil {
		return err
	}

	if isEmpty {
		fmt.Println("Removing:", dir)
		return os.Remove(dir)
	}

	return nil
}

// 检查目录是否为空
func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
