package horm

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
	{{ .Title }} {{ .Type }} {{ .Tag }} {{end}} {{if .IsBelongsTo}}
	{{ .Relation.Name | title }} *{{ .Relation.Table }} {{ .Relation.Tag }} {{end}}
{{- end }}
}
{{end}}
`
const ModelTemplate = `package {{.Package}}

import (
	{{ range .Imports }}
	"{{ . }}"
	{{ end }}
)
type edge{{ .Name | title}} struct { {{ range .Relation }} 
	{{ .Name | title }} {{if .Array}}[]{{end}}*{{ .RefTable | title }}
{{- end }}
}

type {{ .Name | title}} struct { {{ range .Column }} {{if .Title}}
	{{ .Title }} {{ .ModelType }} {{ .Tag }} {{end}}
{{- end }}
	edges edge{{ $.Name | title}}
}

{{ range .Relation }}
func (r *{{ $.Name | title}}) Query{{ .Name | title }}() I{{ .RefTable | title }}Do {
	return Query{{ .RefTable | title }}().Where({{ .RefTable }}s.{{ .ReferenceName }}.Eq(r.{{ .ForeignName }}))
}

func (r *{{ $.Name | title}}) Get{{ .Name | title }}(update ...bool) (data {{if .Array}}[]{{end}}*{{ .RefTable | title }},err error) {
	if len(update) == 0 && r.edges.{{ .Name | title }} != nil {
		return r.edges.{{ .Name | title }},nil
	}
	if r.edges.{{ .Name | title }}, err = r.Query{{ .Name | title }}().{{if .Array}}Find{{else}}First{{end}}(); err != nil {
		return nil, err
	} else {
		return r.edges.{{ .Name | title }}, nil
	}
}

func (r *{{ $.Name | title}}) Get{{ .Name | title }}X(update ...bool) ({{if .Array}}[]{{end}}*{{ .RefTable | title }}) {
	if data, err := r.Get{{ .Name | title }}(update...); err != nil {
		panic(err)
	} else {
		return data
	}
}
{{end}}
`
const FieldTemplate = `package {{.Table}}s

import "gorm.io/gen/field"


var ( {{range .Column}}
 {{.Title}} = field.New{{.Type.Name}}("{{$.Table}}", "{{.Name}}") {{end}}
)

func init(){
	
}
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
	_{{.Table}} := {{.Table}}{}
	_{{.Table}}.{{.Table}}Do.Dao = &do
	_{{.Table}}.table = _{{.Table}}.{{.Table}}Do.TableName()
	_{{.Table}}.fillFieldMap()
	return &_{{.Table}}
}

type {{.Table}} struct {
	{{.Table}}Do
	table string
	fieldMap map[string]any
}


func (a *{{.Table}}) Table(newTableName string) *{{.Table}} {
	a.{{.Table}}Do.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a *{{.Table}}) updateTableName(table string) *{{.Table}} {
	a.table = table
	return a
}

func (a *{{.Table}}) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *{{.Table}}) GetField(fieldName string) (any, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	return _f, ok
}

func (a *{{.Table}}) fillFieldMap() {
	a.fieldMap = make(map[string]any) {{range .Column}}
	a.fieldMap["{{.Name}}"] = field.New{{.Type.Name}}(a.table, "{{.Name}}") {{end}}
}

func (a *{{.Table}}) clone(db *gorm.DB) *{{.Table}} {
	a.{{.Table}}Do.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a *{{.Table}}) replaceDB(db *gorm.DB) *{{.Table}} {
	a.{{.Table}}Do.ReplaceDB(db)
	return a
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

	Key({{.Primary}} {{.PrimaryType}}) I{{.Model}}Do
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
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*{{.Model}}, err error)
	FindInBatches(result *[]*{{.Model}}, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*{{.Model}}) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
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
	With{{.Name | title}}() I{{$.Model}}Do {{end}}
	
}

type {{.Table}}Preload struct{  {{range .Relation}}
	{{.Name}} bool {{end}}
}
type {{.Table}}Do struct{ 
	gen.Dao 
	preload {{.Table}}Preload
}

{{range .Relation}}
func (a {{$.Table}}Do) With{{.Name | title}}() I{{$.Model}}Do {
	a.preload.{{.Name}} = true
	return &{{$.Table}}Do{preload: a.preload, Dao:a.Dao.Debug()}
}
{{end}}

