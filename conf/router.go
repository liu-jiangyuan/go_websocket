package conf

import "github.com/liu-jiangyuan/go_websocket/controller"

//这个配置只针对websocket
func InitRoute() map[string]func(map[string]interface{})map[string]interface{} {
	return map[string]func(map[string]interface{}) map[string]interface{}{
		"Index": controller.Index,
	}
}
//type Handler struct {
//	Func  reflect.Value
//	In   reflect.Type
//	NumIn int
//	Out  reflect.Type
//	NumOut int
//}
//func InitRouter() {
	//handlers := make(map[string]*lib.Handler)
	//v := reflect.ValueOf(&controller.Index{})
	//t := reflect.TypeOf(&controller.Index{})
	//for i := 0; i < v.NumMethod(); i++ {
	//	name := t.Method(i).Name
	//	// 可以根据 i 来获取实例的方法，也可以用 v.MethodByName(name) 获取
	//	m := v.Method(i)
	//	// 这个例子我们只获取第一个输入参数和第一个返回参数
	//	in := m.Type().In(0)
	//	out := m.Type().Out(0)
	//	handlers[name] = &lib.Handler{
	//		Func:  m,
	//		In:   in,
	//		NumIn: m.Type().NumIn(),
	//		Out:  out,
	//		NumOut: m.Type().NumOut(),
	//	}
	//}
	//return handlers
	//inVal := reflect.New(r[c].In).Elem()
	//rtn := r[c].Func.Call([]reflect.Value{inVal})[0]
//}