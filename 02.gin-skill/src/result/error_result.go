package result

type ErrorResult struct {
	err error
}

func (this *ErrorResult) Unwrap() interface{} {
	if this.err != nil {
		panic(this.err.Error())
	}
	return nil
}

func Result(err error) *ErrorResult {
	return &ErrorResult{err: err}
}
