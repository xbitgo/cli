package tpls

import (
	"bytes"
	"text/template"
)

const testingHandlerImpl = `package handler

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/xbitgo/core/tools/tool_json"
	"github.com/xbitgo/core/tools/tool_mock"

	{{- range .OtherPkgList }}
	"{{.}}"
	{{- end}}
)

{{ range .Methods }}
func Test{{.Service}}HandlerImpl_{{.Name}}(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		// "something": "something_val",
	}))
	req := &pb.{{.Request}}{}
	tool_mock.NewDataMocker().Struct(req)
	_req, _ := tool_json.JSON.MarshalToString(req)
	t.Logf("req: %s", _req)
	rs, err := _test{{.Service}}Handler.{{.Name}}(ctx, req)
	if err != nil {
		t.Fatalf("err %+v", err)
	}
	_rs, _ := tool_json.JSON.MarshalToString(rs)
	t.Logf("resp: %s", _rs)
	select {
	case <-time.After(3 * time.Second):
		return
	}
}
{{ end}}`

const testingMainHandlerImpl = `package handler

import (
	"testing"

	"github.com/xbitgo/core/di"
	"github.com/xbitgo/components/eventbus"

	"{{.ProjectName}}/apps/{{.ServiceName}}/conf"
	"{{.ProjectName}}/apps/{{.ServiceName}}/domain/event"
	"{{.ProjectName}}/apps/{{.ServiceName}}/domain/service"
	"{{.ProjectName}}/apps/{{.ServiceName}}/repo_impl"
)
{{- range .Handlers}}
var _test{{.}}Handler *{{.}}HandlerImpl
{{- end}}

func TestMain(m *testing.M) {
	// 配置初始化
	conf.UnitTestingInit()
	// 自定义初始化配置
	conf.CustomInit()
	// 注册DI存储层实现
	repo_impl.DIRegister()
	// 注册DI服务层
	service.DIRegister()
	// DI注入
	di.MustBindALL()
	// eventbus初始化
	eventbus.Init(16)
	// event注册
	event.Register()
	// handler层
	{{- range .Handlers}}
	_test{{.}}Handler = New{{.}}HandlerImpl()
	di.MustBind(_test{{.}}Handler)
	{{- end}}
	// 注册rpc客户端
	conf.RegisterRPCClients()
	m.Run()
}`

// TestingHandler .
type TestingHandler struct {
	ProjectName string
	ServiceName string
	Handler
}

func (t *TestingHandler) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("testing.handler").Parse(testingHandlerImpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, t); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TestingMainHandler .
type TestingMainHandler struct {
	ProjectName string
	ServiceName string
	Handlers    []string
}

func (t *TestingMainHandler) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("testing.main.handler").Parse(testingMainHandlerImpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, t); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
