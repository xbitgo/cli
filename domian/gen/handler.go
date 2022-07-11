package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"xbit/conf"

	"github.com/emicklei/proto"

	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"

	"xbit/domian/gen/tpls"
	"xbit/domian/parser"
)

func (m *Manager) pbServices() []string {
	pbFiles := make([]string, 0)
	fileInfos, err := ioutil.ReadDir(m.Tmpl.PbDir)
	if err != nil {
		return pbFiles
	}
	var reg, _ = regexp.Compile(`service\s+(\w+)\s+{`)
	for _, fi := range fileInfos {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".proto") {
			pbFile := m.Tmpl.PbDir + "/" + fi.Name()
			raw, err := ioutil.ReadFile(pbFile)
			if err != nil {
				continue
			}
			rs := reg.FindAll(raw, -1)
			if len(rs) > 0 {
				pbFiles = append(pbFiles, pbFile)
			}
		}
	}
	return pbFiles
}

func (m *Manager) pbFiles() []string {
	pbFiles := make([]string, 0)
	fileInfos, err := ioutil.ReadDir(m.Tmpl.PbDir)
	if err != nil {
		return pbFiles
	}
	for _, fi := range fileInfos {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".proto") {
			pbFile := fmt.Sprintf("%s/%s/%s", m.Tmpl.AppDir, m.AppName, fi.Name())
			pbFiles = append(pbFiles, pbFile)
		}
	}
	return pbFiles
}

func (m *Manager) Handler() error {
	pbFiles := m.pbServices()
	for _, file := range pbFiles {
		err := m._handlerFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) parsePbFile(pbFile string) []parser.INF {
	reader, err := os.Open(pbFile)
	if err != nil {
		log.Fatalf("os.Open[%s] err: %v", pbFile, err)
	}
	defer reader.Close()

	pbParser := proto.NewParser(reader)
	definition, err := pbParser.Parse()
	if err != nil {
		log.Fatalf("pb Parse[%s] err: %v", pbFile, err)
	}
	var infList = make([]parser.INF, 0)
	proto.Walk(definition,
		proto.WithService(func(s *proto.Service) {
			inf := parser.INF{
				Name:    s.Name,
				Methods: map[string]parser.XMethod{},
			}
			for idx, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				httpMethod, httpRule := "", ""
				for _, option := range r.Options {
					if option.Name == "(google.api.http)" {
						for _, literal := range option.Constant.OrderedMap {
							if literal.Name != "" {
								httpMethod = literal.Name
								if literal.Name != "any" {
									httpMethod = strings.ToUpper(literal.Name)
								}
								if literal.Literal != nil {
									httpRule = literal.Literal.Source
								}
								break
							}
						}
					}
				}
				if ok {
					method := parser.XMethod{
						ImplName: "impl",
						Name:     r.Name,
						Params: []parser.XArg{
							{
								Name: "ctx",
								Type: "context.Context",
							},
							{
								Name: "req",
								Type: "*" + r.RequestType,
							},
						},
						Results: []parser.XArg{
							{
								Name: "resp",
								Type: "*" + r.ReturnsType,
							},
							{
								Name: "err",
								Type: "error",
							},
						},
						Comment:    r.Comment.Message(),
						Sort:       idx,
						HTTPMethod: httpMethod,
						HTTPRule:   httpRule,
					}
					inf.Methods[method.Name] = method
				}
			}
			infList = append(infList, inf)
		}),
	)
	return infList
}

