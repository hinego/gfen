package horm

type (
	Validator struct {
		Name      string // CHECK约束的名称
		Condition string // 验证条件
	}
	Type struct {
		Name       string //字段的名称
		Postgres   string //在Postgres中的类型
		Mysql      string //在Mysql中的类型
		Sqlite     string //在Sqlite中的类型
		Type       any    // 此类型在golang中的类型，直接传入该类型即可，例如 decimal.Decimal{} 或者 int64(0) 然后通过反射获取到类型和包名
		value      any    // 此字段的值，更新时会使用此值
		Native     bool   // 是否是原生类型 不需要进行json序列化
		SetNotNull bool   // 是否可以为空
		Point      bool   // 是否以指针使用指针
	}
	Relation struct {
		Name         string
		RefName      string
		Type         string //belongs_to has_one has_many，一旦创建A->B的关系，B->A的关系也会自动创建
		RefTable     string //关联的表名
		Table        string //此表的表名
		Fake         bool   // 虚拟外键 不会在数据库中创建外键
		Query        bool   //是否创建Query功能
		Foreign      string //外键的字段名
		ForeignPoint bool   //外键是否为指针
		Reference    string //外键引用的字段名 (通常是主键)
		Unique       bool   //是否唯一
		Desc         string
		OnUpdate     string
		OnDelete     string
	}
	Enum struct {
		Name  string
		Value any
		Desc  string
	}
	Enums struct {
		Default any
		Enums   []*Enum
	}
	Check struct {
		Name   string // CHECK约束的名称
		Clause string // 验证条件
	}
	Column struct {
		Name       string            //字段的名称
		Desc       string            //字段的描述
		Type       *Type             //字段的类型
		Tags       map[string]string //给结构体的Tags
		Index      bool              //此字段是否索引
		Unique     bool              //此字段是否唯一
		Checks     []*Check          //验证条件
		Uniques    []string          //联合唯一索引 例如：[]string{"type_index"}
		Primary    bool              //是否主键
		Increment  bool              //是否自增
		Step       int               //自增步长
		Sensitive  bool              //是否敏感字段
		Validators []*Validator
		Relation   *Relation
		modelType  string
		Enums      *Enums
		Default    string
	}
	Table struct {
		Name        string
		Primary     string //主键的字段名
		PrimaryType string //主键的类型
		Column      []*Column
		Mixin       []*Mixin
		Relation    map[string]*Relation
	}
	Mixin struct {
		Column []*Column
	}
	Input struct {
		Table []*Table //模型
		Path  string   //输出路径
	}
)
