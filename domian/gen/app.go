package gen

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"

	"xbit/conf"
	"xbit/domian/gen/tpls"
	"xbit/domian/parser"
)

type genObject struct {
	dir  string
	name string
	fun  func() ([]byte, error)
}

func (m *Manager) App(rootPath string, name string) error {
	appPath := fmt.Sprintf("%s/%s/%s", rootPath, conf.Global.Tmpl.AppDir, name)
	tmpl := conf.Global.GetRealTmpl(rootPath, appPath, name)
	m.Tmpl = tmpl
	rv := reflect.ValueOf(tmpl)
	for i := 0; i < rv.Elem().NumField(); i++ {
		rs := rv.Elem().Field(i).Interface()
		dir := rs.(string)
		if dir != "" {
			if err := os.MkdirAll(dir, 0766); err != nil {
				log.Printf("Mkdir[%s] err: %v \n", dir, err)
			}
		}
	}
	appGen := tpls.App{
		ProjectName:   conf.Global.ProjectName,
		ServiceName:   name,
		ServiceNameUF: tool_str.ToUFirst(name),
		AppsDir:       tmpl.AppDir,
	}
	initFiles := []genObject{
		{
			dir:  tmpl.PbDir,
			name: name + ".proto",
			fun:  appGen.Proto,
		},
		{
			dir:  appPath,
			name: "config.yaml",
			fun:  appGen.CfgFile,
		},
		{
			dir:  appPath,
			name: "Dockerfile",
			fun:  appGen.Dockerfile,
		},
		{
			dir:  tmpl.ConfDir,
			name: "config.go",
			fun:  appGen.Conf,
		},
		{
			dir:  tmpl.ConfDir,
			name: "config_testing.go",
			fun:  appGen.ConfTesting,
		},
		{
			dir:  tmpl.ConfDir,
			name: "custom_init.go",
			fun:  appGen.ConfCustomInit,
		},
		{
			dir:  tmpl.ConfDir,
			name: "rpc_register.go",
			fun:  appGen.ConfRpcRegister,
		},
		{
			dir:  appPath,
			name: "main.go",
			fun:  appGen.Main,
		},
		{
			dir:  tmpl.ServiceDir,
			name: "di_register_gen.go",
			fun:  tpls.NewDI("service").Execute,
		},
		{
			dir:  tmpl.ExtendDir,
			name: "di_register_gen.go",
			fun:  tpls.NewDI("extend").Execute,
		},
		{
			dir:  tmpl.RepoImplDir,
			name: "di_register_gen.go",
			fun:  tpls.NewDI("repo_impl").Execute,
		},
		{
			dir:  tmpl.ConfDir,
			name: "di_register_gen.go",
			fun: tpls.NewSDI("conf", parser.XST{
				Name:      "Config",
				ShortName: "c",
				MPoint:    true,
				FieldList: map[string]parser.XField{
					"Etcd": {
						Name:  "Etcd",
						Type:  "*cfg.Etcd",
						SType: parser.STypeStruct,
						Idx:   0,
						Tag:   `sdi:"Etcd"`,
					},
					"DB": {
						Name:  "DB",
						Type:  "*cfg.DB",
						SType: parser.STypeStruct,
						Idx:   0,
						Tag:   `sdi:"DB"`,
					},
					"Redis": {
						Name:  "Redis",
						Type:  "*cfg.Redis",
						SType: parser.STypeStruct,
						Idx:   0,
						Tag:   `sdi:"Redis"`,
					},
				},
			}).Execute,
		},
		{
			dir:  tmpl.EventDir,
			name: "register_gen.go",
			fun:  tpls.NewEventRegister().Execute,
		},
		{
			dir:  tmpl.EventDir,
			name: "custom_event.go",
			fun: func() ([]byte, error) {
				return []byte(tpls.CustomEventRegisterTpl), nil
			},
		},
	}
	for _, f := range initFiles {
		_ = os.MkdirAll(f.dir, 0766)
		filename := fmt.Sprintf("%s/%s", f.dir, f.name)
		buf, err := f.fun()
		if err != nil {
			log.Printf("app gen [%s] err: %v \n", filename, err)
			continue
		}
		if strings.HasSuffix(f.name, ".go") {
			buf = m.format(buf, filename)
		}
		err = tool_file.WriteFile(filename, buf)
		if err != nil {
			log.Printf("app gen [%s] write file err: %v \n", filename, err)
		}
	}
	// gen handler
	_ = m.Handler()

	return nil
}
