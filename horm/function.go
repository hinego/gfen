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
	valType := reflect.TypeOf(r.Type)
	if r.Package() == "" {
		switch valType.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array:
			return true
		}
		return false
	}
	if r.Native {
		return false
	}
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
	var data = getPackageDeps(reflect.TypeOf(r.Type))
	if len(data) > 0 {
		return data[0]
	}
	return ""
}

func getPackageDeps(t reflect.Type) []string {
	packages := make(map[string]bool)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	pkg := t.PkgPath()
	if pkg != "" {
		packages[pkg] = true
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		for _, p := range getPackageDeps(t.Elem()) {
			packages[p] = true
		}
	case reflect.Map:
		for _, p := range getPackageDeps(t.Key()) {
			packages[p] = true
		}
		for _, p := range getPackageDeps(t.Elem()) {
			packages[p] = true
		}
	// case reflect.Struct:
	// 	for i := 0; i < t.NumField(); i++ {
	// 		for _, p := range getPackageDeps(t.Field(i).Type) {
	// 			packages[p] = true
	// 		}
	// 	}
	case reflect.Chan:
		for _, p := range getPackageDeps(t.Elem()) {
			packages[p] = true
		}
	case reflect.Func:
		for i := 0; i < t.NumIn(); i++ {
			for _, p := range getPackageDeps(t.In(i)) {
				packages[p] = true
			}
		}
		for i := 0; i < t.NumOut(); i++ {
			for _, p := range getPackageDeps(t.Out(i)) {
				packages[p] = true
			}
		}
	}

	var result []string
	for p := range packages {
		if !strings.Contains(p, "builtin") { // filter out built-in packages
			result = append(result, p)
		}
	}
	return result
}
func GenTag(data map[string]any) string {
	var builder strings.Builder
	var tagSlice []string
	for k, v := range data {
		var tags []string
		switch e := v.(type) {
		case map[string]string:
			for k1, v1 := range e {
				if v1 == "" {
					tags = append(tags, k1)
				} else {
					tags = append(tags, fmt.Sprintf("%s:%s", k1, v1))
				}
			}
		case string:
			tags = append(tags, e)
		}

		if len(tags) == 0 {
			continue
		}
		tag := fmt.Sprintf(`%v:"%v"`, k, strings.Join(tags, ";"))
		tagSlice = append(tagSlice, tag)
		// builder.WriteString(tag)
		// builder.WriteString(" ")
	}
	sort.Slice(tagSlice, func(i, j int) bool {
		return tagSlice[i] > tagSlice[j]
	})
	builder.WriteString("`")
	for i, v := range tagSlice {
		if i != 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(v)
	}
	builder.WriteString("`")
	return builder.String()
}

func (c *Relation) Tag() string {
	var mapGorm = map[string]string{
		"foreignKey": ToName(c.Foreign),
		"references": ToName(c.Reference),
	}
	var tags = map[string]any{
		"gorm": mapGorm,
	}
	return GenTag(tags)
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
	var tagMap = map[string]any{
		"gorm": strings.Join(tags, ";"),
		"json": c.Name,
	}
	if c.Desc != "" {
		tagMap["dc"] = c.Desc
	} else {
		tagMap["dc"] = c.Name
	}
	if c.Relation != nil && c.Table != "-" {
		tagMap["table"] = c.Relation.RefTable
	} else if c.Table != "" {
		tagMap["table"] = c.Table
	}

	if c.Filter != "" {
		tagMap["filter"] = c.Filter
	}
	if c.Config != "" {
		tagMap["config"] = c.Config
	}
	if c.Sensitive {
		tagMap["json"] = "-"
	}
	if c.HideTable {
		tagMap["hideTable"] = "true"
	}
	if c.HideSearch {
		tagMap["hideSearch"] = "true"
	}
	if c.Ellipsis {
		tagMap["ellipsis"] = "true"
	}
	if c.Ts != "" {
		tagMap["ts"] = c.Ts
	}
	return GenTag(tagMap)
	// tag := fmt.Sprintf(`gorm:"%v" json:"%v"`, strings.Join(tags, ";"), c.Name)
	// if c.Sensitive {
	// 	tag = fmt.Sprintf(`gorm:"%v" json:"-"`, strings.Join(tags, ";"))
	// }
	// return fmt.Sprintf("`%v`", tag)
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
func (r *Table) HasEnum() bool {
	for _, v := range r.Column {
		if v.Enums != nil {
			return true
		}
	}
	return false
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
func (r *Enum) Typescript() string {
	switch v := r.Value.(type) {
	case string:
		return fmt.Sprintf(`"%v"`, v)
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	default:
		return fmt.Sprintf("%v", v)
	}
}
