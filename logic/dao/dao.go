package dao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/ssr"
	"go/ast"
	"go/parser"
	"go/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Relation struct {
	Type  string // hasOne, hasMany, belongsTo
	Table string // 从哪个表
	Field string // 条件字段
	Value string // 条件值字段
	Name  string // 关联的名称
}
type Schema struct {
	*schema.Schema
	Data     []*genx.FieldType
	Relation []*Relation
}
type sDao struct {
	db *gorm.DB
	*genx.DaoInput
	Map     *sync.Map
	schemas []*Schema
	Objects []string
	Import  map[string]string
}

func (r *sDao) doImport(data *Schema) []string {
	var (
		Imap    = map[string]struct{}{}
		objects = []string{}
	)
	for _, object := range data.Data {
		var doType = object.DoType
		if gstr.Contains(doType, ".") {
			// 获取object.DoType的包名
			var pkg = strings.Split(doType, ".")[0]
			pkg = gstr.ReplaceByMap(pkg, map[string]string{
				"*":  "",
				"[]": "",
			})
			if pack, ok := r.Import[pkg]; ok {
				Imap[pack] = struct{}{}
			} else {
				log.Println("pack not fund", pkg)
			}
		}
	}
	for k := range Imap {
		objects = append(objects, k)
	}
	return objects
}
func (r *sDao) entityImport(data *Schema) []string {
	var (
		Imap    = map[string]struct{}{}
		objects = []string{}
	)
	for _, object := range data.Data {
		var doType = object.FieldType
		if gstr.Contains(doType, ".") {
			// 获取object.DoType的包名
			var pkg = strings.Split(doType, ".")[0]
			pkg = gstr.ReplaceByMap(pkg, map[string]string{
				"*":  "",
				"[]": "",
			})
			if pack, ok := r.Import[pkg]; ok {
				Imap[pack] = struct{}{}
			} else {
				log.Println("pack not fund", pkg)
			}
		}
	}
	if len(data.Relation) != 0 {
		Imap["github.com/gogf/gf/v2/frame/g"] = struct{}{}
	}
	for k := range Imap {
		objects = append(objects, k)
	}
	return objects
}
func (r *sDao) connect() (err error) {
	const link = "host=192.168.32.130 user=postgres password=postgres dbname=sock port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var config = postgres.Config{
		DSN:                  link,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}
	r.db, err = gorm.Open(postgres.New(config), &gorm.Config{})
	return r.db.AutoMigrate(r.Data...)
}
func (r *sDao) syncSchema() (err error) {
	r.Map = &sync.Map{}
	for _, object := range r.Data {
		var data *schema.Schema
		if data, err = schema.Parse(object, r.Map, r.db.NamingStrategy); err != nil {
			return err
		}
		r.schemas = append(r.schemas, &Schema{
			Schema: data,
			Data:   make([]*genx.FieldType, 0),
		})
	}
	return
}
func (r *sDao) genDao(data *Schema) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: TemplateGenDaoIndexContent,
		Data: g.Map{
			"Name":    data.Name,
			"DaoPath": r.DaoPath,
		},
		File: fmt.Sprintf("%s/%s.go", r.DaoPath, strings.ToLower(data.Name)),
		Must: true,
	})
}
func (r *sDao) genDaoInternal(data *Schema) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: TemplateGenDaoInternalContent,
		Data: g.Map{
			"Name":      data.Name,
			"DaoPath":   r.DaoPath,
			"Group":     "default",
			"TableName": data.Table,
			"Data":      data.Data,
		},
		File: fmt.Sprintf("%s/internal/%s.go", r.DaoPath, strings.ToLower(data.Name)),
		Must: true,
	})
}
func (r *sDao) genModelDo(data *Schema) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: TemplateGenDaoDoContent,
		Data: g.Map{
			"Name":      data.Name,
			"TableName": data.Table,
			"Data":      data.Data,
			"Imports":   r.doImport(data),
		},
		File: fmt.Sprintf("%s/do/%s.go", r.ModelPath, strings.ToLower(data.Table)),
		Must: true,
	})
}
func (r *sDao) genModelEntity(data *Schema) (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: TemplateGenDaoEntityContent,
		Data: g.Map{
			"Name":      data.Name,
			"TableName": data.Table,
			"Data":      data.Data,
			"Relations": data.Relation,
			"Imports":   r.entityImport(data),
		},
		File:  fmt.Sprintf("%s/entity/%s.go", r.ModelPath, strings.ToLower(data.Table)),
		Must:  true,
		Debug: false,
	})
}

