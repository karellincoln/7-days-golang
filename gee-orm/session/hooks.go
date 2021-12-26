package session

import (
	"github.com/karellincoln/7-day-golang/gee-orm/log"
	"reflect"
)

// Hooks constants
const (
	BeforeQuery = "BeforeQuery"
	AfterQuery  = "AfterQuery"

	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"

	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"

	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// 这种通过反射的方法，使用者基本不感知，其调用流程
// CallMethod calls the registered hooks
func (s *Session) CallMethod(method string, value interface{}) {
	// 找到Session的model对象中的method方法。
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	// 接受一个参数为session，并且期待返回为一个err的函数类型
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() { // 如果value的类型定义了这个方法，则调用。
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}
