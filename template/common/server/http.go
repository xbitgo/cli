package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

type httpSvc struct {
	Engine *gin.Engine
	Addr   string
}

func (s *httpSvc) Type() string {
	return "http"
}

func (s *httpSvc) Start() error {
	fmt.Printf(" -- starting http server: [%s] ... \n", s.Addr)
	return s.Engine.Run(s.Addr)
}

func (a *App) InitHTTP(addr string, process func(r *gin.Engine)) {
	if addr == "" {
		return
	}
	gin.DefaultWriter = os.Stderr
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	process(engine)
	a.httpSvc = &httpSvc{
		Engine: engine,
		Addr:   addr,
	}
}
