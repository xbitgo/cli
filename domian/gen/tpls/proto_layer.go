package tpls

import (
	"bytes"
	"text/template"
)

const protoLayerTpl = `
syntax = "proto3";

package {{.AppName}};

import "google/api/annotations.proto";
import "google/api/gogo.proto";
import "apps/{{.AppName}}/{{.AppName}}_message_gen.proto";
import "common/base.proto";

option go_package = "{{.ProjectName}}/proto/apps/{{.AppName}};{{.AppName}}";

service {{.ServiceName}}Svc {
{{- range .EntityList }}
  // {{.EntityName}}List 
  rpc {{.EntityName}}List ({{.EntityName}}ListRequest) returns ({{.EntityName}}ListResponse)  {
    option (google.api.http) = {
      post: "/api/{{$.AppName}}/{{.EntityName}}List"
    };
  }

  // Create{{.EntityName}}
  rpc Create{{.EntityName}} (Create{{.EntityName}}Request) returns (Create{{.EntityName}}Response)  {
    option (google.api.http) = {
      post: "/api/{{$.AppName}}/Create{{.EntityName}}"
    };
  }

  // Update{{.EntityName}} 
  rpc Update{{.EntityName}} (Update{{.EntityName}}Request) returns (Update{{.EntityName}}Response)  {
    option (google.api.http) = {
      post: "/api/{{$.AppName}}/Update{{.EntityName}}"
    };
  }

  // Delete{{.EntityName}} 
  rpc Delete{{.EntityName}} (Delete{{.EntityName}}Request) returns (Delete{{.EntityName}}Response)  {
    option (google.api.http) = {
      post: "/api/{{$.AppName}}/Delete{{.EntityName}}"
    };
  }
{{- end}}
}

{{- range .EntityList }}
message {{.EntityName}}ListRequest {
  base.Page page = 1 [json_name = "page", (gogoproto.moretags) = 'validate:"required" label:"分页"'];
  base.OrderBy orderBy = 2 [json_name = "order_by", (gogoproto.moretags) = 'validate:"required" label:"排序"'];
  repeated base.Filtering filters = 3 [json_name = "filters", (gogoproto.moretags) = 'validate:"required" label:"筛选项"'];
}

message {{.EntityName}}ListResponse {
  base.PageInfo PageInfo = 1 [json_name = "page_info", (gogoproto.moretags) = 'label:"分页信息"'];
  repeated {{.EntityName}} list = 2 [json_name = "list", (gogoproto.moretags) = 'label:"数据列表"'];
}

message Create{{.EntityName}}Request {
  {{.EntityName}} {{.VarName}} = 1 [json_name = "{{.VarName}}", (gogoproto.moretags) = 'validate:"required" label:"数据"'];
}

message Create{{.EntityName}}Response {
  {{.EntityName}} {{.VarName}} = 1 [json_name = "{{.VarName}}", (gogoproto.moretags) = 'validate:"required" label:"操作完成数据"'];
}

message Update{{.EntityName}}Request {
  int64 {{.VarName}}Id = 1 [json_name = "{{.VarName}}Id", (gogoproto.moretags) = 'validate:"required" label:"ID"'];
  {{.EntityName}} {{.VarName}} = 2 [json_name = "{{.VarName}}", (gogoproto.moretags) = 'validate:"required" label:"数据"'];
  repeated string fields = 3 [json_name = "fields", (gogoproto.moretags) = 'validate:"required" label:"更新字段"'];	
}

message Update{{.EntityName}}Response {
  {{.EntityName}} {{.VarName}} = 1 [json_name = "{{.VarName}}", (gogoproto.moretags) = 'validate:"required" label:"操作完成数据"'];
}

message Delete{{.EntityName}}Request {
  int64 {{.VarName}}Id = 1 [json_name = "{{.VarName}}Id", (gogoproto.moretags) = 'validate:"required" label:"ID"'];
}

message Delete{{.EntityName}}Response {}

{{- end}}
`

