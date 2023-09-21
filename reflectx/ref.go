package reflectx

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	_ "github.com/hinego/gfen/logic"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/mod/modfile"
)

var module string
var Mapping = map[string]string{
	"big.Int":               "string",
	"decimal.Decimal":       "string",
	"soft_delete.DeletedAt": "number",
}
var TypeMapping = map[string]any{}

func SetTypeMapping(m map[string]any) {
	for k, v := range m {
		TypeMapping[k] = v
	}
}
func IsDecimal(goType string) bool {
	switch goType {
	case "big.Int", "decimal.Decimal":
		return true
	default:
		return false
	}
}

var ImportMap = map[string]string{
	"decimal": "import { Decimal } from 'decimal.js';",
}
var Namer = func(s string) string {
	return s
}

func SetMapping(m map[string]string) {
	for k, v := range m {
		Mapping[k] = v
	}
}
func SetImport(m map[string]string) {
	for k, v := range m {
		ImportMap[k] = v
	}
}

type Table struct {
	Name  string
	Table string
	Key   string
	Enums []*Enum
}
type Enum struct {
	Name       string
	Value      any
	Desc       string
	Typescript string
}

var EnumMap = map[string][]*Table{}

type Field struct {
	Json       string   `json:"json,omitempty"`
	Name       string   `json:"name"`
	Type       string   `json:"type,omitempty"`
	Typescript string   `json:"typescript,omitempty"`
	Package    string   `json:"package,omitempty"`
	Enum       string   `json:"enum,omitempty"`
	Path       string   `json:"path,omitempty"`
	Data       []*Field `json:"data,omitempty"`
	Import     string   `json:"import,omitempty"`
	Array      bool     `json:"array,omitempty"`
	Decimal    bool     `json:"decimal,omitempty"`
	Optional   bool     `json:"optional,omitempty"`
	Desc       string   `json:"desc,omitempty"`
	Table      string   `json:"table,omitempty"`
}

func (r *Field) IsOptional() bool {
	var i = 0
	for _, v := range r.Data {
		if !v.Optional {
			i += 1
		}
	}
	return i == 0
}

func (r *Field) Have() bool {
	return len(r.Data) != 0
}
func (r *Field) TypeName() string {
	var arr = strings.Split(r.Type, ".")
	Name := arr[len(arr)-1]
	if len(r.Data) == 0 && r.Name == "" {
		return Name
	}
	if len(r.Data) == 0 {
		return r.Typescript
	}
	return Name
}
func (r *Field) TypeNameArray() string {
	var name = r.TypeName()
	if r.Array {
		return name + "[]"
	}
	return name
}
func Fields(field *Field) []*Field {
	var data = []*Field{}
	var arr = map[string]*Field{}
	for _, v := range field.Data {
		if len(v.Data) != 0 {
			arr[v.TypeName()] = v
			for _, vv := range Fields(v) {
				arr[vv.TypeName()] = vv
			}
		}
	}
	for _, v := range arr {
		data = append(data, v)
	}
	return data
}

type FunName struct {
	Version string
	API     string
	File    string
	Fun     string
}
type Func struct {
	FunName
	Path   string    `json:"path"`
	Method string    `json:"method"`
	In     *Field    `json:"in"`
	Out    *Field    `json:"out"`
	Func   *Function `json:"func"`
}

func (r *Func) Fields() []*Field {
	var data = []*Field{}
	data = append(data, Fields(r.In)...)
	data = append(data, Fields(r.Out)...)
	data = append(data, r.In)
	data = append(data, r.Out)
	return data
}

type Object struct {
	FunName
	Name   string   `json:"name"`
	Func   []*Func  `json:"func"`
	Enum   []*Table `json:"enum"`
	Data   *Field   `json:"data"`
	Create *Field   `json:"create"`
	Update *Field   `json:"update"`
}

func (r *Object) Sync() {
	for _, v := range r.Func {
		r.FunName = v.FunName
		break
	}
	for _, v := range r.Func {
		if r.Data == nil {
			r.Data = v.Func.GetData()
		}
		if r.Create == nil {
			r.Create = v.Func.GetCreate()
		}
		if r.Update == nil {
			r.Update = v.Func.GetUpdate()
		}
	}
}
func (r *Object) Fields() []*Field {
	var data = []*Field{}
	var exists = map[string]bool{}
	for _, v := range r.Func {
		var d2 = v.Fields()
		for _, v2 := range d2 {
			if _, ok := exists[v2.TypeName()]; !ok {
				exists[v2.TypeName()] = true
				data = append(data, v2)
			}
		}
	}
	return data
}

