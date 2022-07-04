package tpls

import (
	"bytes"
	"text/template"
)

const handlerTmpl = `package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"

{{- range .OtherPkgList }}
	"{{.}}"
{{- end}}
	pb "{{.ProjectName}}/{{.Package}}"
)

// {{.Service}}HandlerImpl @IMPL[{{.Service}}]
type {{.Service}}HandlerImpl struct {
	pb.Unimplemented{{.Service}}Server
}

func New{{.Service}}HandlerImpl() *{{.Service}}HandlerImpl {
	return &{{.Service}}HandlerImpl{}
}
{{ range .Methods }}

func (impl *{{.Service}}HandlerImpl) {{.Name}}(ctx context.Context, req *pb.{{.Request}}) (*pb.{{.Reply}}, error) {
	return &pb.{{.Reply}}{}, nil
}
{{ end }}`

const handlerProxyTmpl = `// Code generated by xbit. DO NOT EDIT.
package entry

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"

	"github.com/xbitgo/core/di"
	"github.com/xbitgo/core/log"

{{- range .OtherPkgList }}
	"{{.ProjectName}}/{{.}}"
{{- end}}
	. "{{.ProjectName}}/{{.Package}}"
	"{{.ProjectName}}/apps/{{.AppName}}/handler"
)

type {{.Service}}Handler struct {
	Unimplemented{{.Service}}Server
	impl {{.Service}}Server
}

func New{{.Service}}Handler() *{{.Service}}Handler {
	impl := handler.New{{.Service}}HandlerImpl()
	di.MustBind(impl)
	return &{{.Service}}Handler{
		impl: impl,
	}
}
{{ range .Methods }}
func (s *{{.Service}}Handler) {{.Name}}(ctx context.Context, req *{{.Request}}) (*{{.Reply}}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "{{.Service}}Handler_{{.Name}}")
	defer span.Finish()

	meta, _ :=  metadata.FromIncomingContext(ctx)
	log.With().TraceID(ctx).Field("request", req).Field("metadata",meta).{{.LogLevel}}("[{{.Name}}] on request")
	resp, err := s.impl.{{.Name}}(ctx, req)
	log.With().
		Field("request", req).
		Field("response", resp).
		Field("error", err).
		Field("metadata",meta).
		{{.LogLevel}}("[{{.Name}}] on response")
	return resp, err
}

{{ end }}`

const handlerHttpImpl = `// Code generated by xbit. DO NOT EDIT.
package entry

import (
	"github.com/gin-gonic/gin"

	pb "{{.ProjectName}}/{{.Package}}"

	"{{.ProjectName}}/common/ecode"
	"{{.ProjectName}}/common/http_io"
)

func {{.Service}}HttpInit(r *gin.Engine, impl pb.{{.Service}}Server) {
	{{ range .Methods }}
	// {{.Comment}}
	r.{{.HTTPMethod}}("{{.HTTPRule}}", func(ctx *gin.Context) {
		var req = &pb.{{.Request}}{}
		if err := http_io.BindBody(ctx, &req); err != nil {
			http_io.JSONError(ctx, ecode.ErrParams)
			return
		}
		res, err := impl.{{.Name}}(http_io.Metadata(ctx), req)
		http_io.JSON(ctx, res, err)
	})
	{{end}}
}`

// Handler .
type Handler struct {
	OtherPkgList []string
	Package      string
	Service      string
	AppName      string
	ProjectName  string
	Methods      []*HMethod
	Filepath     string
}

// HMethod is a proto method.
type HMethod struct {
	Service    string
	Name       string
	AppName    string
	Request    string
	Reply      string
	LogLevel   string
	Comment    string
	HTTPRule   string
	HTTPMethod string
}

func (s *Handler) ExecuteProxy() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("Handler.proxy").Parse(handlerProxyTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Handler) ExecuteIMPL() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("Handler.impl").Parse(handlerTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Handler) ExecuteHttp() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("Handler.http").Parse(handlerHttpImpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
