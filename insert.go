package WebOrm

//
import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//
func (c *Conn) Insert(tableName string, value interface{}) (int, error) {
	if value, ok := c.hook[BeforeInsert]; ok {
		value(c)
	}
	defer func(c *Conn) {
		if value, ok := c.hook[AfterInsert]; ok {
			value(c)
		}
	}(c)
	switch reflect.ValueOf(value).Kind() {
	case reflect.Struct, reflect.Ptr:
		return c.insertOne(tableName, value)
	case reflect.Slice, reflect.Array:
		return c.insertLot(tableName, value)

	}
	return 0, errors.New("插入失败")
}
func (c *Conn) insertLot(tableName string, value interface{}) (int, error) {
	var join string
	var sql string
	var bufStr strings.Builder
	bufStr.WriteString("insert into " + tableName)

	join = parseSlice(value)
	bufStr.WriteString(" value " + join)
	sql = bufStr.String()

	line, err := c.Exec(sql)
	if err != nil {
		return 0, err
	}
	return int(line), err
}
func parseStu(value interface{}) string {
	fields := reflect.Indirect(reflect.ValueOf(value))
	str := make([]string, fields.NumField())
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		if field.Kind() == reflect.Int {
			str[i] = strconv.Itoa(int(field.Int()))
		} else {
			str[i] = "'" + field.String() + "'"
		}

	}

	join := strings.Join(str, ",")
	return join
}
func parseSlice(stus interface{}) string {
	var bulider strings.Builder
	value := reflect.ValueOf(stus)
	if value.Index(0).Kind() == reflect.Ptr {
		return parseSlicePtr(value, &bulider)
	} else {
		return parseSliceStu(value, &bulider)
	}

}

func parseSliceStu(value reflect.Value, bulider *strings.Builder) string {
	len := value.Len()
	for i := 0; i < len; i++ {

		elem := value.Index(i)
		num := elem.NumField()
		str := make([]string, num)
		for j := 0; j < num; j++ {
			switch elem.Field(j).Kind() {
			case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32,
				reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				str[j] = strconv.Itoa(int(elem.Field(j).Int()))
			case reflect.String:
				str[j] = "'" + elem.Field(j).String() + "'"
			case reflect.Float32, reflect.Float64:
				str[j] = fmt.Sprintf("%f", elem.Field(j).Float())
			case reflect.Bool:
				if elem.Field(j).Bool() == true {
					str[j] = "1"
				} else {
					str[j] = "0"
				}
			}
		}
		join := strings.Join(str, ",")
		if i == len-1 {
			bulider.WriteString("(" + join + ")" + ";")
		} else {
			bulider.WriteString("(" + join + ")" + ",")
		}

	}
	return bulider.String()
}

func parseSlicePtr(value reflect.Value, bulider *strings.Builder) string {
	len := value.Len()
	for i := 0; i < len; i++ {

		elem := value.Index(i).Elem()
		num := elem.NumField()
		str := make([]string, num)
		for j := 0; j < num; j++ {
			switch elem.Field(j).Kind() {
			case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32,
				reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				str[j] = strconv.Itoa(int(elem.Field(j).Int()))
			case reflect.String:
				str[j] = "'" + elem.Field(j).String() + "'"
			case reflect.Float32, reflect.Float64:
				str[j] = fmt.Sprintf("%f", elem.Field(j).Float())
			case reflect.Bool:
				if elem.Field(j).Bool() == true {
					str[j] = "1"
				} else {
					str[j] = "0"
				}
			}
		}
		join := strings.Join(str, ",")
		if i == len-1 {
			bulider.WriteString("(" + join + ")" + ";")
		} else {
			bulider.WriteString("(" + join + ")" + ",")
		}

	}
	return bulider.String()
}

//插入一条数据数据,可直接传入指针或者结构体对象，例如：User{}，&User{}
func (c *Conn) insertOne(tableName string, value interface{}) (int, error) {
	var line int64
	var join string
	var sql string
	var bufStr strings.Builder
	bufStr.WriteString("insert into " + tableName)

	join = parseStu(value)
	bufStr.WriteString(" value (" + join + ");")
	sql = bufStr.String()

	line, err := c.Exec(sql)
	if err != nil {
		return 0, err
	}
	return int(line), err
}
