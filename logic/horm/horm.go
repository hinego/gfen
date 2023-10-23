package horm

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/horm"
	"github.com/hinego/gfen/ssr"
)

type sHorm struct {
	*horm.Input
	relation map[string]map[string]*horm.Relation
}

func (r *sHorm) Unique(in *horm.Input) (data *horm.Input) {
	var unique = make(map[string]map[string]struct{})
	data = &horm.Input{
		Table: make([]*horm.Table, 0),
		Path:  in.Path,
	}
	for _, v := range in.Table {
		if _, ok := unique[v.Name]; !ok {
			unique[v.Name] = make(map[string]struct{})
			var tab = &horm.Table{
				Name:        v.Name,
				Column:      make([]*horm.Column, 0),
				Primary:     v.Primary,
				PrimaryType: v.PrimaryType,
				Mixin:       v.Mixin,
				CacheLevel:  v.CacheLevel,
				Relation:    v.Relation,
			}
			for _, v1 := range v.Column {
				if _, ok1 := unique[v.Name][v1.Name]; !ok1 {
					unique[v.Name][v1.Name] = struct{}{}
					tab.Column = append(tab.Column, v1)
				}
			}
			data.Table = append(data.Table, tab)
		}
	}
	return
}
func (r *sHorm) Generate(in *horm.Input) (err error) {
	in = r.Unique(in)
	r.Input = in
	r.fill()
	if err = r.field(); err != nil {
		return err
	}
	if err = r.model(); err != nil {
		return err
	}
	if err = r.migrate(); err != nil {
		return err
	}
	if err = r.enmus(); err != nil {
		return err
	}
	if err = r.typescript(); err != nil {
		return err
	}
	if err = r.mapping(); err != nil {
		return err
	}

	if err = r.dao(); err != nil {
		return err
	}
	if err = r.gen(); err != nil {
		return err
	}
	return
}
func (r *sHorm) migrate() (err error) {
	var imports = make(map[string]string)
	for _, v := range r.Table {
		for _, v1 := range v.Column {
			path := v1.Type.Package()
			if path != "" {
				imports[path] = path
			}
		}
	}
	return ssr.Gen().Execute(&genx.Execute{
		Code: MigrateTemplate,
		Data: r.data(),
		Map: map[string]any{
			"Imports": imports,
		},
		File: fmt.Sprintf("%s/migrate/migrate.gen.go", r.Path),
		Must: true,
	})
}
func (r *sHorm) enmus() (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: EnumTemplate,
		Data: r.data(),
		File: fmt.Sprintf("%s/enums/enums.gen.go", r.Path),
		Must: true,
	})
}
func (r *sHorm) typescript() (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: EnumTypeTemplate,
		Data: r.data(),
		File: fmt.Sprintf("%s/enums/enums.gen.ts", r.Path),
		Must: true,
	})
}
func (r *sHorm) mapping() (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: MappingTemplate,
		Data: r.data(),
		Map: map[string]any{
			"Name": "db",
			"Dao":  gfile.Basename(r.Path),
		},
		File: fmt.Sprintf("%s/db/db.gen.go", r.Path),
		Must: true,
	})
}
func (r *sHorm) model() (err error) {
	self := ssr.Gen().Path(r.Path)
	base := gfile.Basename(self)
	for k, v := range r.Table {
		var imports = make(map[string]string)
		for _, v1 := range v.Relation {
			str := ssr.Gen().Path(r.Path, "field", v1.RefTable+"s")
			imports[str] = str
		}
		for k1, v1 := range v.Column {
			path := v1.Type.Package()
			if self == path {
				r.Table[k].Column[k1].SetModelType(strings.ReplaceAll(v1.Type.String(), base+".", ""))
			}
			if self != path && path != "" {
				imports[path] = path
			}
		}
		if err = ssr.Gen().Execute(&genx.Execute{
			Code: ModelTemplate,
			Data: v,
			Map: map[string]any{
				"Imports":  imports,
				"Package":  gfile.Basename(r.Path),
				"CacheKey": v.CacheKey(),
			},
			File: fmt.Sprintf("%s/%s.mod.gen.go", r.Path, v.Name),
			Must: true,
		}); err != nil {
			return err
		}
	}
	return err
}
func (r *sHorm) setRelation(data *horm.Relation) {
	if r.relation == nil {
		r.relation = make(map[string]map[string]*horm.Relation)
	}
	if r.relation[data.Table] == nil {
		r.relation[data.Table] = map[string]*horm.Relation{}
	}
	r.relation[data.Table][data.Name] = data
}
func (r *sHorm) dao() (err error) {
	for _, v := range r.Table {
		if err = ssr.Gen().Execute(&genx.Execute{
			Code: DaoTemplate,
			Data: map[string]any{
				"Model":       strings.Title(v.Name),
				"Table":       v.Name,
				"Primary":     v.Primary,
				"PrimaryType": v.PrimaryType,
				"Column":      v.Column,
				"Relation":    v.Relation,
				"Imports": []string{
					ssr.Gen().Path(r.Path, "field", v.Name+"s"),
				},
			},
			Map: map[string]any{
				"Package": gfile.Basename(r.Path),
			},
			File: fmt.Sprintf("%s/%s.dao.gen.go", r.Path, v.Name),
			Must: true,
		}); err != nil {
			return err
		}
	}
	return err
}
func (r *sHorm) field() (err error) {
	for _, v := range r.Table {
		if err = ssr.Gen().Execute(&genx.Execute{
			Code: FieldTemplate,
			Data: map[string]any{
				"Model":       strings.Title(v.Name),
				"Table":       v.Name,
				"TableName":   v.TableName(),
				"Primary":     v.Primary,
				"PrimaryType": v.PrimaryType,
				"Column":      v.Column,
				"CacheKey":    v.CacheKey(),
			},
			File: fmt.Sprintf("%s/field/%ss/%v.gen.go", r.Path, v.Name, v.Name),
			Must: true,
			// Debug: true,
		}); err != nil {
			return err
		}
	}
	return err
}
func (r *sHorm) gen() (err error) {
	return ssr.Gen().Execute(&genx.Execute{
		Code: GenTemplate,
		Data: r.data(),
		Map: map[string]any{
			"Package": gfile.Basename(r.Path),
		},
		File: fmt.Sprintf("%s/gen.gen.go", r.Path),
		Must: true,
	})
}
func (r *sHorm) fill() {
	for k, v := range r.Table {
		if v.DefaultMixin() {
			r.Table[k].Mixin = append(r.Table[k].Mixin, &horm.Model)
		}
	}
	for k := range r.Table {
		if len(r.Table[k].Mixin) == 0 {
			continue
		}
		for k1 := range r.Table[k].Mixin {
			r.Table[k].Column = append(r.Table[k].Mixin[k1].Column, r.Table[k].Column...)
		}
	}
	for k, v := range r.Table {
		for _, v1 := range v.Column {
			if v1.Relation != nil {
				v1.Relation.Type = horm.BelongsTo
				v1.Relation.Table = v.Name
				v1.Relation.Foreign = v1.Name
				if v1.Relation.Reference == "" {
					v1.Relation.Reference = "id"
				}
				v1.Relation.RefName = strings.Title(v.Name) + "s"
				var ref = v1.Relation.Reverse()
				r.setRelation(ref)
				if v1.Type.Point {
					v1.Relation.ForeignPoint = true
				}
				r.setRelation(v1.Relation)
			}
			if v1.Primary {
				r.Table[k].Primary = v1.Title()
				r.Table[k].PrimaryType = v1.Type.String()
			}
		}
	}
	for k, v := range r.Table {
		if data, ok := r.relation[v.Name]; ok {
			r.Table[k].Relation = data
		}
	}
}
func (r *sHorm) data() (data map[string]any) {
	return map[string]any{
		"Table": r.Table,
		"Path":  r.Path,
	}
}
