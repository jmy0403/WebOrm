package WebOrm

const (
	BeforeInsert = "insertBefore"
	AfterInsert  = "afterInsert"
	BeforeUpdate = "beforeUpdate"
	AfterUpdate  = "afterUpdate"
	BeforeQuery  = "beforeQuery"
	AfterQuery   = "afterQuery"
	BeforeDelete = "beforeDelete"
	AfterDelete  = "afterDelete"
)

type fn func(*Conn)

func (e *Engine) SetHook(key string, hook fn) {
	e.hookMap[key] = hook
}
