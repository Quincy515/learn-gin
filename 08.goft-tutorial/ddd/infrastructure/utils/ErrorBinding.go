package utils

type ErrorResult struct {
	data interface{}
	err  error
}

func NewErrorResult(data interface{}, err error) *ErrorResult {
	return &ErrorResult{data: data, err: err}
}

// Unwrap 有错误 panic 没有错误就返回 data
func (e *ErrorResult) Unwrap() interface{} {
	if e.err != nil {
		panic(e.err.Error())
	}
	return e.data
}

type BindFunc func(v interface{}) error

// Exec 统一处理 gin 里 shouldBuild 的 error
func Exec(f BindFunc, value interface{}) *ErrorResult {
	err := f(value)
	return NewErrorResult(value, err)
}
