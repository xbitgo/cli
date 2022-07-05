package cfg

import (
	"github.com/xbitgo/core/config"
	//"github.com/xbitgo/components/cfg_adapter"
	//agolloCfg "github.com/apolloconfig/agollo/v4/env/config"
)

type Cfg struct {
	//*cfg_adapter.Apollo
	*config.Yaml
}

func NewCfg(arg string) *Cfg {
	// 默认使用yaml文件配置
	return &Cfg{Yaml: &config.Yaml{ConfigFile: "config.yaml"}}
	// 使用携程apollo配置中心配置
	//apoCfg, err := cfg_adapter.NewApollo(&agolloCfg.AppConfig{
	//	AppID:   "",
	//	Cluster: "",
	//	IP:      "",
	//	Secret:  "",
	//}, arg)
	//if err != nil {
	//	panic(err)
	//}
	//return &Cfg{Apollo: apoCfg}
}
