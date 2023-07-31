package genx

type API struct {
	Name string     //对应Hello目录
	Data []*Version //对应Hello目录下的版本
}
type Version struct {
	Name string  //对应v1目录
	Data []*File //对应v1目录下的API文件，例如:user.go hello.go
}
type File struct {
	Name string      //对应user.go hello.go
	Data []*Function //对应user.go hello.go下的API
}
type Function struct {
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
}
type Execute struct {
	Code string
	Data any
	File string
	Must bool
}
type Fun struct {
	API      *API
	Version  *Version
	File     *File
	Function *Function
}
type ApiInput struct {
	Data       []*API
	Controller string // 控制器输出路径
	API        string // API输出路径
	Init       string // 初始化输出路径
}

type LogicInput struct {
	LogicPath   string // 逻辑层路径
	ServicePath string // 服务层路径
}