type Function struct {
	FunName
	In  reflect.Type
	Out reflect.Type
}

func (r *Function) GetData() *Field {
	var fun = strings.ToLower(r.Fun)
	if fun == "get" {
		return Get(r.Out, r.FunName)
	}
	return nil
}
func (r *Function) GetCreate() *Field {
	var fun = strings.ToLower(r.Fun)
	if fun == "create" {
		var field reflect.StructField
		var ok bool
		var tp = r.In
		for tp.Kind() != reflect.Struct {
			tp = tp.Elem()
		}
		if field, ok = tp.FieldByName("Data"); !ok {
			return nil
		}
		return Get(field.Type, r.FunName)
	}
	return nil
}
func (r *Function) GetUpdate() *Field {
	var fun = strings.ToLower(r.Fun)
	if fun == "update" {
		var field reflect.StructField
		var ok bool
		var tp = r.In
		for tp.Kind() != reflect.Struct {
			tp = tp.Elem()
		}
		if field, ok = tp.FieldByName("Data"); !ok {
			return nil
		}
		return Get(field.Type, r.FunName)
	}
	return nil
}

type FieldParams struct {
	InType  any
	InName  string
	OutType any
	OutName string
}

func (r *Field) SetFiled(name string, data *Field) *Field {
	var ret = &Field{
		Json:       r.Json,
		Name:       r.Name,
		Type:       r.Type,
		Typescript: r.Typescript,
		Package:    r.Package,
		Enum:       r.Enum,
		Path:       r.Path,
		Import:     r.Import,
		Data:       make([]*Field, 0),
	}
	ret.Type = data.Type
	ret.Name = data.Name
	for _, v := range r.Data {
		if strings.EqualFold(v.Name, name) {
			ret.Data = append(ret.Data, &Field{
				Json:       strings.ToLower(name),
				Name:       name,
				Type:       data.Type + "Data",
				Typescript: data.Typescript,
				Package:    data.Package,
				Enum:       data.Enum,
				Path:       data.Path,
				Data:       data.Data,
				Import:     data.Import,
			})
		} else {
			ret.Data = append(ret.Data, v)
		}
	}
	return ret
}

func (r *FieldParams) InData() *Field {
	if r.InName == "" {
		return nil
	}
	return Get(reflect.TypeOf(r.InType), FunName{})
}
func (r *FieldParams) OutData() *Field {
	if r.OutName == "" {
		return nil
	}
	return Get(reflect.TypeOf(r.OutType), FunName{})
}
func Parse(data []*Function, params *FieldParams) []*Func {
	var (
		inData  = params.InData()
		outData = params.OutData()
	)
	var funs = []*Func{}
	for _, v := range data {
		var in = Get(v.In, v.FunName)
		var out = Get(v.Out, v.FunName)
		if inData != nil {
			in = inData.SetFiled(params.InName, in)
		}
		if outData != nil {
			out = outData.SetFiled(params.OutName, out)
		}
		funs = append(funs, &Func{
			FunName: v.FunName,
			In:      in,
			Out:     out,
			Path:    in.Path,
			Func:    v,
		})
	}
	return funs
}
func ParseObject(data []*Function, namer func(name FunName) string, params *FieldParams, maps map[string]*Object) []*Object {
	// var maps = map[string]*Object{}
	var fs = Parse(data, params)
	for _, v := range fs {
		var name = namer(v.FunName)
		var key = camelToSnake(name)
		if _, ok := maps[key]; !ok {
			maps[key] = &Object{
				Name: key,
			}
			if vv, ok1 := EnumMap[key]; ok1 {
				// log.Println("key111", len(key), key)
				maps[key].Enum = vv
			}
		}
		maps[key].Func = append(maps[key].Func, v)
	}

	for k, v := range EnumMap {
		if _, ok := maps[k]; !ok {
			maps[k] = &Object{
				Name: k,
				Enum: v,
			}
		} else {
		}
	}
	var funName *FunName
	var datas = []*Object{}
	for _, v := range maps {
		if v.FunName.Version != "" {
			funName = &v.FunName
		}
		v.Sync()
		datas = append(datas, v)
	}
	for _, v := range datas {
		if v.FunName.Version == "" && funName != nil {
			v.FunName = *funName
		}
	}
	return datas
}
func GetFunName(fun reflect.Type) FunName {
	for fun.Kind() != reflect.Struct {
		fun = fun.Elem()
	}
	if data, ok := fun.FieldByName("Meta"); ok {
		var path = data.Tag.Get("path")
		arr := gstr.Split(path, "/")
		var ext = []string{}
		for _, v := range arr {
			if v == "" {
				continue
			}
			ext = append(ext, v)
		}
		if len(ext) != 4 {
			return FunName{}
		}
		return FunName{
			Version: ext[0],
			API:     ext[1],
			File:    ext[2],
			Fun:     ext[3],
		}
	} else {
		return FunName{}
	}
}

