package horm

import (
	"time"

	"github.com/hinego/decimal"
	"gorm.io/plugin/soft_delete"
)

var (
	DeletedAt = &Type{
		Name:     "Uint",
		Postgres: "bigint",
		Mysql:    "bigint",
		Sqlite:   "bigint",
		Type:     soft_delete.DeletedAt(0),
		Native:   true,
	}
	Int64 = &Type{
		Name:     "Int64",
		Postgres: "bigint",
		Mysql:    "bigint",
		Sqlite:   "bigint",
		Type:     int64(0),
	}
	String = &Type{
		Name:     "String",
		Postgres: "text",
		Mysql:    "varchar",
		Sqlite:   "varchar",
		Type:     string(""),
	}
	Decimal = &Type{
		Name:     "Field",
		Postgres: "decimal",
		Mysql:    "decimal",
		Sqlite:   "decimal",
		Type:     decimal.Decimal{},
	}
	Bool = &Type{
		Name:     "Bool",
		Postgres: "boolean",
		Mysql:    "tinyint(1)", // MySQL 使用 tinyint(1) 表示布尔值
		Sqlite:   "INTEGER",    // SQLite 使用 INTEGER 表示布尔值，0 代表 false，非0 代表 true
		Type:     false,
		value:    nil,
	}
	Int = &Type{
		Name:     "Int",
		Postgres: "integer",
		Mysql:    "int",
		Sqlite:   "integer",
		Type:     int(0),
	}
	Int32 = &Type{
		Name:     "Int32",
		Postgres: "integer",
		Mysql:    "int",
		Sqlite:   "integer",
		Type:     int32(0),
	}
	Float32 = &Type{
		Name:     "Float32",
		Postgres: "real",
		Mysql:    "float",
		Sqlite:   "real",
		Type:     float32(0),
	}
	Float64 = &Type{
		Name:     "Float64",
		Postgres: "double precision",
		Mysql:    "double",
		Sqlite:   "double",
		Type:     float64(0),
	}
	Byte = &Type{
		Name:     "Byte",
		Postgres: "smallint",
		Mysql:    "tinyint",
		Sqlite:   "tinyint",
		Type:     byte(0),
	}
	Bytes = &Type{
		Name:     "Bytes",
		Postgres: "bytea",
		Mysql:    "blob",
		Sqlite:   "blob",
		Type:     []byte{},
	}
	Time = &Type{
		Name:     "Time",
		Postgres: "timestamp",
		Mysql:    "datetime",
		Sqlite:   "datetime",
		Type:     time.Time{},
	}
	Uint = &Type{
		Name:     "Uint",
		Postgres: "integer",
		Mysql:    "int unsigned",
		Sqlite:   "integer",
		Type:     uint(0),
	}
	Uint8 = &Type{
		Name:     "Uint8",
		Postgres: "smallint",
		Mysql:    "tinyint unsigned",
		Sqlite:   "tinyint",
		Type:     uint8(0),
	}
	Uint16 = &Type{
		Name:     "Uint16",
		Postgres: "integer",
		Mysql:    "smallint unsigned",
		Sqlite:   "smallint",
		Type:     uint16(0),
	}
	Uint32 = &Type{
		Name:     "Uint32",
		Postgres: "integer",
		Mysql:    "int unsigned",
		Sqlite:   "integer",
		Type:     uint32(0),
	}
	Uint64 = &Type{
		Name:     "Uint64",
		Postgres: "bigint",
		Mysql:    "bigint unsigned",
		Sqlite:   "bigint",
		Type:     uint64(0),
	}
)

func Json(typ any) *Type {
	return &Type{
		Name:     "Field",
		Postgres: "jsonb", // Postgres 推荐使用 jsonb 类型，因为它有更好的性能特性
		Mysql:    "json",  // MySQL 提供了原生的 json 数据类型
		Sqlite:   "text",  // SQLite 并不支持原生的 json 类型，但你可以使用文本字段来存储 JSON 数据
		Type:     typ,     // 这里我们使用 Go 的 string 类型，但在实际应用中，你可能需要对应的结构或 map 来处理这个数据
	}
}