func (a {{.Table}}Do) doPreload(data ...*{{.Table | title}}) (err error) {
	for _,v :=range data	{	{{range .Relation}}
		if a.preload.{{.Name}} {
			if _,err = v.Get{{.Name | title}}(); err != nil {
				return
			}
		} {{end}}
	}
	return nil
}

func (a {{.Table}}Do) Debug() I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Debug()}
}

func (a {{.Table}}Do) WithContext(ctx context.Context) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.WithContext(ctx)}
}

func (a {{.Table}}Do) ReadDB() I{{.Model}}Do {
	return a.Clauses(dbresolver.Read)
}

func (a {{.Table}}Do) WriteDB() I{{.Model}}Do {
	return a.Clauses(dbresolver.Write)
}

func (a {{.Table}}Do) Session(config *gorm.Session) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Session(config)}
}

func (a {{.Table}}Do) Clauses(conds ...clause.Expression) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Clauses(conds...)}
}

func (a {{.Table}}Do) Returning(value interface{}, columns ...string) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Returning(value, columns...)}
}

func (a {{.Table}}Do) Not(conds ...gen.Condition) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Not(conds...)}
}

func (a {{.Table}}Do) Or(conds ...gen.Condition) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Or(conds...)}
}

func (a {{.Table}}Do) Select(conds ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Select(conds...)}
}

func (a {{.Table}}Do) Where(conds ...gen.Condition) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Where(conds...)}
}

func (a {{.Table}}Do) Key({{.Primary}} {{.PrimaryType}}) I{{.Model}}Do {
	return a.Where({{.Table}}s.{{.Primary}}.Eq({{.Primary}}))
}
func (a {{.Table}}Do) Get({{.Primary}} {{.PrimaryType}}) (*{{.Model}}, error) {
	return a.Key({{.Primary}}).First()
}

func (a {{.Table}}Do) MustGet({{.Primary}} {{.PrimaryType}}) *{{.Model}} {
	data, _ := a.Key({{.Primary}}).First()
	return data
}

func (a {{.Table}}Do) MustDelete({{.Primary}} {{.PrimaryType}}) (err error) {
	_, err = a.Key({{.Primary}}).Delete()
	return
}

func (a {{.Table}}Do) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) I{{.Model}}Do {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a {{.Table}}Do) Order(conds ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Order(conds...)}
}

func (a {{.Table}}Do) Distinct(cols ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Distinct(cols...)}
}

func (a {{.Table}}Do) Omit(cols ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Omit(cols...)}
}

func (a {{.Table}}Do) Join(table schema.Tabler, on ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Join(table, on...)}
}

func (a {{.Table}}Do) LeftJoin(table schema.Tabler, on ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.LeftJoin(table, on...)}
}

func (a {{.Table}}Do) RightJoin(table schema.Tabler, on ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.RightJoin(table, on...)}
}

func (a {{.Table}}Do) Group(cols ...field.Expr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Group(cols...)}
}

func (a {{.Table}}Do) Having(conds ...gen.Condition) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Having(conds...)}
}

func (a {{.Table}}Do) Limit(limit int) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Limit(limit)}
}

func (a {{.Table}}Do) Offset(offset int) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Offset(offset)}
}

func (a {{.Table}}Do) Scopes(funcs ...func(gen.Dao) gen.Dao) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Scopes(funcs...)}
}

func (a {{.Table}}Do) Unscoped() I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Unscoped()}
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

func (a {{.Table}}Do) First() (*{{.Model}}, error) {
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
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Attrs(attrs...)}
}

func (a {{.Table}}Do) Assign(attrs ...field.AssignExpr) I{{.Model}}Do {
	return &{{.Table}}Do{preload: a.preload, Dao:a.Dao.Assign(attrs...)}
}

func (a {{.Table}}Do) Joins(fields ...field.RelationField) I{{.Model}}Do {
	var data = a.Dao
	for _, _f := range fields {
		data.Joins(_f)
	}
	return &{{.Table}}Do{preload: a.preload, Dao:data}	
}

func (a {{.Table}}Do) Preload(fields ...field.RelationField) I{{.Model}}Do {
	var data = a.Dao
	for _, _f := range fields {
		data.Preload(_f)
	}
	return &{{.Table}}Do{preload: a.preload, Dao:data}		
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
	"context"
	"database/sql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	_db          *gorm.DB
) 
{{range .Table}}
func Query{{ .Name | title}}() I{{ .Name | title}}Do {
	return new{{ .Name | title}}(_db)
}
{{end}}
func RegisterDB(db *gorm.DB) {
	_db = db.Session(&gorm.Session{NewDB: true})
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
