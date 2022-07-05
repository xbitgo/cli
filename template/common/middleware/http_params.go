package middleware

import (
	"{{ProjectName}}/common/ecode"
	"{{ProjectName}}/common/http_io"
	"github.com/gin-gonic/gin"
)

func HTTPParams(parser ...func(c *gin.Context) ([]byte, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		_parser := defaultParser
		if len(parser) > 0 {
			_parser = parser[0]
		}
		raw, err := _parser(c)
		if err != nil {
			http_io.JSONError(c, ecode.ErrParams)
		}
		c.Set(http_io.HTTPBodyKEy, raw)
		c.Next()
	}
}

func defaultParser(c *gin.Context) ([]byte, error) {
	raw, err := c.GetRawData()
	if err != nil {
		return nil, err
	}
	return raw, nil
}
