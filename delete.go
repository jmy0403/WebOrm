package WebOrm

type delSql struct {
	table string
	rule  string
	args  []interface{}
	db    *Conn
}

func (c *Conn) DeleteFrom(table string) *delSql {
	return &delSql{table: table, db: c}
}
func (del *delSql) Rule(value string, args ...interface{}) {
	del.rule = value
	del.args = args
	conn := del.db
	if value, ok := conn.hook[BeforeDelete]; ok {
		value(conn)
	}
	defer func(c *Conn) {
		if value, ok := c.hook[AfterDelete]; ok {
			value(conn)
		}
	}(conn)
	delExec(del)

}

func delExec(del *delSql) (int64, error) {
	sql := "delete From " + del.table + " where " + del.rule
	result, err := del.db.DB().Exec(sql, del.args...)
	ShowSql(sql, err, del.args...)
	if err != nil {
		return 0, err
	}
	affected, err1 := result.RowsAffected()
	if err1 != nil {
		return affected, err
	}
	return affected, nil
}
