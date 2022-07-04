package gen

import (
	"bytes"
	"fmt"
	"github.com/emicklei/proto"
	"github.com/xbitgo/core/tools/tool_file"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"xbit/conf"
	"xbit/domian/gen/tpls"
)

func (m *Manager) Client(pbFile string) {
	var reg, _ = regexp.Compile(`service\s+(\w+)\s+{`)
	raw, err := ioutil.ReadFile(pbFile)
	if err != nil {
		return
	}
	rs := reg.FindAll(raw, -1)
	if len(rs) <= 0 {
		return
	}
	reader := bytes.NewReader(raw)
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	var (
		pkg string
		res []*tpls.Service
	)
	proto.Walk(definition,
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				pkg = strings.Split(o.Constant.Source, ";")[0]
			}
		}),
		proto.WithService(func(s *proto.Service) {
			tmp := strings.Split(pkg, "/")
			cs := &tpls.Service{
				ProjectName: conf.Global.ProjectName,
				Package:     pkg,
				PName:       tmp[len(tmp)-1],
				Service:     s.Name,
				RawName:     strings.TrimSuffix(s.Name, "Service"),
				Filepath:    pbFile,
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if ok {
					method := &tpls.Method{
						Service: s.Name,
						Name:    r.Name,
						PName:   cs.PName,
						RawName: cs.RawName,
						Request: r.RequestType,
						Reply:   r.ReturnsType,
					}
					if tmp := strings.Split(method.Request, "."); len(tmp) == 1 {
						method.Request = method.PName + "." + method.Request
					}
					if tmp := strings.Split(method.Reply, "."); len(tmp) == 1 {
						method.Reply = method.PName + "." + method.Reply
					}
					cs.Methods = append(cs.Methods, method)
				}
			}
			res = append(res, cs)
		}),
	)

	targetDir := m.Project.RootPath() + "/proto/rpc_client"
	for _, s := range res {
		buf, err := s.Execute()
		if err != nil {
			log.Fatal(err)
		}
		filename := fmt.Sprintf("%s/client_%s_%s", targetDir, strings.ToLower(s.RawName), "gen.go")
		m.format(buf, filename)
		err = tool_file.WriteFile(filename, buf)
		if err != nil {
			log.Printf("app gen [%s] write file err: %v \n", filename, err)
		}
	}
	return
}