const protoLayerImplTpl = `
package handler

import (
	"context"
	"{{.ProjectName}}/apps/{{.AppName}}/domain/service"
	"{{.ProjectName}}/apps/{{.AppName}}/handler/converter"
	pb "{{.ProjectName}}/proto/apps/{{.AppName}}"

	"github.com/pkg/errors"
)

// {{.ServiceName}}SvcHandlerImpl @IMPL[{{.ServiceName}}Svc]
type {{.ServiceName}}SvcHandlerImpl struct {
	pb.Unimplemented{{.ServiceName}}SvcServer
	{{.ServiceName}}Service *service.{{.ServiceName}} ` + "`" + `di:"service.{{.ServiceName}}"` + "`" + `
}

func New{{.ServiceName}}SvcHandlerImpl() *{{.ServiceName}}SvcHandlerImpl {
	return &{{.ServiceName}}SvcHandlerImpl{}
}

{{- range .EntityList }}

func (impl *{{$.ServiceName}}SvcHandlerImpl) Create{{.EntityName}}(ctx context.Context, req *pb.Create{{.EntityName}}Request) (resp *pb.Create{{.EntityName}}Response, err error) {
	resp = &pb.Create{{.EntityName}}Response{}
	var (
		input = converter.To{{.EntityName}}Entity(req.{{.EntityName}})
	)
	if input == nil {
		return nil, errors.New("params error")
	}
	output, err := impl.{{$.ServiceName}}Service.Create{{.EntityName}}(ctx, input)
	if err != nil {
		return nil, err
	}
	resp.{{.EntityName}} = converter.From{{.EntityName}}Entity(output)
	return
}

func (impl *{{$.ServiceName}}SvcHandlerImpl) {{.EntityName}}List(ctx context.Context, req *pb.{{.EntityName}}ListRequest) (resp *pb.{{.EntityName}}ListResponse, err error) {
	resp = &pb.{{.EntityName}}ListResponse{}

	filterPage := converter.ToFilterPage(req.Page, req.OrderBy)
	filterList := converter.ToFilteringList(req.Filters)
	list, count, err := impl.{{$.ServiceName}}Service.Query{{.EntityName}}(ctx, filterList, filterPage)
	if err != nil {
		return nil, err
	}
	resp.List = converter.From{{.EntityName}}List(list)
	resp.PageInfo = converter.ToPageInfo(filterPage, count)
	return resp, nil
}

func (impl *{{$.ServiceName}}SvcHandlerImpl) Update{{.EntityName}}(ctx context.Context, req *pb.Update{{.EntityName}}Request) (resp *pb.Update{{.EntityName}}Response, err error) {
	resp = &pb.Update{{.EntityName}}Response{}
	if req.{{.EntityName}}Id == 0 {
		return nil, errors.New("params error")
	}
	{{.VarName}}, err := impl.{{$.ServiceName}}Service.Get{{.EntityName}}(ctx, req.{{.EntityName}}Id)
	if err != nil {
		return nil, errors.New("{{.VarName}} not found")
	}
	toModify := converter.To{{.EntityName}}Entity(req.{{.EntityName}})
	modifyMap := {{.VarName}}.ModifyDBMap(toModify, req.Fields)
	new{{.EntityName}}, err := impl.{{$.ServiceName}}Service.Set{{.EntityName}}(ctx, req.{{.EntityName}}Id, modifyMap)
	if err != nil {
		return nil, err
	}
	resp.{{.EntityName}} = converter.From{{.EntityName}}Entity(new{{.EntityName}})
	return resp, nil
}

func (impl *{{$.ServiceName}}SvcHandlerImpl) Delete{{.EntityName}}(ctx context.Context, req *pb.Delete{{.EntityName}}Request) (resp *pb.Delete{{.EntityName}}Response, err error) {
	resp = &pb.Delete{{.EntityName}}Response{}
	if req.{{.EntityName}}Id == 0 {
		return nil, errors.New("params error")
	}
	{{.VarName}}, err := impl.{{$.ServiceName}}Service.Get{{.EntityName}}(ctx, req.{{.EntityName}}Id)
	if err != nil {
		return nil, errors.New("{{.VarName}} not found")
	}
	err = impl.{{$.ServiceName}}Service.Delete{{.EntityName}}(ctx, {{.VarName}}.ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

{{- end}}
`

type ProtoLayer struct {
	ProjectName string
	AppName     string
	ServiceName string
	EntityList  []ProtoLayerItem
}

type ProtoLayerItem struct {
	EntityName string
	VarName    string
}

func (s *ProtoLayer) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("ProtoLayer").Parse(protoLayerTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *ProtoLayer) ExecuteImpl() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("ProtoLayer").Parse(protoLayerImplTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
