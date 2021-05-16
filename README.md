* 创建配置结构体
```go
cfg = WebOrm.NewConfig()
cfg.User = "用户名称"
cfg.Password = "密码"
cfg.Database = "数据库名称"
//ip和port 有默认值
```
* 连接数据库
```go
engine err := WebOrm.NewEngine(cfg);
//获取会话
conn = engine.NewConn()
```
* 钩子函数
```go
    engine.SetHook("AfterInsert",fn)
```
* 创建数据表
```go
//value是与数据表相映射的结构体，arg是表的一些设置，比如"engine = innodb"
func (c *Conn) CreateTable(value interface{}, arg string) error 
type User struct{
//如果是字符串类型要在tag写上varchar()
    id int 'orm:"primary key"'
    name string `orm:"varchar(20) , not null"`
    age int
}
conn.CreateTable(&User{},"")
```
* Conn用来控制数据的增删改查
```go
select查询只需要传入一个与查询结果映射的结构体,结构体数组或者切片.然后用链式操作生成sql语句
delete删除只需要传入结构体，结构体切片。
update需要拼接sql语句
```

* 事务
```go
//fn为具体的事务操作，Begin后会自动回滚或者提交。
    conn.Begin(fn)
```