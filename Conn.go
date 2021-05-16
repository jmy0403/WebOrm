package WebOrm

import (
	"database/sql"
	"fmt"

	logger "github.com/jmy0403/WebOrm/log"
	"github.com/jmy0403/WebOrm/stuTrans"
	"github.com/jmy0403/WebOrm/typeTrans"
	"reflect"
	"strings"
)

type Conn struct {
	db       *sql.DB
	refTable *stuTrans.DataTable
	hook     map[string]fn
}
type txFun func(*Conn) (interface{}, error)

//开启事务，自动回滚和提交
func (c *Conn) Begin(fn txFun) (interface{}, error) {
	begin, err := c.db.Begin()

	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = begin.Rollback()
			panic(p)
		} else if err != nil {
			_ = begin.Rollback()
		} else {
			err = begin.Commit()
		}
	}()
	return fn(c)
}

func NewConn(db *sql.DB, hoo map[string]fn) *Conn {

	return &Conn{
		db:   db,
		hook: hoo,
	}
}

func (c *Conn) DB() *sql.DB {
	return c.db
}

//执行SQL语句，返回影响到的行数
func (c *Conn) Exec(SqlStr string, args ...interface{}) (line int64, err error) {
	result, err := c.db.Exec(SqlStr, args...)

	if err != nil {
		logger.Error("SQL语句执行失败")
		return 0, err
	}
	//var line int64
	line, _ = result.RowsAffected()
	ShowSql(SqlStr, err, args...)
	return line, nil
}

func ShowSql(str string, err error, args ...interface{}) {
	if IsShow == true {
		all := strings.ReplaceAll(str, "?", "%v")
		sprintf := fmt.Sprintf(all, args...)
		if err == nil {
			logger.Info("SQL语句执行成功: ", sprintf)
		} else {
			logger.Error("SQL语句执行失败：", sprintf)
		}

	}

}

func (c *Conn) stuParse(value interface{}) *Conn {
	if c.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(c.refTable.Model) {
		c.refTable = stuTrans.Parse(value, typeTrans.GetTrans())
		return c
	}
	logger.Error("结构体转换数据表失败")
	return c
}
func (c *Conn) GetTable() *stuTrans.DataTable {
	if c.refTable == nil {
		logger.Error("无数据表")
	}
	return c.refTable
}

//value是要创建表的结构体，args是表后面的属性 ，默认是ENGINE=INNODB CHARACTER SET utf8mb4
func (c *Conn) CreateTable(value interface{}, arg string) error {

	c.stuParse(value)
	table := c.refTable
	var builder strings.Builder
	str := make([]string, len(table.Fields))

	builder.WriteString("create table " + table.TableName + "(")
	for i := 0; i < len(table.Fields); i++ {
		field := table.Fields[i]
		str[i] = field.Name + " " + field.Type + " " + field.Tag

	}
	join := strings.Join(str, ",")

	if arg == "" {
		builder.WriteString(join + ")" + "ENGINE=INNODB CHARACTER SET utf8mb4;")
	} else {
		builder.WriteString(join + ")" + arg + ";")
	}

	sql := builder.String()
	fmt.Println("Sql:", sql)
	_, err := c.Exec(sql)
	return err
}
func (c *Conn) DroupTable(tableName string) error {
	str := fmt.Sprintf("drop table if exists %s ;", tableName)
	_, err := c.Exec(str)
	return err
}

//判断表是否存在
func (c *Conn) HasTable(database string, tableName string) bool {

	//var str = "select TABLE_NAME from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA='geeorm' and TABLE_NAME='person'"
	var str = "select TABLE_NAME from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA='" + database + "' and TABLE_NAME='" + tableName + "'"
	var has string

	row := c.db.QueryRow(str)
	ShowSql(str, nil)
	err := row.Scan(&has)
	if err != nil {
		logger.Error(err)
	}

	return has == tableName

}

//删除表中的数据，不删除表
func (c *Conn) Truncate(tableName string) (int64, error) {
	line, err := c.Exec("truncate table ?", tableName)
	if err != nil {
		return 0, err

	}
	return line, err
}
