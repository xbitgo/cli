package tpls

import (
	"bytes"
	"text/template"
)

const AppConfigTestingTmpl = `package conf

import "{{.ProjectName}}/common/cfg"

// UnitTestingInit 单元测试时的配置初始化
func UnitTestingInit() {
	//Init() 一般配置中心模式可以直接调用正常初始化
	App = &Config{
		Name: "{{.ServiceName}}_server",
		Env:  "local",
		Server: &cfg.Server{
			HTTPAddr: "0.0.0.0:9901",
			GRPCAddr: "0.0.0.0:9902",
		},
		Tracing: &cfg.Tracing{
			ServiceName:                "{{.ServiceName}}_server",
			SamplerType:                "const",
			SamplerParam:               1,
			ReporterLocalAgentHostPort: "127.0.0.1:6831",
			LogSpans:                   false,
		},
		Etcd: &cfg.Etcd{
			Endpoints: []string{"127.0.0.1:2379"},
			Timeout:   5,
		},
		DB: &cfg.DB{
			DSN:             "root:@tcp(localhost:3306)/test?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
			MaxOpenConn:     10,
			MaxIdleConn:     10,
			ConnMaxLifetime: 300,
			ConnMaxIdleTime: 300,
		},
		Redis: &cfg.Redis{
			Addr:     "127.0.0.1:6379",
			Password: "",
			PoolSize: 4,
			DB:       0,
		},
	}
	_dIRegister()
}`

const AppConfigCustomInitImpl = `package conf

import (
	"os"

	"github.com/xbitgo/core/log"
	"github.com/xbitgo/components/tracing"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go/extra"
)

// CustomInit 业务自定义配置
func CustomInit() {
	// JSON配置
	binding.EnableDecoderUseNumber = true
	extra.RegisterFuzzyDecoders()

	// 链路追踪配置
	tracing.InitJaegerTracer(App.Tracing.ToJaegerCfg())

	// log 配置
	log.InitLogger(os.Stderr)
	log.SetLevel(log.DebugLevel)
	log.SetTraceIdFunc(tracing.TraceID)
}`

const AppConfigTmplRpcRegisterImpl = `package conf

// RegisterRPCClients 注册Rpc客户端
func RegisterRPCClients() {
	//rpc_client.RegisterRPC(conf.Namespace,clients.User) // 在这里注册依赖的其他服务 注册后 代码中可以直接使用
}`

const AppConfigTmpl = `package conf

import (
    "{{.ProjectName}}/common/cfg"
)

const Namespace = "{{.ServiceName}}"

var App = &Config{}

type Config struct {
    Name      string      ` + "`" + `json:"name" yaml:"name"` + "`" + `
    Env       string      ` + "`" + `json:"env" yaml:"env"` + "`" + `
    Server    *cfg.Server ` + "`" + `json:"server" yaml:"server"` + "`" + `
	Tracing   *cfg.Tracing ` + "`" + `json:"trace" yaml:"tracing"` + "`" + `
    Etcd      *cfg.Etcd   ` + "`" + `json:"etcd" yaml:"etcd" sdi:"etcd""` + "`" + `
    DB        *cfg.DB     ` + "`" + `json:"DB" yaml:"DB" sdi:"DB"` + "`" + `
    Redis     *cfg.Redis ` + "`" + `json:"redis" yaml:"redis" sdi:"redis"` + "`" + `
}

func Init() {
    c := cfg.NewCfg(Namespace)
    err := c.Apply(App)
    if err != nil {
        log.Panic(err)
    }
    _dIRegister()
}`

const AppDockerTmpl = `FROM frolvlad/alpine-glibc:glibc-2.34

COPY ./start /opt/start

ENV SERVICE_NAME="{{.ServiceName}}_server"

WORKDIR /opt

CMD ["./start"]
// todo 根据实际情况配置
`

const ApiTestHttpTmpl = `POST http://localhost:9901/api/{{.ServiceName}}/test
Content-Type: application/json

{}

###`

