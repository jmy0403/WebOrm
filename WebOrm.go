package WebOrm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	logger "github.com/jmy0403/WebOrm/log"
)

type Engine struct {
	db *sql.DB
	//Conn
	hookMap map[string]fn
}

func NewConfig() *Config {
	return &Config{}
}

var IsShow bool

func NewEngine(cfg *Config) (*Engine, error) {
	//logger.SetLevel(cfg.logLervel)
	//"mysql", "root:qwe!23@tcp(127.0.0.1:3306)/nulige?charset=utf8"

	if cfg.Ip == "" {
		cfg.Ip = "127.0.0.1"
	}
	IsShow = cfg.IsShowSql
	if cfg.Port == "" {
		cfg.Port = "3306"
	}

	db, err := sql.Open("mysql", cfg.User+":"+cfg.Password+"@tcp("+cfg.Ip+cfg.Port+")"+cfg.DatabaseName)
	if err != nil {
		logger.Panic("数据库连接失败", err)
	}
	err1 := db.Ping()
	if err1 != nil {
		logger.Error("数据库ping失败", err)
		return nil, err1
	}
	e := &Engine{db: db, hookMap: make(map[string]fn)}
	logger.Info("数据库连接成功")
	return e, nil

}
func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		logger.Error("数据库连接关闭失败", err)
	}
	logger.Info("数据库成功关闭")
}
func (e *Engine) NewConn() *Conn {
	return NewConn(e.db, e.hookMap)
}