func (r *sDao) Generate(data *genx.DaoInput) (err error) {
	r.DaoInput = data
	if err = r.syncImports(); err != nil {
		return
	}
	if err = r.connect(); err != nil {
		return
	}
	if err = r.syncSchema(); err != nil {
		return
	}
	var basename = gfile.Basename(r.TypePath)
	for _, object := range r.schemas {
		var add = func(field *schema.Field) {
			var v = &genx.FieldType{
				Name:      field.Name,
				DBName:    field.DBName,
				FieldType: field.FieldType.String(),
				DataType:  string(field.DataType),
				DoType:    "any",
			}
			if gstr.Contains(v.DataType, "time") {
				v.DoType = "*gtime.Time"
				v.FieldType = "*gtime.Time"
			}
			v.FieldType = gstr.Replace(v.FieldType, basename+".", "")
			object.Data = append(object.Data, v)
		}
		for _, field := range object.Fields {
			if field.DBName == "" {
				continue
			}
			add(field)
		}
		for _, field := range object.Fields {
			if field.DBName != "" {
				continue
			}
			add(field)
		}
		for _, field := range object.Relationships.Relations {
			if field.Schema.Table != object.Table {
				continue
			}
			if len(field.References) != 1 {
				continue
			}
			var ref = field.References[0]
			rale := &Relation{
				Type:  string(field.Type),
				Table: field.FieldSchema.Table,
				Name:  field.Name,
			}
			log.Println(field.Field.Name, field.Field.DBName, field.Field.DataType, field.Field.FieldType)
			switch field.Type {
			case schema.HasOne, schema.HasMany:
				rale.Field = ref.ForeignKey.DBName
				rale.Value = ref.PrimaryKey.Name
			case schema.BelongsTo:
				rale.Field = ref.PrimaryKey.DBName
				rale.Value = ref.ForeignKey.Name
			default:
				continue
			}
			object.Relation = append(object.Relation, rale)
		}
		if err = r.genDao(object); err != nil {
			return
		}
		if err = r.genDaoInternal(object); err != nil {
			return
		}
		if err = r.genModelDo(object); err != nil {
			return
		}
		if err = r.genModelEntity(object); err != nil {
			return
		}
	}
	return
}
func (r *sDao) syncObjects() (err error) {
	return filepath.Walk(r.TypePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			if gstr.Contains(path, "types.go") {
				return nil
			}
			astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				//fmt.Printf("解析文件 %s 错误: %v\n", path, err)
				return nil
			}
			ast.Inspect(astFile, func(node ast.Node) bool {
				if typeSpec, ok := node.(*ast.TypeSpec); ok {
					//fmt.Printf("文件 %s 中的类型：%s\n", path, typeSpec.Name.Name)
					r.Objects = append(r.Objects, typeSpec.Name.Name)
				}
				return true
			})
		}
		return nil
	})
}
func (r *sDao) genObject() (err error) {
	var (
		name = gfile.Basename(r.TypePath)
	)
	return ssr.Gen().Execute(&genx.Execute{
		Code: ModelContent,
		Data: g.Map{
			"Package": name,
			"Data":    r.Objects,
		},
		File: fmt.Sprintf("%s/object_types.go", r.TypePath),
		Must: true,
	})
}
func (r *sDao) genMain() (err error) {
	var (
		path = filepath.Dir(r.TypePath)
	)
	return ssr.Gen().Execute(&genx.Execute{
		Code: MainContent,
		Data: r.DaoInput,
		File: fmt.Sprintf("%s/main.go", path),
		Must: true,
	})
}
func (r *sDao) GenModel(data *genx.DaoInput) (err error) {
	r.DaoInput = data
	if err = r.syncObjects(); err != nil {
		return
	}
	if err = r.genObject(); err != nil {
		return
	}
	return r.genMain()
}
func (r *sDao) syncImports() (err error) {
	r.Import = make(map[string]string)
	r.Import["gtime"] = "github.com/gogf/gf/v2/os/gtime"
	r.Import["g"] = "github.com/gogf/gf/v2/frame/g"
	return filepath.Walk(r.TypePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				fmt.Printf("解析文件 %s 错误: %v\n", path, err)
				return nil
			}
			for _, importSpec := range astFile.Imports {
				importPath := strings.Trim(importSpec.Path.Value, "\"")
				var baseName = gfile.Basename(importPath)
				r.Import[baseName] = importPath
			}
		}
		return nil
	})
}
