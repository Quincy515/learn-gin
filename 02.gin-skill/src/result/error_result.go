package result

import (
	"fmt"
	"ginskill/src/validators"
)

type ErrorResult struct {
	data interface{}
	err  error
}

func (this *ErrorResult) Unwrap() interface{} {
	if this.err != nil {
		validators.CheckErrors(this.err) // 如果匹配到这里就会 panic
		panic(this.err.Error())          // 没有匹配到继续走这个 panic
	}
	return this.data
}

func (this *ErrorResult) UnwrapOr(v interface{}) interface{} {
	if this.err != nil {
		return v // 如果有错误直接返回 v
	}
	return this.data
}

func (this *ErrorResult) UnwrapOrElse(f func() interface{}) interface{} {
	if this.err != nil {
		return f() // 如果有错误执行函数 f，返回 interface{}
	}
	return this.data
}

func Result(vs ...interface{}) *ErrorResult {
	if len(vs) == 1 {
		if vs[0] == nil {
			return &ErrorResult{nil, nil}
		}
		if e, ok := vs[0].(error); ok {
			return &ErrorResult{nil, e}
		}
	}
	if len(vs) == 2 { // 如果可变参数有2个判断第2个参数是否有错
		if vs[1] == nil {
			return &ErrorResult{vs[0], nil}
		}
		if e, ok := vs[1].(error); ok {
			return &ErrorResult{vs[0], e}
		}
	}
	return &ErrorResult{nil, fmt.Errorf("error result format")}
}
