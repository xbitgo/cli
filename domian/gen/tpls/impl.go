package tpls

import (
	"bytes"
	"text/template"
	"xbit/domian/parser"
)

const implFuncTmpl = `{{.XName}}) {{.Name}}({{- range $i,$v:= .Params}}  {{- if ne $i 0}},{{- end}} {{$v.Name}} {{$v.Type}}{{- end}}) {{- if gt (len .Results) 0}} ({{- range $i,$v:= .Results}}  {{- if gt $i 0}},{{- end}} {{$v.Name}} {{$v.Type}}{{- end}} ) {{- end}} {`

const implTmpl = `

{{range .MethodList}}
func ({{.ImplName}} {{$.Name}}) {{.Name}}({{- range $i,$v:= .Params}}  {{- if ne $i 0}},{{- end}} {{$v.Name}} {{$v.Type}}{{- end}}) {{- if gt (len .Results) 0}} ({{- range $i,$v:= .Results}}  {{- if gt $i 0}},{{- end}} {{$v.Name}} {{$v.Type}}{{- end}} ) {{- end}}{
	panic("implement me")
}
{{end}}
`

type IMPLFunc struct {
	XName string
	parser.XMethod
}

func (s *IMPLFunc) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New(s.Name + "INFFun").Parse(implFuncTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type INF struct {
	Body       []byte
	Name       string
	MethodList []parser.XMethod
}

func (s *INF) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New(s.Name + "INF").Parse(implTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return append(s.Body, buf.Bytes()...), nil
}