const AppMainTmpl = `package main

import (
    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"

	"github.com/xbitgo/components/tracing"
	"github.com/xbitgo/core/di"

	"github.com/xbitgo/components/dtx"
	"github.com/xbitgo/components/sequence"

	"{{.ProjectName}}/common/middleware"
	"{{.ProjectName}}/common/server"
	

    pb "{{.ProjectName}}/proto/{{.AppsDir}}/{{.ServiceName}}"
    "{{.ProjectName}}/apps/{{.ServiceName}}/conf"
	"{{.ProjectName}}/apps/{{.ServiceName}}/domain/event"
	"{{.ProjectName}}/apps/{{.ServiceName}}/domain/extend"
    "{{.ProjectName}}/apps/{{.ServiceName}}/domain/service"
    "{{.ProjectName}}/apps/{{.ServiceName}}/handler/entry"
    "{{.ProjectName}}/apps/{{.ServiceName}}/repo_impl"
)

func main() {
    // 初始化
	Init()
	// handler层代理
	handlerProxy := entry.New{{.ServiceNameUF}}Handler()
    // 创建服务
    app := server.NewApp()
    // grpc server
    app.InitGRPC(conf.App.Server.GRPCAddr, func(grpcService grpc.ServiceRegistrar) {
        pb.Register{{.ServiceNameUF}}Server(grpcService, handlerProxy)
    }, tracing.GrpcServerTrace())
    // http server
    app.InitHTTP(conf.App.Server.HTTPAddr, func(r *gin.Engine) {
        // 可选 全局中间件
        // r.Use(middleware.HTTPCors())
        // 可选 全局路由
        // r.OPTIONS("/*wild", func(c *gin.Context) {
        //	return
        // })
        // 必须；解析参数可以定制
        r.Use(middleware.HTTPParams())
        {
            entry.{{.ServiceNameUF}}HttpInit(r, handlerProxy)
        }
    })
    // 服务停止前处理
    app.OnClose(func() {
        // something close ...
    })
    err := app.Start()
    if err != nil {
        panic(err)
    }
}

func Init()  {
	// 配置初始化
	conf.Init()
	// 自定义初始化配置
	conf.CustomInit()
	// 分布式ID生成器初始化
	sequence.Init()
	// 分布式事务管理器初始化
	dtx.Init(nil, 0) // 需要配置mq开启跨服务
	// 注册rpc客户端
	conf.RegisterRPCClients()
	// 注册DI存储层实现
	repo_impl.DIRegister()
	// 注册DI业务拓展层
	extend.DIRegister()
	// 注册DI服务层
	service.DIRegister()
	// event注册
	event.Register()
	// DI注入
	di.MustBindALL()
}
`

const AppCfgFileTmpl = `name: {{.ServiceName}}
env: local
server:
  http_addr: "0.0.0.0:9901"
  grpc_addr: "0.0.0.0:9902"
tracing:
  service_name: "{{.ServiceName}}"
  sampler_type: "const"
  sampler_param: 1
  reporter_local_agent_host_port: "localhost:6831"
  log_spans: false

etcd:
  endpoints:
    - 127.0.0.1:2379
  timeout: 5

DB:
  type: mysql
  dsn: "root:@tcp(localhost:3306)/test?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local"
  max_open_conn: 10
  max_idle_conn: 10
  conn_max_lifetime: 300
  conn_max_idle_time: 300

redis:
  addr: "127.0.0.1:6379"
  password: ""
  pool_size: 4
  db: 0`

const AppProtoTmpl = `syntax = "proto3";

package {{.ServiceName}};

import "google/api/annotations.proto";
//import "google/api/gogo.proto";
//import "{{.AppsDir}}/{{.ServiceName}}/{{.ServiceName}}_message_gen.proto";

option go_package = "proto/{{.AppsDir}}/{{.ServiceName}};{{.ServiceName}}";

service {{.ServiceNameUF}} {
    // Test
    rpc Test (TestRequest) returns (TestResponse)  {
        option (google.api.http) = {
            post: "/api/{{.ServiceName}}/test"
        };
    }
}

message TestRequest {
}

message TestResponse {
    string msg = 1;
}`

type App struct {
	ProjectName   string
	ServiceName   string
	ServiceNameUF string
	AppsDir       string
}

func (s *App) Conf() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.Conf").Parse(AppConfigTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) ConfTesting() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.ConfTesting").Parse(AppConfigTestingTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) ConfCustomInit() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.ConfCustomInit").Parse(AppConfigCustomInitImpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) ConfRpcRegister() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.ConfRpcRegister").Parse(AppConfigTmplRpcRegisterImpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) Dockerfile() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.Dockerfile").Parse(AppDockerTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) ApiTestHttp() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.ApiTestHttp").Parse(ApiTestHttpTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) Main() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.Main").Parse(AppMainTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) CfgFile() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.CfgFile").Parse(AppCfgFileTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *App) Proto() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("app.Proto").Parse(AppProtoTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
