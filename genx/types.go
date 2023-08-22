package genx

import (
	"reflect"
)

type (
	Field struct {
		Name string
		Type any
	}
	Code struct {
		Req  []*Field       //请求结构体的参数
		Res  []*Field       //响应结构体的参数
		Uses []any          // 需要使用那些东西
		Code string         // 代码模板
		Data map[string]any // 代码模板的数据
	}
	API struct {
		Name string     //对应Hello目录
		Data []*Version //对应Hello目录下的版本
	}
	Version struct {
		Name string  //对应v1目录
		Data []*File //对应v1目录下的API文件，例如:user.go hello.go
	}
	File struct {
		Name string      //对应user.go hello.go
		Data []*Function //对应user.go hello.go下的API
	}
	Function struct {
		Name        string //对应 方法的名称 例如：HelloReq中的Hello
		Path        string //对应 方法的路径 例如：HelloReq中的/hello
		Tags        string //对应 方法的标签 例如：HelloReq中的hello
		Summary     string //对应 方法的简介 例如：HelloReq中的hello
		Method      string //对应 方法的请求方式 例如：HelloReq中的GET
		Deprecated  string //对应 方法的废弃 例如：HelloReq中的hello
		Description string //对应 方法的详细描述 例如：HelloReq中的hello
		Mime        string //对应 方法的MIME类型 例如：HelloReq中的text/html
		Type        string //对应 方法的类型 例如：HelloReq中的HelloRes
		In          string //对应 方法的提交方式 例如：HelloReq中的header/path/query/cookie
		Default     string //对应 方法的默认值 例如：HelloReq中的string
		Code        *Code  //对应 方法的代码 如果为空则使用默认的代码
		Must        bool   //是否覆盖已有的代码
	}
	Execute struct {
		Code    string
		Data    any
		Replace map[string]string
		Map     map[string]any
		File    string
		Must    bool
		Debug   bool
	}
	Fun struct {
		API      *API
		Version  *Version
		File     *File
		Function *Function
	}
	ApiInput struct {
		Data       []*API
		Controller string // 控制器输出路径
		API        string // API输出路径
		Init       string // 初始化输出路径
		Clear      bool   // 是否情况非自动生成的文件
	}
	LogicInput struct {
		LogicPath   string // 逻辑层路径
		ServicePath string // 服务层路径
	}
)

func (r *Code) WithData(data map[string]any) *Code {
	return &Code{
		Uses: r.Uses,
		Code: r.Code,
		Data: data,
	}
}
func (r *Code) Import() []string {
	var data = make(map[string]bool)
	for _, use := range r.Uses {
		switch e := use.(type) {
		case string:
			data[e] = true
		default:
			ss := getPack(use)
			if ss != "" {
				data[ss] = true
			}
		}
	}
	return toSlice(data)
}
func (r *Code) ApiImport() []string {
	var data = make(map[string]bool)
	for _, use := range append(r.Req, r.Res...) {
		ss := getPack(use.Type)
		if ss != "" {
			data[ss] = true
		}
	}
	return toSlice(data)
}
func (r *Code) Request() []string {
	var data = make(map[string]bool)
	for _, v := range r.Req {
		var name string
		if v.Name != "" {
			name = v.Name + " "
		}
		name += getName(v.Type)
		data[name] = true
	}
	return toSlice(data)
}
func (r *Code) Response() []string {
	var data = make(map[string]bool)
	for _, v := range r.Res {
		var name string
		if v.Name != "" {
			name = v.Name + " "
		}
		name += getName(v.Type)
		data[name] = true
	}
	return toSlice(data)
}
func toSlice(data map[string]bool) (ret []string) {
	for k := range data {
		ret = append(ret, k)
	}
	return ret
}
func getPack(t any) string {
	ref := reflect.ValueOf(t)
	for ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	return ref.Type().PkgPath()
}
func getName(t any) string {
	ref := reflect.ValueOf(t)
	if ref.Kind() == reflect.Ptr {
		return "*" + ref.Elem().Type().String()
	}
	return ref.Type().String()
}
