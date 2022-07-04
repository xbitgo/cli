package tpls

import (
	"bytes"
	"text/template"
)

const do2Entity = `
package entity

// @AutoHelper 
type {{.Name}} struct {
{{- range .Fields}}
	{{.Name}} {{.Type}} {{.Tag }} // {{.Comment}}
{{- end}}
}
`

//TypeForMysqlToGo map for converting mysql type to golang types
var TypeForMysqlToGo = map[string]string{
	"int":                "int32",
	"integer":            "int32",
	"tinyint":            "int32",
	"smallint":           "int32",
	"mediumint":          "int32",
	"bigint":             "int64",
	"int unsigned":       "int32",
	"integer unsigned":   "int32",
	"tinyint unsigned":   "int32",
	"smallint unsigned":  "int32",
	"mediumint unsigned": "int32",
	"bigint unsigned":    "int64",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time",
	"datetime":           "time.Time",
	"timestamp":          "time.Time",
	"time":               "time.Time",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

type GenEntity struct {
	Name   string
	Fields []DoField
}

func (s *GenEntity) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New(s.Name + "Do2Entity").Parse(do2Entity)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
