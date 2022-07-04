package gen

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/xbitgo/core/tools/tool_file"

	"xbit/domian/gen/tpls"
	"xbit/domian/parser"
)

func (m *Manager) Impl(dirs []string) error {
	infMap := map[string]parser.INF{}
	iprList := make([]*parser.IParser, 0)
	for _, dir := range dirs {
		ipr, err := parser.Scan(dir, parser.ParseTypeImpl)
		if err != nil {
			log.Fatalf("impl: parse dir[%s], err: %v", dir, err)
		}
		pkgName := ""
		if strings.HasPrefix(dir, m.Tmpl.DomainDir) {
			tmp := strings.Replace(dir, m.Tmpl.DomainDir+"/", "", 1)
			pkgName = strings.Replace(tmp, "/", ".", -1)
		}
		for _, inf := range ipr.INFList {
			infMap[pkgName+"."+inf.Name] = inf
		}
		if len(ipr.StructList) > 0 {
			iprList = append(iprList, ipr)
		}
	}

	for _, iPr := range iprList {
		for _, xst := range iPr.StructList {
			if xst.ImplINF != "" {
				if inf, ok := infMap[xst.ImplINF]; ok {
					_ = m.writeImplFunc(xst, inf)
				}
			}
		}
	}
	return nil
}

func (m *Manager) genImpl(implPr *parser.IParser, inf parser.INF) {
	for _, xst := range implPr.StructList {
		if xst.ImplINF == inf.Name {
			_ = m.writeImplFunc(xst, inf)
		}
	}
}

func (m *Manager) rewriteSameFunc(body []byte, xName string, changeFunc []parser.XMethod) []byte {
	for _, method := range changeFunc {
		rex := regexp.MustCompile(fmt.Sprintf(`%s\)\s*%s\(`, xName, method.Name) + `.*{`)
		gf := tpls.IMPLFunc{XName: xName, XMethod: method}
		buf, _ := gf.Execute()
		body = rex.ReplaceAll(body, buf)
	}
	return body
}

func (m *Manager) writeImplFunc(xst parser.XST, inf parser.INF) error {
	ok, noImplFuncs, changeFunc := xst.IsImpl(inf)
	if ok {
		return nil
	}
	body, _ := os.ReadFile(xst.File)

	// 判断是否有实现同名方法 替换
	if len(changeFunc) > 0 {
		body = m.rewriteSameFunc(body, xst.Name, changeFunc)
	}
	name := xst.Name
	if xst.MPoint {
		name = "*" + name
	}
	infTpl := tpls.INF{
		Body:       body,
		Name:       name,
		MethodList: noImplFuncs,
	}
	buf, err := infTpl.Execute()
	if err != nil {
		log.Printf("gen IMPL file[%s] err: %s \n", xst.File, err)
		return err
	}
	buf = m.format(buf, xst.File)
	log.Printf("gen IMPL file %s \n", xst.File)
	err = tool_file.WriteFile(xst.File, buf)
	if err != nil {
		return err
	}
	return nil
}
