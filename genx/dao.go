package genx

type DaoInput struct {
	DaoPath   string
	ModelPath string
	TypePath  string
	Data      []any
}
type FieldType struct {
	Name      string
	DBName    string
	FieldType string
	DataType  string
	DoType    string
}
