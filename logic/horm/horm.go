package horm

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/horm"
	"github.com/hinego/gfen/ssr"
	"strings"
)

type sHorm struct {
	*horm.Input
	relation map[string]map[string]*horm.Relation
}

func (r *sHorm) Generate(in *horm.Input) (err error) {
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
				"Imports": imports,
				"Package": gfile.Basename(r.Path),
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
	return
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
			},
			File: fmt.Sprintf("%s/field/%ss/%v.gen.go", r.Path, v.Name, v.Name),
			Must: true,
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
			r.Table[k].Mixin = append(r.Table[k].Mixin, horm.Model)
		}
	}
	for k, _ := range r.Table {
		if len(r.Table[k].Mixin) == 0 {
			continue
		}
		for k1, _ := range r.Table[k].Mixin {
			r.Table[k].Column = append(r.Table[k].Mixin[k1].Column, r.Table[k].Column...)
		}
	}
	for k, v := range r.Table {
		for _, v1 := range v.Column {
			if v1.Relation != nil && v1.Relation.Type == horm.BelongsTo {
				var ref = v1.Relation.Reverse()
				r.setRelation(ref)
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
	return
}
func (r *sHorm) data() (data map[string]any) {
	return map[string]any{
		"Table": r.Table,
		"Path":  r.Path,
	}
}
