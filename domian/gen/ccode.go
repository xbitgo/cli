package gen

import (
	"fmt"
	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"
	"log"
	"strings"
	"xbit/conf"
	"xbit/domian/gen/tpls"
	"xbit/domian/parser"
)

func (m *Manager) CRepo() error {
	ipr, err := parser.Scan(m.Tmpl.EntityDir, parser.ParseTypeWatch)
	if err != nil {
		log.Fatalf("CRepo: parse dir[%s], err: %v", m.Tmpl.EntityDir, err)
	}
	entityList := make([]string, 0)
	for _, it := range ipr.StructList {
		for _, field := range it.FieldList {
			tag := field.GetTag("db")
			if tag != nil {
				entityList = append(entityList, it.Name)
				break
			}
		}
	}
	for _, s := range entityList {
		tpl := tpls.Repo{
			ProjectName: conf.Global.ProjectName,
			AppName:     m.AppName,
			EntityName:  s,
		}
		buf, err := tpl.Execute()
		if err != nil {
			log.Printf("gen Repo %s err: %v \n", s, err)
			return err
		}
		filename := fmt.Sprintf("%s/%s_repo.go", m.Tmpl.RepoDir, tool_str.ToSnakeCase(s))
		buf = m.format(buf, filename)
		log.Printf("gen IMPL file %s \n", filename)
		err = tool_file.WriteFile(filename, buf)
		if err != nil {
			return err
		}

		buf, err = tpl.ExecuteImpl()
		if err != nil {
			log.Printf("gen Repo.Impl %s err: %v \n", s, err)
			return err
		}
		filename = fmt.Sprintf("%s/%s_repo_impl.go", m.Tmpl.RepoImplDir, tool_str.ToSnakeCase(s))
		buf = m.format(buf, filename)
		log.Printf("gen IMPL file %s \n", filename)
		err = tool_file.WriteFile(filename, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) CService() error {
	ipr, err := parser.Scan(m.Tmpl.EntityDir, parser.ParseTypeWatch)
	if err != nil {
		log.Fatalf("CRepo: parse dir[%s], err: %v", m.Tmpl.EntityDir, err)
	}
	entityList := make([]string, 0)
	for _, it := range ipr.StructList {
		for _, field := range it.FieldList {
			tag := field.GetTag("db")
			if tag != nil {
				entityList = append(entityList, it.Name)
				break
			}
		}
	}
	entityServiceMap := map[string][]tpls.ServiceLayerItem{}
	for _, s := range entityList {
		tmp := tool_str.ToSnakeCase(s)
		service := tool_str.ToUFirst(strings.Split(tmp, "_")[0])
		if _, ok := entityServiceMap[service]; !ok {
			entityServiceMap[service] = []tpls.ServiceLayerItem{}
		}
		entityServiceMap[service] = append(entityServiceMap[service], tpls.ServiceLayerItem{
			EntityName: s,
			VarName:    tool_str.ToLFirst(s),
		})
	}

	for s, items := range entityServiceMap {
		tpl := tpls.ServiceLayer{
			ProjectName: conf.Global.ProjectName,
			AppName:     m.AppName,
			ServiceName: s,
			EntityList:  items,
		}
		buf, err := tpl.Execute()
		if err != nil {
			log.Printf("gen Repo %s err: %v \n", s, err)
			return err
		}
		filename := fmt.Sprintf("%s/%s.go", m.Tmpl.ServiceDir, tool_str.ToSnakeCase(s))
		buf = m.format(buf, filename)
		log.Printf("gen IMPL file %s \n", filename)
		err = tool_file.WriteFile(filename, buf)
		if err != nil {
			return err
		}
	}
	return nil
}