func (m *Manager) _handlerFile(pbFile string) error {
	infList := m.parsePbFile(pbFile)
	for _, inf := range infList {
		err := m._handler(inf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) _handler(inf parser.INF) error {
	_ = m.handlerProxy(m.toHandlerGen(inf))
	_ = m.handlerHttp(m.toHandlerGen(inf))
	// 生成Impl 不存在
	_ = m.handlerImpl(m.toHandlerGen(inf))
	// 生成Impl 已存在
	hPr, err := parser.Scan(m.Tmpl.HandlerDir, parser.ParseTypeImpl)
	if err != nil {
		log.Fatalf("parser scan[%s] err: %v", m.Tmpl.HandlerDir, err)
	}
	withPb := false
	for _, xst := range hPr.StructList {
		if xst.ImplINF == inf.Name {
			for _, method := range xst.Methods {
				if len(method.Params) == 2 {
					if strings.Contains(method.Params[1].Type, "pb.") {
						withPb = true
					}
					break
				}
			}
		}
	}
	if withPb {
		methods := make(map[string]parser.XMethod)
		for s, mth := range inf.Methods {
			_mth := mth
			_mth.Params[1].Type = strings.Replace(_mth.Params[1].Type, "*", "*pb.", 1)
			_mth.Results[0].Type = strings.Replace(_mth.Results[0].Type, "*", "*pb.", 1)
			methods[s] = _mth
		}
		inf.Methods = methods
	}
	m.genImpl(hPr, inf)

	return nil
}

func (m *Manager) toHandlerGen(inf parser.INF) tpls.Handler {
	pkgPath := fmt.Sprintf("proto/%s/%s", m.Tmpl.AppDir, m.AppName)
	handler := tpls.Handler{
		OtherPkgList: m.otherPkgList(inf),
		Package:      pkgPath,
		Service:      inf.Name,
		AppName:      m.AppName,
		ProjectName:  conf.Global.ProjectName,
		Methods:      make([]*tpls.HMethod, 0),
		Filepath:     "",
	}

	methodList := make([]parser.XMethod, 0)
	for _, method := range inf.Methods {
		if method.HTTPMethod == "" || method.HTTPRule == "" {
			continue
		}
		methodList = append(methodList, method)
	}
	sort.SliceStable(methodList, func(i, j int) bool {
		return methodList[i].Sort < methodList[j].Sort
	})
	for _, method := range methodList {
		req := strings.Replace(method.Params[1].Type, "*pb.", "", 1)
		resp := strings.Replace(method.Results[0].Type, "*pb.", "", 1)
		req = strings.TrimPrefix(req, "*")
		resp = strings.TrimPrefix(resp, "*")
		hm := &tpls.HMethod{
			Service:    inf.Name,
			Name:       method.Name,
			AppName:    m.AppName,
			Request:    req,
			Reply:      resp,
			LogLevel:   "Debugf",
			Comment:    method.Comment,
			HTTPMethod: method.HTTPMethod,
			HTTPRule:   method.HTTPRule,
		}
		handler.Methods = append(handler.Methods, hm)
	}
	return handler
}

func (m *Manager) handlerProxy(handler tpls.Handler) error {
	buf, err := handler.ExecuteProxy()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/%s_handler_proxy_gen.go", m.Tmpl.HandlerEntryDir, tool_str.ToSnakeCase(handler.Service))
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("app gen [%s] write file err: %v \n", filename, err)
	}
	return nil
}

func (m *Manager) handlerImpl(handler tpls.Handler) error {
	filename := fmt.Sprintf("%s/%s_handler.go", m.Tmpl.HandlerDir, tool_str.ToSnakeCase(handler.Service))
	if tool_file.Exists(filename) {
		return nil
	}
	buf, err := handler.ExecuteIMPL()
	if err != nil {
		return err
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("app gen [%s] write file err: %v \n", filename, err)
	}
	return nil
}

func (m *Manager) handlerHttp(handler tpls.Handler) error {
	buf, err := handler.ExecuteHttp()
	if err != nil {
		fmt.Println(err)
		return err
	}
	filename := fmt.Sprintf("%s/%s_handler_http_gen.go", m.Tmpl.HandlerEntryDir, tool_str.ToSnakeCase(handler.Service))
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("app gen [%s] write file err: %v \n", filename, err)
	}
	return nil
}

func (m *Manager) otherPkgList(inf parser.INF) []string {
	pkgList := make([]string, 0)
	tmpMap := map[string]struct{}{}
	for _, method := range inf.Methods {
		if len(method.Params) == 2 {
			if strings.Contains(method.Params[1].Type, ".") {
				tmp := strings.Split(strings.TrimPrefix(method.Params[1].Type, "*"), ".")
				tmpMap[tmp[0]] = struct{}{}
			}
		}
		if len(method.Results) == 2 {
			if strings.Contains(method.Results[0].Type, ".") {
				tmp := strings.Split(strings.TrimPrefix(method.Results[0].Type, "*"), ".")
				tmpMap[tmp[0]] = struct{}{}
			}
		}
	}

	for s := range tmpMap {
		if s == "base" {
			pkgList = append(pkgList, "proto/base")
			continue
		}
		if s == "pb" {
			continue
		}
		for _, app := range m.Project.List() {
			if strings.HasSuffix(app, s) {
				pkgList = append(pkgList, fmt.Sprintf("proto/%s/%s", m.Tmpl.AppDir, app))
			}
		}
	}
	return pkgList
}
