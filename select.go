package WebOrm

import (
	sqll "database/sql"
	"fmt"
	logger "github.com/jmy0403/WebOrm/log"
	"reflect"
	"strings"
)

type QuerySql struct {
	alias     []string
	resault   interface{}
	tableName string
	arg       []interface{}
	rule      string
	conn      *Conn
}

func (c *Conn) Select(value interface{}, alias ...string) *QuerySql {

	return &QuerySql{
		alias:   alias,
		resault: value,
		conn:    c,
	}
}
func (query *QuerySql) From(table string) *QuerySql {
	query.tableName = table
	return query
}
func (c *QuerySql) Rule(queryRule string, args ...interface{}) {

	resault := c.resault
	c.rule = queryRule
	conn := c.conn

	if value, ok := conn.hook[BeforeQuery]; ok {
		value(conn)
	}
	defer func(c *Conn) {
		if value, ok := conn.hook[AfterQuery]; ok {
			value(c)
		}
	}(conn)
	switch reflect.ValueOf(resault).Elem().Kind() {
	case reflect.Struct:
		{

			c.exec(c, args...)
		}
	case reflect.Slice:
		{
			c.execs(c, args...)
		}

	}
}

func (query *QuerySql) exec(c *QuerySql, args ...interface{}) {
	var sql strings.Builder
	var join string
	sql.WriteString("select ")
	if c.alias != nil {
		sql.WriteString(c.alias[0])
	} else {
		values := reflect.TypeOf(c.resault).Elem().Elem()
		v := reflect.New(values)
		value := v.Elem()
		str := make([]string, value.NumField())
		for i := 0; i < value.NumField(); i++ {
			str[i] = value.Type().Field(i).Name
		}
		join = strings.Join(str, ",")
	}
	var row *sqll.Row

	sql.WriteString(join)
	sql.WriteString(" from " + c.tableName + " " + c.rule)

	if args == nil {

		row = c.conn.DB().QueryRow(sql.String() + ";")
	} else {
		row = c.conn.DB().QueryRow(sql.String()+";", args...)
	}

	fmt.Println(sql.String(), args)
	queryParse(c.resault, row)

}

func (query *QuerySql) execs(c *QuerySql, args ...interface{}) {
	var sql strings.Builder
	var join string
	sql.WriteString("select ")
	if c.alias != nil {
		sql.WriteString(c.alias[0])
	} else {
		value := reflect.TypeOf(c.resault).Elem().Elem()

		str := make([]string, value.NumField())
		for i := 0; i < value.NumField(); i++ {
			str[i] = value.Field(i).Name
		}
		join = strings.Join(str, ",")
	}

	sql.WriteString(join)
	sql.WriteString(" from " + c.tableName + " " + c.rule + ";")

	r, err := c.conn.DB().Query(sql.String(), args...)
	ShowSql(sql.String(), err, args...)
	if err != nil {
		logger.Error(err)
	}
	rows := r

	queryParseSlice(c.resault, rows)
}

func queryParseSlice(value interface{}, rows *sqll.Rows) {
	typeOf := reflect.TypeOf(value).Elem()

	slice := reflect.Indirect(reflect.ValueOf(value))

	elem := typeOf.Elem()
	//fmt.Println(elem.Kind())
	for rows.Next() {
		resault := reflect.New(elem)

		stu := resault.Elem()

		var args []interface{}
		for i := 0; i < stu.NumField(); i++ {
			iface := stu.Field(i).Addr().Interface()
			args = append(args, iface)
		}
		rows.Scan(args...)
		slice.Set(reflect.Append(slice, stu))
	}
}
func queryParse(value interface{}, row *sqll.Row) {

	var arr []interface{}
	valueOf := reflect.ValueOf(value)
	typeOf := valueOf.Type()
	v := reflect.New(typeOf.Elem())

	for i := 0; i < v.Elem().NumField(); i++ {
		iface := v.Elem().Field(i).Addr().Interface()
		arr = append(arr, iface)
	}

	err := row.Scan(arr...)
	if err != nil {
		fmt.Println(err)
	}

	valueOf.Elem().Set(v.Elem())

}
