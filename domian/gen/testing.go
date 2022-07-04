package gen

import (
	"fmt"
	"log"
	"xbit/conf"
	"xbit/domian/parser"

	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"

	"xbit/domian/gen/tpls"
)

func (m *Manager) TestHandlers() error {
	// pb gen
	pbFiles := m.pbServices()
	tplHandlerTestingList := []tpls.TestingHandler{}
	for _, file := range pbFiles {
		infList := m.parsePbFile(file)
		for _, inf := range infList {
			tplHandlerTesting, err := m._testHandler(inf)
			if err != nil {
				return err
			}
			tplHandlerTestingList = append(tplHandlerTestingList, tplHandlerTesting)
		}
	}

	// testMain
	if len(tplHandlerTestingList) > 0 {
		tM := tpls.TestingMainHandler{
			ProjectName: conf.Global.ProjectName,
			ServiceName: tplHandlerTestingList[0].ServiceName,
			Handlers:    []string{},
		}
		for _, handler := range tplHandlerTestingList {
			tM.Handlers = append(tM.Handlers, handler.Service)
		}
		buf, err := tM.Execute()
		if err != nil {
			fmt.Println(err)
			return err
		}
		filename := fmt.Sprintf("%s/main_test.go", m.Tmpl.HandlerDir)
		buf = m.format(buf, filename)
		err = tool_file.WriteFile(filename, buf)
		if err != nil {
			log.Printf("app gen [%s] write file err: %v \n", filename, err)
		}
		return nil
	}
	return nil
}

func (m Manager) _testHandler(inf parser.INF) (tpls.TestingHandler, error) {
	tplHandler := m.toHandlerGen(inf)
	tplHandlerTesting := tpls.TestingHandler{
		ProjectName: conf.Global.ProjectName,
		ServiceName: tool_str.ToSnakeCase(tplHandler.Service),
		Handler:     tplHandler,
	}
	buf, err := tplHandlerTesting.Execute()
	if err != nil {
		fmt.Println(err)
		return tplHandlerTesting, err
	}
	filename := fmt.Sprintf("%s/%s_handler_test.go", m.Tmpl.HandlerDir, tool_str.ToSnakeCase(tplHandler.Service))
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("app gen [%s] write file err: %v \n", filename, err)
	}
	return tplHandlerTesting, nil
}
