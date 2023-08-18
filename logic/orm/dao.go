package orm

import (
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/hinego/gfen/genx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
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
type sOrm struct {
	db *gorm.DB
	*genx.DaoInput
	Map     *sync.Map
	schemas []*Schema
	Objects []string
	Import  map[string]string
}

func (r *sOrm) syncSchema() (err error) {
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
func (r *sOrm) Generate(data *genx.DaoInput) (err error) {
	r.DaoInput = data
	if err = r.syncSchema(); err != nil {
		return
	}
	var basename = gfile.Basename(r.TypePath)
	for _, object := range r.schemas {
		log.Println(basename, object)
	}
	return
}
func (r *sOrm) GenModel(data *genx.DaoInput) (err error) {
	return err
}