func mapping(goType string) string {
	switch goType {
	case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8":
		return "number"
	case "float64", "float32":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	default:
		if v, ok := Mapping[goType]; ok {
			return v
		}
		return "any"
	}
}
func mappingImport(goType string) string {
	if v, ok := ImportMap[strings.ToLower(goType)]; ok {
		return v
	}
	return ""
}
func GetModule() string {
	if module != "" {
		return module
	}
	file, err := modfile.Parse("go.mod", gfile.GetBytes("go.mod"), nil)
	if err != nil {
		log.Fatal(err)
	}
	module = file.Module.Mod.Path
	return file.Module.Mod.Path
}
func camelToSnake(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func Ref(t reflect.Type) (reflect.Type, bool) {
	// var t = reflect.TypeOf(data)
	var isArray = false
	var i = 0
	for t.Kind() == reflect.Ptr {
		i += 1
		t = t.Elem()
		for t.Kind() == reflect.Slice {
			t = t.Elem()
			isArray = true
		}
	}
	for t.Kind() == reflect.Slice {
		t = t.Elem()
		isArray = true
	}
	for t.Kind() == reflect.Ptr {
		i += 1
		t = t.Elem()
		for t.Kind() == reflect.Slice {
			t = t.Elem()
			isArray = true
		}
	}
	return t, isArray
}

func inspectStruct(t reflect.Type, name FunName) *Field {
	if t.Kind() != reflect.Struct {
		t = t.Elem()
	}

	field := &Field{
		Type:    t.String(),
		Package: t.PkgPath(),
	}
	field.Typescript = mapping(t.String())
	field.Decimal = IsDecimal(t.String())
	field.Import = mappingImport(field.Typescript)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			if !ft.IsExported() {
				continue
			}
			fieldName := ft.Name
			if ft.Anonymous {
				// 如果该字段是一个匿名字段，则将其展开
				field.Data = append(field.Data, inspectStruct(ft.Type, name).Data...)
				if ft.Name == "Meta" {
					field.Path = ft.Tag.Get("path")
				}
				continue
			}

			var arrayFlag bool
			ft.Type, arrayFlag = Ref(ft.Type)
			if t.Name() == "PageRes" && ft.Name == "Data" {
				arrayFlag = true
				if tt, ok := TypeMapping[name.File]; ok {
					ft.Type, _ = Ref(reflect.TypeOf(tt))
				}
			}

			childField := &Field{
				Name:    fieldName,
				Type:    ft.Type.String(),
				Json:    ft.Tag.Get("json"),
				Enum:    ft.Tag.Get("enum"),
				Desc:    ft.Tag.Get("dc"),
				Table:   ft.Tag.Get("table"),
				Package: ft.Type.PkgPath(),
				Array:   arrayFlag,
			}
			var dcArr = strings.Split(childField.Desc, ":")
			if len(dcArr) > 1 {
				childField.Desc = dcArr[0]
			}

			childField.Typescript = mapping(childField.Type)
			childField.Decimal = IsDecimal(childField.Type)

			childField.Import = mappingImport(childField.Typescript)
			if strings.Contains(strings.ToLower(childField.Json), "omitempty") {
				childField.Optional = true
			}
			var aar2 = strings.Split(childField.Json, ",")
			if len(aar2) > 1 {
				childField.Json = aar2[0]
			}

			if childField.Json == "-" {
				continue
			}
			if childField.Json == "" {
				childField.Json = camelToSnake(fieldName)
			}
			if strings.Contains(childField.Package, module) || childField.Package == "main" {
				// ft.Type, arrayFlag = Ref(ft.Type)
				// childField.Array = arrayFlag

				if ft.Type.Kind() == reflect.Struct {
					childField.Data = inspectStruct(ft.Type, name).Data
				}
			}
			field.Data = append(field.Data, childField)
		}
	}
	return field
}
func Get(data reflect.Type, name FunName) *Field {
	_ = GetModule()
	ff := inspectStruct(data, name)
	// var arr = strings.Split(ff.Type, ".")
	// ff.Name = arr[len(arr)-1]
	return ff
}
