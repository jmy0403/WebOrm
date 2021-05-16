package stuTrans

import (
	"github.com/jmy0403/WebOrm/typeTrans"
	"go/ast"
	"reflect"
)

//结构体的每一个字段
type Field struct {
	Name string
	Type string
	Tag  string
}
type DataTable struct {
	Model      interface{}
	TableName  string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func Parse(Stu interface{}, transForm typeTrans.Trans) *DataTable {
	//获取表名称
	modelType := reflect.Indirect(reflect.ValueOf(Stu)).Type()
	tableName := modelType.Name()
	table := &DataTable{
		Model:     Stu,
		TableName: tableName,
		fieldMap:  make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)

		if !ast.IsExported(p.Name) {

			typeof := transForm.DataTypeof(reflect.Indirect(reflect.New(p.Type)))
			if typeof == "text" {
				typeof = ""
			}
			field := &Field{
				Name: p.Name,
				Type: typeof,
			}
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}
			table.Fields = append(table.Fields, field)
			table.FieldNames = append(table.FieldNames, p.Name)
			table.fieldMap[p.Name] = field
		}

	}
	return table

}
