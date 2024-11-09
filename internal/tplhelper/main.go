package tplhelper

import (
	"html/template"
	"io"
)

func RenderWithPanic(tpl *template.Template, wr io.Writer, name string, data any) {
	err := tpl.ExecuteTemplate(wr, name, data)
	if err != nil {
		panic(err)
	}
}
