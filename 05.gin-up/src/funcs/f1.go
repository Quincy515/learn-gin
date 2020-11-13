package funcs

import "html/template"

func Strong(txt string) template.HTML {
	return template.HTML("<strong>" + txt + "</strong>")
}
func Test() string {
	return "test"
}
