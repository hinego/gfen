package test

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"reflect"
)

func Convent(data any) []any {
	h := &handler{}
	h.getAllMethods(reflect.ValueOf(data))
	return h.data
}

type handler struct {
	data []any
}

func (r *handler) getAllMethods(dataValue reflect.Value) {
	if r.checkAll(dataValue) == nil {
		r.data = append(r.data, dataValue.Interface())
		return
	}
	dataType := dataValue.Type()
	for i := 0; i < dataType.NumMethod(); i++ {
		if dataValue.Method(i).Type().NumIn() != 0 || dataValue.Method(i).Type().NumOut() != 1 {
			continue
		}
		results := dataValue.Method(i).Call(nil)
		for _, result := range results {
			r.getAllMethods(result)
		}
	}
	return
}
func (r *handler) checkAll(dataValue reflect.Value) (err error) {
	dataType := dataValue.Type()
	for i := 0; i < dataType.NumMethod(); i++ {
		method := dataType.Method(i)
		if err = r.check(method.Type); err != nil {
			return
		}
	}
	return
}
func (r *handler) check(reflectType reflect.Type) (err error) {
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		return gerror.New("invalid handler")
	}

	if !reflectType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return gerror.New("invalid handler")
	}

	if !reflectType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return gerror.New("invalid handler")
	}

	// The request struct should be named as `xxxReq`.
	if !gstr.HasSuffix(reflectType.In(1).String(), `Req`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for request: defined as "%s", but it should be named with "Req" suffix like "XxxReq"`,
			reflectType.In(1).String(),
		)
		return
	}

	// The response struct should be named as `xxxRes`.
	if !gstr.HasSuffix(reflectType.Out(0).String(), `Res`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for response: defined as "%s", but it should be named with "Res" suffix like "XxxRes"`,
			reflectType.Out(0).String(),
		)
		return
	}
	return
}
