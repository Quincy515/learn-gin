package goft

// Fairing 规范中间件代码和功能的接口
type Fairing interface {
	OnRequest() error
}
