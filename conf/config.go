package conf

import (
	"github.com/xbitgo/core/config"
	"github.com/xbitgo/core/di"
	"github.com/xbitgo/core/log"
)

type Config struct {
	ProjectName string            `yaml:"project_name"`
	DBList      map[string]string `yaml:"database"`
	Tmpl        *Tmpl             `yaml:"tmpl"`
}

var (
	Global      = &Config{}
	defaultTmpl = &Tmpl{
		AppDir:          "apps",
		ConfDir:         "{appPath}/conf",
		DomainDir:       "{appPath}/domain",
		EntityDir:       "{appPath}/domain/entity",
		DoDir:           "{appPath}/repo_impl/do",
		ConvDoDir:       "{appPath}/repo_impl/converter",
		DaoDir:          "{appPath}/repo_impl/dao",
		PbDir:           "{rootPath}/proto/apps/{appName}",
		ConvPbDir:       "{appPath}/handler/converter",
		ServiceDir:      "{appPath}/domain/service",
		ExtendDir:       "{appPath}/domain/extend",
		RepoDir:         "{appPath}/domain/repo",
		RepoImplDir:     "{appPath}/repo_impl",
		EventDir:        "{appPath}/domain/event",
		HandlerDir:      "{appPath}/handler",
		HandlerEntryDir: "{appPath}/handler/entry",
		SQLDir:          "{appPath}/repo_impl/sql",
	}
)

func Init() {
	di.PrintLog = false
	cfg := config.Yaml{ConfigFile: "xbit.yaml"}
	err := cfg.Apply(Global)
	if err != nil {
		log.Panic(err)
	}
	if Global.Tmpl == nil {
		Global.Tmpl = defaultTmpl
	}
}
