package WebOrm

import (
	"fmt"
)

type updateSql struct {
	table    string
	setSql   string
	setArgs  []interface{}
	ruleSql  string
	ruleArgs []interface{}
	db       *Conn
}

func (c *Conn) Update(tableName string) *updateSql {
	return &updateSql{
		table: tableName,
		db:    c,
	}
}
func (up *updateSql) Set(key string, args ...interface{}) *updateSql {
	up.setSql = key
	up.setArgs = args
	return up

}
func (up *updateSql) Where(key string, args ...interface{}) {
	up.ruleSql = key
	up.ruleArgs = args
	conn := up.db
	if value, ok := conn.hook[BeforeUpdate]; ok {
		value(conn)
	}
	defer func(c *Conn) {
		if value, ok := c.hook[AfterInsert]; ok {
			value(conn)
		}
	}(conn)
	upexes(up)
}
func upexes(up *updateSql) (int64, error) {
	sql := "update " + up.table + " Set " + up.setSql + " Where " + up.ruleSql
	up.setArgs = append(up.setArgs, up.ruleArgs...)

	exec, err := up.db.DB().Exec(sql, up.setArgs...)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	ShowSql(sql, err, up.setArgs...)
	affected, err1 := exec.RowsAffected()
	if err1 != nil {
		return 0, err1
	}
	return affected, nil
}
