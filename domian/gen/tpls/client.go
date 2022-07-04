package tpls

import (
	"bytes"
	"text/template"
)

const clientTpl = `// Code generated by xbit. DO NOT EDIT.
package rpc_client

import (
	"context"
	"time"

	//"google.golang.org/grpc/metadata"
	"google.golang.org/grpc"
	"github.com/pkg/errors"

	"github.com/xbitgo/components/prometheus"

	"{{.ProjectName}}/proto/apps/{{.PName}}"
)

var {{.RawName}} = RPCService{Name: "{{.PName}}"}

func Get{{.RawName}}Client() *{{.RawName}}Client {
	conn, ok := RpcConn[{{.RawName}}.Name]
	if !ok {
		panic(errors.Errorf("No Register RPC Client[%s]", {{.RawName}}.Name))
	}
	return New{{.RawName}}Client(conn)
}

type {{.RawName}}Client struct {
	cli {{.PName}}.{{.Service}}Client
}

func New{{.RawName}}Client(cc *grpc.ClientConn) *{{.RawName}}Client {
	return &{{.RawName}}Client{
		cli: {{.PName}}.New{{.Service}}Client(cc),
	}
}

{{ range .Methods }}

func (s *{{.RawName}}Client) {{.Name}}(ctx context.Context, req *{{.Request}}, timeout ...time.Duration) (*{{.Reply}}, error) {
	_timeout := 2 * time.Second
	if len(timeout) > 0 {
		_timeout = timeout[0]
	}
	ctx, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()
	st := time.Now()
	resp, err := s.cli.{{.Name}}(ctx, req)
	prometheus.HistogramVec.Timing("rpc_client_{{.RawName}}_{{.Name}}", []string{"from",From,"ret", prometheus.RetLabel(err)}, st)
	return resp, err
}
{{- end }}
`

// Service is a proto service.
type Service struct {
	ProjectName string
	Package     string
	PName       string
	CName       string
	Service     string
	RawName     string
	Methods     []*Method
	Filepath    string
}

// Method is a proto method.
type Method struct {
	Service  string
	Name     string
	PName    string
	RawName  string
	Request  string
	Reply    string
	LogLevel string
}

func (s *Service) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("client").Parse(clientTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
