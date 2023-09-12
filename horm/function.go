package horm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gorm.io/gorm/schema"
)

func (r *Type) NotNull() *Type {
	return &Type{
		Name:       r.Name,
		Postgres:   r.Postgres,
		Mysql:      r.Mysql,
		Sqlite:     r.Sqlite,
		Type:       r.Type,
		value:      nil,
		Native:     r.Native,
		SetNotNull: true,
		Point:      r.Point,
	}
}
func (r *Type) PointAble() *Type {
	return &Type{
		Name:       r.Name,
		Postgres:   r.Postgres,
		Mysql:      r.Mysql,
		Sqlite:     r.Sqlite,
		Type:       r.Type,
		Native:     r.Native,
		SetNotNull: r.SetNotNull,
		Point:      true,
	}
}
func (r *Type) Get() (value any) {
	return r.value
}
func (r *Type) Valid() (ok bool) {
	return r.value != nil
}
func (r *Type) Serializer() (value bool) {
	if r.Package() == "" {
		return false
	}
	if r.Native {
		return false
	}
	valType := reflect.TypeOf(r.Type)
	if valType.Kind() != reflect.Ptr {
		valType = reflect.PtrTo(valType)
	}
	ok1 := valType.Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem())
	ok2 := valType.Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem())
	return !(ok1 && ok2)
}

func (r *Type) String() string {
	ref := reflect.ValueOf(r.Type)
	if ref.Kind() == reflect.Ptr {
		return "*" + ref.Elem().Type().String()
	}
	if r.Point {
		return "*" + ref.Type().String()
	}
	return ref.Type().String()
}
func (r *Type) Model() string {
	ref := reflect.ValueOf(r.Type)
	for ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	return ref.Type().String()
}
func (r *Type) Package() string {
	ref := reflect.ValueOf(r.Type)
	for ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	return ref.Type().PkgPath()
}
func (c *Relation) Tag() string {
	var tags []string
	tags = append(tags, fmt.Sprintf("foreignKey:%s", ToName(c.Foreign)))
	tags = append(tags, fmt.Sprintf("references:%s", ToName(c.Reference)))
	if len(tags) == 0 {
		return ""
	}
	return "`gorm:\"" + strings.Join(tags, ";") + "\"`"
}
func (c *Relation) ForeignName() string {
	return ToName(c.Foreign)
}
func (c *Relation) TableName() string {
	return ToName(c.Table)
}
func (c *Relation) ReferenceName() string {
	return ToName(c.Reference)
}
func (c *Relation) IsBelongsTo() bool {
	return c.Type == BelongsTo
}
func (c *Relation) Array() bool {
	return c.Type != BelongsTo && !c.Unique
}
func (c *Relation) Reverse() *Relation {
	if c.Type != BelongsTo {
		return nil
	}
	var typ = HasOne
	if c.Array() {
		typ = HasMany
	}
	if c.RefName == "" {
		c.RefName = ToName(c.Table)
		if !c.Array() {
			c.RefName = c.RefName + "s"
		}
	}
	return &Relation{
		Name:      c.RefName,
		RefName:   c.Name,
		Type:      typ,
		Table:     c.RefTable,
		RefTable:  c.Table,
		Query:     c.Query,
		Foreign:   c.Reference,
		Reference: c.Foreign,
		Unique:    c.Unique,
	}
}
func (c *Column) IsBelongsTo() bool {
	return c.Relation != nil && c.Relation.Type == BelongsTo
}
func (c *Column) Title() string {
	return ToName(c.Name)
}
func (c *Column) Tag() string {
	var tags []string

	// Add GORM-specific tags based on column properties
	if c.Primary {
		tags = append(tags, "primaryKey")
	}
	if c.Name == "" {
		return ""
	}
	tags = append(tags, "column:"+c.Name)
	tags = append(tags, "type:"+c.Type.Postgres)
	if c.Increment {
		tags = append(tags, "autoIncrement")
	}
	if c.Default != "" {
		tags = append(tags, "default:"+c.Default)
	}
	if c.Unique {
		tags = append(tags, "unique")
	}
	for _, u := range c.Uniques {
		tags = append(tags, fmt.Sprintf("uniqueIndex:%s", u))
	}
	if c.Index {
		tags = append(tags, "index")
	}
	if c.Type.Serializer() {
		tags = append(tags, "serializer:json;")
	}
	if c.Type.SetNotNull {
		tags = append(tags, "not null")
	}
	// Add user-defined tags
	for k, v := range c.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}
	if len(tags) == 0 {
		return ""
	}

	tag := fmt.Sprintf(`gorm:"%v" json:"%v"`, strings.Join(tags, ";"), c.Name)
	if c.Sensitive {
		tag = fmt.Sprintf(`gorm:"%v" json:"-"`, strings.Join(tags, ";"))
	}
	return fmt.Sprintf("`%v`", tag)
}

func (c *Column) Clone() *Column {
	return &Column{
		Name:       c.Name,
		Desc:       c.Desc,
		Type:       c.Type,
		Tags:       c.Tags,
		Index:      c.Index,
		Unique:     c.Unique,
		Checks:     c.Checks,
		Uniques:    c.Uniques,
		Caches:     c.Caches,
		Primary:    c.Primary,
		Increment:  c.Increment,
		Step:       c.Step,
		Sensitive:  c.Sensitive,
		Validators: c.Validators,
		Relation:   c.Relation,
		modelType:  c.modelType,
		Enums:      c.Enums,
		Default:    c.Default,
	}
}

func (c *Column) SetModelType(typ string) {
	c.modelType = typ
}
func (c *Column) ModelType() string {
	if c.modelType != "" {
		return c.modelType
	}
	return c.Type.String()
}
func (r *Table) DefaultMixin() bool {
	for _, v := range r.Column {
		if v.Primary {
			return false
		}
	}
	for _, v := range r.Mixin {
		for _, v1 := range v.Column {
			if v1.Primary {
				return false
			}
		}
	}
	return true
}

var namer = schema.NamingStrategy{}

func (r *Table) TableName() string {
	return namer.TableName(r.Name)
}
func (r *Table) CacheKey() []*Table {
	var data = make(map[string]*Table, 0)
	var ret = make([]*Table, 0)
	for _, ve := range r.Column {
		v := ve.Clone()
		if v.Primary {
			data[v.Name] = &Table{
				Name: ToName(v.Name),
				Column: []*Column{
					v,
				},
			}
		}
		for _, v1 := range v.Caches {
			var arr = strings.Split(v1, ":")
			var name = arr[0]
			key := ToName(name)
			if len(arr) > 1 {
				v.Default = arr[1]
			} else {
				v.Default = ""
			}
			if _, ok := data[key]; !ok {
				data[key] = &Table{
					Name: key,
					Column: []*Column{
						v,
					},
				}
			} else {
				data[key].Column = append(data[key].Column, v)
			}
		}
	}
	for _, v := range data {
		ret = append(ret, v)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})
	return ret
}
func (r *Enums) Type() string {
	ref := reflect.ValueOf(r.Default)
	if ref.Kind() == reflect.Ptr {
		return "*" + ref.Elem().Type().String()
	}
	return ref.Type().String()
}
func (r *Enum) String() string {
	switch v := r.Value.(type) {
	case string:
		return fmt.Sprintf(`"%v"`, v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
