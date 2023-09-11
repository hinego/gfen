package horm

const MappingTemplate = `package {{.Name}}

import (
	"{{.Module}}/internal/{{$.Dao}}"
)
var ( {{range .Table}}
	{{ .Name | title}} = {{$.Dao}}.Query{{ .Name | title }} {{end}}
)

`
const MigrateTemplate = `package migrate

import (
	"gorm.io/gorm"
	{{ range .Imports }}
	"{{ . }}"
	{{ end }}
)

func Migrate(db *gorm.DB) (err error) {
	return db.AutoMigrate({{range .Table}}&{{ .Name }}{},{{end}})
}

{{range .Table}}
type {{ .Name }} struct { {{ range .Column }} {{if .Title}}
	{{ .Title }} {{ .Type }} {{ .Tag }} {{end}} {{if .IsBelongsTo}} {{if not .Relation.Fake}}
	{{ .Relation.Name | ToName }}Ref *{{ .Relation.RefTable }} {{ .Relation.Tag }} {{end}} {{end}}
{{- end }}
}
{{end}}
`
const ModelTemplate = `package {{.Package}}

import (
	"context"
	"gorm.io/gen"
	{{ range .Imports }}
	"{{ . }}"
	{{ end }}
)
type edge{{ .Name | title}} struct { {{ range .Relation }} 
	{{ .Name | ToName }} {{if .Array}}[]{{end}}*{{ .RefTable | title }} 
{{- end }}
}

type {{ .Name | title}} struct { {{ range .Column }} {{if .Title}}
	{{ .Title }} {{ .ModelType }} {{ .Tag }} {{end}}
{{- end }}
	edges edge{{ $.Name | title}}
	query *Query
}

func (r *{{ $.Name | title}}) getCacheKeys() []string {
	var keys = make([]string, 0) {{range  $k, $v := .CacheKey}}
	keys = append(keys,gen.CacheKey({{range $k1, $v1 :=  $v.Column}} {{if $k1}},{{end}} {{if $v1.Title}}r.{{ $v1.Title }}{{end}}{{end}})) {{end}}
	return keys
}

func (r *{{ $.Name | title}}) deleted() bool {
	return r.DeletedAt != 0
}

func (r *{{ $.Name | title}}) Tx(tx *Query) *{{ $.Name | title}} {
	return &{{ $.Name | title}}{ {{ range .Column }} {{if .Title}}
		{{ .Title }}: r.{{ .Title }}, {{end}} {{end}}
		edges: r.edges,
		query: tx,
	}
}
func (r *{{ $.Name | title}}) getQuery() *Query  {
	if r.query != nil {
		return r.query
	}
	var ctx = context.Background()
	return Ctx(ctx)
}
func (r *{{ $.Name | title}}) Query() I{{ $.Name | title}}Do {
	return r.getQuery().{{ .Name | title}}().ID(r.{{ .Primary }})
}
func (r *{{ $.Name | title}}) Delete() (info gen.ResultInfo,err error) {
	return r.Query().Delete()
}
{{ range .Relation }}
func (r *{{ $.Name | title}}) Query{{ .Name | ToName }}() I{{ .RefTable | title }}Do { {{if .ForeignPoint}}
	if r.{{ .ForeignName }} == nil {
		return Query{{ .RefTable | title }}().Where({{ .RefTable }}s.{{ .ReferenceName }}.Eq(-9999))
	}  {{end}}
	return Query{{ .RefTable | title }}().Where({{ .RefTable }}s.{{ .ReferenceName }}.Eq({{if .ForeignPoint}}*{{end}}r.{{ .ForeignName }}))
}

func (r *{{ $.Name | title}}) Load{{ .Name | ToName }}(fun ...func(do I{{ .RefTable | title }}Do) I{{ .RefTable | title }}Do) (err error) {
	var do = r.Query{{ .Name | ToName }}()
	for _, v := range fun {
		do = v(do)
	}
	if r.edges.{{ .Name | ToName }}, err = do.{{if .Array}}Find{{else}}First{{end}}(); err != nil {
		return 
	} else {
		return 
	}
}

func (r *{{ $.Name | title}}) Get{{ .Name | ToName }}(update ...bool) (data {{if .Array}}[]{{end}}*{{ .RefTable | title }},err error) {
	if len(update) == 0 && r.edges.{{ .Name | ToName }} != nil {
		return r.edges.{{ .Name | title }},nil
	}
	if r.edges.{{ .Name | ToName }}, err = r.Query{{ .Name | ToName }}().{{if .Array}}Find{{else}}First{{end}}(); err != nil {
		return nil, err
	} else {
		return r.edges.{{ .Name | ToName }}, nil
	}
}

func (r *{{ $.Name | title}}) Get{{ .Name | ToName }}X(update ...bool) ({{if .Array}}[]{{end}}*{{ .RefTable | title }}) {
	if data, err := r.Get{{ .Name | ToName }}(update...); err != nil {
		panic(err)
	} else {
		return data
	}
}
{{end}}
`
const FieldTemplate = `package {{.Table}}s

import (
	"errors"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

var ( {{range .Column}}
 {{.Title}} = field.New{{.Type.Name}}("{{$.TableName}}", "{{.Name}}") {{end}}
)

{{range $k1,$v1 := .Column}} {{if .Enums}}
const ( {{range $k, $v := $v1.Enums.Enums}}
	{{$v1.Title}}{{$v.Name | ToName}} {{$v1.Enums.Type}} = {{$v.String}}	{{end}}
)
{{end}}{{end}}

type ( {{range  $k, $v := .CacheKey}}
		cacheWhere{{$v.Name | title}} struct { {{range $v.Column}} {{if .Title}}
			{{ .Title }} {{ .Type }} {{ .Tag }} {{end}} {{end}}
		} {{end}}
)

{{range  $k, $v := .CacheKey}}

func Cache{{$v.Name | title}} ({{range $k1, $v1 :=  $v.Column}} {{if $k1}},{{end}} {{if $v1.Title}}{{ $v1.Title | lower }} {{ .Type }}{{end}}{{end}}) gen.Condition {
	return &cacheWhere{{$v.Name | title}}{
		{{range $k1, $v1 :=  $v.Column}} {{if $v1.Title}}{{ $v1.Title }}:{{ $v1.Title  | lower  }},{{end}}{{end}}
	}
}

func (c *cacheWhere{{$v.Name | title}})  BeCond() any {
	return &gen.CacheWhere{
		Code: gen.CacheKey({{range $k1, $v1 :=  $v.Column}} {{if $k1}},{{end}} {{if $v1.Title}}c.{{ $v1.Title }}{{end}}{{end}}),
		Where: []gen.Condition{ {{range $v.Column}} {{if .Title}}
		{{ .Title }}.Eq(c.{{ .Title }}), {{end}} {{end}}
		},
	}
}
func (c *cacheWhere{{$v.Name | title}})  CondError() error {
	return errors.New("cache key")
}
{{end}}

`
const DaoTemplate = `package {{.Package}}

import (
	"context" {{ range .Imports }}
	"{{ . }}" {{ end }}
	"github.com/gogf/gf/util/gconv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func new{{.Model}}(db *gorm.DB, opts ...gen.DOOption) I{{.Model}}Do {
	do := gen.NewDo(db, &{{.Model}}{})
	_{{.Table}}Do := {{.Table}}Do{}
	_{{.Table}}Do.Dao = &do
	_{{.Table}}Do.table = _{{.Table}}Do.TableName()
	_{{.Table}}Do.fieldMap = map[string]any{ {{range .Column}}
		"{{.Name}}":field.New{{.Type.Name}}(_{{$.Table}}Do.table, "{{.Name}}"), {{end}}
	} 
	return &_{{.Table}}Do
}


type I{{.Model}}Do interface { 
	gen.SubQuery
	Debug() I{{.Model}}Do
	WithContext(ctx context.Context) I{{.Model}}Do
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() I{{.Model}}Do
	WriteDB() I{{.Model}}Do
	As(alias string) gen.Dao
	Session(config *gorm.Session) I{{.Model}}Do
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) I{{.Model}}Do
	Not(conds ...gen.Condition) I{{.Model}}Do
	Or(conds ...gen.Condition) I{{.Model}}Do
	Select(conds ...field.Expr) I{{.Model}}Do
	Where(conds ...gen.Condition) I{{.Model}}Do

	ID({{.Primary}} {{.PrimaryType}}) I{{.Model}}Do
	Get({{.Primary}} {{.PrimaryType}}) (*{{.Model}}, error)
	MustGet({{.Primary}} {{.PrimaryType}}) *{{.Model}}
	MustDelete({{.Primary}} {{.PrimaryType}}) (err error)
	Order(conds ...field.Expr) I{{.Model}}Do
	Distinct(cols ...field.Expr) I{{.Model}}Do
	Omit(cols ...field.Expr) I{{.Model}}Do
	Join(table schema.Tabler, on ...field.Expr) I{{.Model}}Do
	LeftJoin(table schema.Tabler, on ...field.Expr) I{{.Model}}Do
	RightJoin(table schema.Tabler, on ...field.Expr) I{{.Model}}Do
	Group(cols ...field.Expr) I{{.Model}}Do
	Having(conds ...gen.Condition) I{{.Model}}Do
	Limit(limit int) I{{.Model}}Do
	Offset(offset int) I{{.Model}}Do
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) I{{.Model}}Do
	Unscoped() I{{.Model}}Do
	Create(values ...*{{.Model}}) error
	CreateAny(values ...any) error
	CreateInBatches(values []*{{.Model}}, batchSize int) error
	Save(values ...*{{.Model}}) error
	First() (*{{.Model}}, error)
	Take() (*{{.Model}}, error)
	Last() (*{{.Model}}, error)
	Find() ([]*{{.Model}}, error)
	DoCache(date int64,store func(data cacheKey)) (err error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*{{.Model}}, err error)
	FindInBatches(result *[]*{{.Model}}, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*{{.Model}}) (info gen.ResultInfo, err error)
	Update(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) I{{.Model}}Do
	Assign(attrs ...field.AssignExpr) I{{.Model}}Do
	Joins(fields ...field.RelationField) I{{.Model}}Do
	Preload(fields ...field.RelationField) I{{.Model}}Do
	FirstOrInit() (*{{.Model}}, error)
	FirstOrCreate() (*{{.Model}}, error)
	FindByPage(offset int, limit int) (result []*{{.Model}}, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) I{{.Model}}Do
	UnderlyingDB() *gorm.DB
	schema.Tabler {{range .Relation}}
	With{{.Name | title}}(fun ...func(do I{{ .RefTable | title }}Do) I{{ .RefTable | title }}Do) I{{$.Model}}Do {{end}}
	GetFieldByName(fieldName string) (field.OrderExpr, bool)
	GetField(fieldName string) (any, bool)
	Query(face queryFace) (err error)
	WhereRaw(data any, args ...any) I{{$.Model}}Do
}

type {{.Table}}Preload struct{  {{range .Relation}}
	{{.Name | title }} bool 
	{{.Name  | title }}Expr []func(do I{{ .RefTable | title }}Do) I{{ .RefTable | title }}Do{{end}}
}
type {{.Table}}Do struct{ 
	gen.Dao 
	preload {{.Table}}Preload
	table string
	fieldMap map[string]any
	cacheWhere *gen.CacheWhere
}
func (a {{.Table}}Do) Query(face queryFace) (err error) {
	data := &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao: a.Dao}
	face.SetFun(a.GetFieldByName)
	var result any
	var total int64
	if result, total, err = data.Order(face.GetOrder()...).WhereRaw(face.GetWhere()).FindByPage(face.GetOffset(), face.GetSize()); err != nil {
		return err
	}
	face.SetData(result)
	face.SetTotal(total)
	return
}

func (a {{.Table}}Do) WhereRaw(data any, args ...any) I{{$.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao: a.Dao.WhereRaw(data,args...)}
}
func (a {{.Table}}Do) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a {{.Table}}Do) GetField(fieldName string) (any, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	return _f, ok
}

{{range .Relation}}
func (a {{$.Table}}Do) With{{.Name | title}}(fun ...func(do I{{ .RefTable | title }}Do) I{{ .RefTable | title }}Do) I{{$.Model}}Do {
	a.preload.{{.Name | title}} = true
	a.preload.{{.Name | title}}Expr = append(a.preload.{{.Name | title}}Expr, fun...)
	return &{{$.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Debug()}
}
{{end}}

func (a {{.Table}}Do) doPreload(data ...*{{.Table | title}}) (err error) {
	{{if .Relation}}
	for _,v :=range data	{	{{range .Relation}}
		if a.preload.{{.Name | title}} {
			if err = v.Load{{.Name | title}}(a.preload.{{.Name | title}}Expr...); err != nil {
				return
			}
		} {{end}}
	} {{end}}
	return nil
}

func (a {{.Table}}Do) Debug() I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Debug()}
}

func (a {{.Table}}Do) WithContext(ctx context.Context) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.WithContext(ctx)}
}

func (a {{.Table}}Do) ReadDB() I{{.Model}}Do {
	return a.Clauses(dbresolver.Read)
}

func (a {{.Table}}Do) WriteDB() I{{.Model}}Do {
	return a.Clauses(dbresolver.Write)
}

func (a {{.Table}}Do) Session(config *gorm.Session) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Session(config)}
}

func (a {{.Table}}Do) Clauses(conds ...clause.Expression) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Clauses(conds...)}
}

func (a {{.Table}}Do) Returning(value interface{}, columns ...string) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Returning(value, columns...)}
}

func (a {{.Table}}Do) Not(conds ...gen.Condition) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Not(conds...)}
}

func (a {{.Table}}Do) Or(conds ...gen.Condition) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Or(conds...)}
}

func (a {{.Table}}Do) Select(conds ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Select(conds...)}
}

func (a {{.Table}}Do) Where(conds ...gen.Condition) I{{.Model}}Do {
	var conditions = make([]gen.Condition, 0)
	var cacheWhere *gen.CacheWhere
	var num int
	for _, v := range conds {
		var err error
		if err = v.CondError();err != nil {
			var data any
			if data = v.BeCond(); data != nil {
				var dataWhere *gen.CacheWhere
				var ok bool
				if dataWhere, ok = data.(*gen.CacheWhere); ok {
					cacheWhere = dataWhere
					conditions = append(conditions, dataWhere.Where...)
				}
			}
		}else {
			num += 1
			conditions = append(conditions, v)
		}
	}
	if cacheWhere != nil && num == 0 {
		a.cacheWhere = cacheWhere
	}else {
		a.cacheWhere = nil
	}
	if len(conditions) == 0 {
		return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao}
	}
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Where(conditions...)}
}

func (a {{.Table}}Do) ID({{.Primary}} {{.PrimaryType}}) I{{.Model}}Do {
	return a.Where({{.Table}}s.Cache{{.Primary}}({{.Primary}}))
}
func (a {{.Table}}Do) Get({{.Primary}} {{.PrimaryType}}) (*{{.Model}}, error) {
	return a.ID({{.Primary}}).First()
}

func (a {{.Table}}Do) MustGet({{.Primary}} {{.PrimaryType}}) *{{.Model}} {
	data, _ := a.ID({{.Primary}}).First()
	return data
}

func (a {{.Table}}Do) MustDelete({{.Primary}} {{.PrimaryType}}) (err error) {
	_, err = a.ID({{.Primary}}).Delete()
	return
}

func (a {{.Table}}Do) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) I{{.Model}}Do {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a {{.Table}}Do) Order(conds ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Order(conds...)}
}

func (a {{.Table}}Do) Distinct(cols ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Distinct(cols...)}
}

func (a {{.Table}}Do) Omit(cols ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Omit(cols...)}
}

func (a {{.Table}}Do) Join(table schema.Tabler, on ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Join(table, on...)}
}

func (a {{.Table}}Do) LeftJoin(table schema.Tabler, on ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.LeftJoin(table, on...)}
}

func (a {{.Table}}Do) RightJoin(table schema.Tabler, on ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.RightJoin(table, on...)}
}

func (a {{.Table}}Do) Group(cols ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Group(cols...)}
}

func (a {{.Table}}Do) Having(conds ...gen.Condition) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Having(conds...)}
}

func (a {{.Table}}Do) Limit(limit int) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Limit(limit)}
}

func (a {{.Table}}Do) Offset(offset int) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Offset(offset)}
}

func (a {{.Table}}Do) Scopes(funcs ...func(gen.Dao) gen.Dao) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Scopes(funcs...)}
}

func (a {{.Table}}Do) Unscoped() I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Unscoped()}
}

func (a {{.Table}}Do) Create(values ...*{{.Model}}) error {
	if len(values) == 0 {
		return nil
	}
	return a.Dao.Create(values)
}

func (a {{.Table}}Do) CreateAny(values ...any) error {
	if len(values) == 0 {
		return nil
	}
	var data = make([]*{{.Model}}, 0)
	if err := gconv.Scan(values, data); err != nil {
		return err
	}
	return a.Dao.Create(data)
}

func (a {{.Table}}Do) CreateInBatches(values []*{{.Model}}, batchSize int) error {
	return a.Dao.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a {{.Table}}Do) Save(values ...*{{.Model}}) error {
	if len(values) == 0 {
		return nil
	}
	return a.Dao.Save(values)
}
func (a {{.Table}}Do) getByCache() (result *{{.Model}}, vaild bool) {
	if a.cacheWhere == nil || !{{.Table | lower}}Cacher.UseCache() {
		return nil, false
	}
	var v any
	var ok bool
	if v,ok = {{.Table | lower}}Cacher.Get(a.cacheWhere.Code); !ok {
		return nil, true
	}
	if result, ok = v.(*{{.Model}}); !ok {
		return nil, true
	}
	return result, true
}
func (a {{.Table}}Do) First() (*{{.Model}}, error) {
	if result, ok := a.getByCache(); ok {
		return result, nil
	}

	if result, err := a.Dao.First(); err != nil {
		return nil, err
	} else {
		var data = result.(*{{.Model}})
		if err = a.doPreload(data); err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (a {{.Table}}Do) Take() (*{{.Model}}, error) {

	if result, ok := a.getByCache(); ok {
		return result, nil
	}

	if result, err := a.Dao.Take(); err != nil {
		return nil, err
	} else {
		var data = result.(*{{.Model}})
		if err = a.doPreload(data); err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (a {{.Table}}Do) Last() (*{{.Model}}, error) {

	if result, ok := a.getByCache(); ok {
		return result, nil
	}

	if result, err := a.Dao.Last(); err != nil {
		return nil, err
	} else {
		var data = result.(*{{.Model}})
		if err = a.doPreload(data); err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (a {{.Table}}Do) Find() ([]*{{.Model}}, error) {
	if result, err := a.Dao.Find();err !=nil {
		return nil, err
	}else {
		var data = result.([]*{{.Model}})
		if err = a.doPreload(data...); err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (a {{.Table}}Do) DoCache(date int64,store func(data cacheKey)) (err error) {
	var results = make([]*{{.Model}}, 0)
	if err = a.Where({{.Table}}s.UpdatedAt.Gt(date)).FindInBatches(&results, 500, func(tx gen.Dao, batch int) (err error) {
		for _, v := range results {
			store(v)
		}
		return
	}); err != nil {
		return
	}
	return nil
}

func (a {{.Table}}Do) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*{{.Model}}, err error) {
	buf := make([]*{{.Model}}, 0, batchSize)
	err = a.Dao.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a {{.Table}}Do) FindInBatches(result *[]*{{.Model}}, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.Dao.FindInBatches(result, batchSize, fc)
}

func (a {{.Table}}Do) Attrs(attrs ...field.AssignExpr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Attrs(attrs...)}
}

func (a {{.Table}}Do) Assign(attrs ...field.AssignExpr) I{{.Model}}Do {
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:a.Dao.Assign(attrs...)}
}

func (a {{.Table}}Do) Joins(fields ...field.RelationField) I{{.Model}}Do {
	var data = a.Dao
	for _, _f := range fields {
		data.Joins(_f)
	}
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:data}	
}

func (a {{.Table}}Do) Preload(fields ...field.RelationField) I{{.Model}}Do {
	var data = a.Dao
	for _, _f := range fields {
		data.Preload(_f)
	}
	return &{{.Table}}Do{table: a.table,cacheWhere:a.cacheWhere, fieldMap: a.fieldMap, preload: a.preload, Dao:data}		
}

func (a {{.Table}}Do) FirstOrInit() (*{{.Model}}, error) {
	if result, err := a.Dao.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*{{.Model}}), nil
	}
}

func (a {{.Table}}Do) FirstOrCreate() (*{{.Model}}, error) {
	if result, err := a.Dao.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*{{.Model}}), nil
	}
}

func (a {{.Table}}Do) FindByPage(offset int, limit int) (result []*{{.Model}}, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a {{.Table}}Do) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a {{.Table}}Do) Scan(result interface{}) (err error) {
	return a.Dao.Scan(result)
}

func (a {{.Table}}Do) Delete(models ...*{{.Model}}) (result gen.ResultInfo, err error) {
	return a.Dao.Delete(models)
}

`
const GenTemplate = `package {{.Package}}

import (
	"time"
	"log"
	"strings"
	"errors"
	"sync"
	"context"
	"database/sql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"gorm.io/gen/field"

)

var (
	_db          *gorm.DB
) 

type CacheWhere interface {
	GetCondition() []gen.Condition
	GetCacheKey() string
}

type queryFace interface {
	GetOffset() int
	GetSize() int
	GetOrder() (data []field.Expr)
	GetWhere() map[string]any
	SetFun(fun func(fieldName string) (field.OrderExpr, bool))
	SetData(data any)
	SetTotal(total int64)
}



type cacheKey interface {
	getCacheKeys() []string
	deleted() bool
}
type cacheValue struct {
	Valid bool
	Value any
}
type cacheDriver struct {
	level              int
	err                error
	cache              sync.Map
	lastExecuteTime    int64
	nominalExecuteTime int64
	fun                func(date int64, store func(data cacheKey)) (err error)
	lock               sync.Mutex
}

func (r *cacheDriver) Init() (err error) {
	switch r.level {
	case 1:
		go r.Reload()
	case 2:
		return r.Reload()
	}
	return nil
}
func (r *cacheDriver) UseCache() bool {
	if r.level == 0 {
		return false
	}
	return r.err == nil
}
func (r *cacheDriver) get(key string) (value *cacheValue, ok bool) {
	var v any
	if v, ok = r.cache.Load(key); ok {
		value = v.(*cacheValue)
	}
	return
}
func (r *cacheDriver) Get(key string) (value any, ok bool) {
	var v *cacheValue
	if v, ok = r.get(key); ok {
		if v.Valid {
			return v.Value, true
		} else {
			r.cache.Delete(key)
			return nil, false
		}
	}
	return
}
func (r *cacheDriver) Put(data cacheKey) {
	for _, v := range data.getCacheKeys() {
		r.store(v, data)
	}
}
func (r *cacheDriver) store(key string, value cacheKey) {
	var (
		v  *cacheValue
		ok bool
	)
	if v, ok = r.get(key); ok {
		v.Valid = false
		r.cache.Delete(key)
	}
	if value.deleted() { // 如果数据已经被删除了 那么就不要缓存了
		return
	}
	r.cache.Store(key, &cacheValue{
		Valid: true,
		Value: value,
	})
}
func (r *cacheDriver) reload() (check bool, err error) {
	if r.fun == nil {
		return
	}
	if !r.lock.TryLock() {
		r.nominalExecuteTime = time.Now().Unix()
		return
	}
	defer r.lock.Unlock()
	var now = time.Now().Unix()
	var date = now
	if r.lastExecuteTime == 0 {
		date = 0
	}
	if err = r.fun(date, r.Put); err != nil {
		r.err = err // 如果出错 那么则暂时不使用缓存了
		return
	}
	if r.err != nil {
		r.err = nil
	}
	r.lastExecuteTime = now
	return true, nil
}
func (r *cacheDriver) Reload() (err error) {
	for {
		var check bool
		if check, err = r.reload(); err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if check && r.nominalExecuteTime > r.lastExecuteTime {
			return r.Reload()
		}
		break
	}
	return
}


var ( {{range .Table}} 
	{{ .Name | lower }}Cacher = &cacheDriver{
		err: errors.New("wait load"),
		level: {{.CacheLevel}},

	} {{end}}
)




func CallbackPlugin() gorm.Plugin {
	return &sCallback{}
}

type sCallback struct{}

func (r *sCallback) Name() string {
	return ""
}
func (r *sCallback) Initialize(db *gorm.DB) (err error) {
	err = db.Callback().Create().After("gorm:create").Register("gorm:cache:after_create", r.after)
	if err != nil {
		return err
	}

	err = db.Callback().Delete().After("gorm:delete").Register("gorm:cache:after_delete", r.after)
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("gorm:update").Register("gorm:cache:after_update", r.after)
	if err != nil {
		return err
	}
	return
}
func (r *sCallback) after(db *gorm.DB) {
	if db.RowsAffected == 0 {
		return // no rows affected, no need to invalidate cache
	}
	var tableName string
	if db.Statement.Schema != nil {
		tableName = strings.ToLower(db.Statement.Schema.Table)
	} else {
		tableName = strings.ToLower(db.Statement.Table)
	}
	var fun func() (err error)
	switch tableName { {{range .Table}} {{if .CacheLevel}}
	case "{{.TableName}}":
		fun = {{ .Name | lower }}Cacher.Reload {{end}} {{end}}
	default:
		fun = func() (err error) {
			return
		}
	}
	if err := fun(); err != nil {
		log.Println("刷新缓存失败", err)
	}
}









{{range .Table}}
func Query{{ .Name | title}}() I{{ .Name | title}}Do {
	return new{{ .Name | title}}(_db)
}
{{end}}
func RegisterDB(db *gorm.DB) {
	_db = db.Session(&gorm.Session{NewDB: true})
	initCache()
}
func initCache() { {{range .Table}} {{if .CacheLevel}}
	{{ .Name | lower }}Cacher.fun = Query{{ .Name | title}}().DoCache {{end}} {{end}}

	var cacheFunc = []func() error{ {{range .Table}} {{if .CacheLevel}}
	{{ .Name | lower }}Cacher.Init,{{end}} {{end}}
	}
	for _, v := range cacheFunc {
		if err := v(); err != nil {
			panic(err)
		}
	}
}
func DB() *gorm.DB {
	return _db.Session(&gorm.Session{NewDB: true})
}
func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	_db = db
}


func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db: db,
	}
}

type Query struct{
	db *gorm.DB 
	ctx context.Context
}
func Ctx(ctx context.Context) *Query {
	return &Query{
		db: _db,
		ctx: ctx,
	}
}
{{range .Table}}
func (q *Query) {{ .Name | title}}() I{{ .Name | title}}Do {
	if q.ctx != nil {
		return new{{ .Name | title}}(q.db).WithContext(q.ctx)
	}
	return new{{ .Name | title}}(q.db)
}
{{end}}
func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db: db, 
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db: db, 
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}


`
