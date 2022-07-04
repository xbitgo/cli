package http_io

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	"github.com/xbitgo/core/tools/tool_json"

	"{{ProjectName}}/common/ecode"
)

const HTTPBodyKEy = "_BODY_"

type HTTPResponse struct {
	Code  int32       `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

func BindBody(ctx *gin.Context, obj interface{}) error {
	raw, ok := ctx.Get(HTTPBodyKEy)
	if !ok {
		return errors.New("not found params")
	}
	body, ok := raw.([]byte)
	if !ok {
		return errors.New("error params")
	}
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := tool_json.JSON.NewDecoder(r)
	if binding.EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if binding.EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj interface{}) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}

func Metadata(ctx *gin.Context) context.Context {
	return metadata.NewIncomingContext(ctx, map[string][]string(ctx.Request.Header))
}

// JSONSuccess 成功返回
func JSONSuccess(ctx *gin.Context, data interface{}) {
	res := HTTPResponse{
		Data: data,
	}
	toResp(ctx, res)
}

// JSONError 失败返回
func JSONError(ctx *gin.Context, err ecode.Error) {
	res := HTTPResponse{
		Code:  err.Code,
		Error: err.Err,
		Data:  struct{}{},
	}
	toResp(ctx, res)
}

// JSONResp .
func toResp(ctx *gin.Context, res interface{}) {
	ctx.Render(http.StatusOK, render.JSON{Data: res})
	ctx.Abort()
	return
}

func JSON(ctx *gin.Context, res interface{}, err error) {
	// resp
	if err != nil {
		ec := ecode.ErrSystem
		if v, ok := err.(ecode.Error); ok {
			ec = v
		}
		JSONError(ctx, ec)
		return
	}
	JSONSuccess(ctx, res)
}
