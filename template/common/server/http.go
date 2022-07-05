package server

import (
	"fmt"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

type httpSvc struct {
	Engine *gin.Engine
	Addr   string
}

func (s *httpSvc) Type() string {
	return "http"
}

func (s *httpSvc) StartAddr() string {
	addr, err := net.ResolveTCPAddr("", s.Addr)
	if err != nil {
		panic(err)
	}
	return ExposedAddr(addr).String()
}

func (s *httpSvc) Start() error {
	fmt.Printf(" -- starting http server: [%s] ... \n", s.Addr)
	return s.Engine.Run(s.Addr)
}

func (s *httpSvc) Stop() {

}

func (a *App) InitHTTP(process func(r *gin.Engine)) {
	addr := a.cfg.HTTPAddr
	if addr == "" {
		return
	}
	gin.DefaultWriter = os.Stderr
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	process(engine)
	a.svcList = append(a.svcList, &httpSvc{
		Engine: engine,
		Addr:   addr,
	})
}
