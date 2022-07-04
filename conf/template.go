package conf

import "strings"

type Tmpl struct {
	AppDir          string `json:"app_dir" yaml:"app_dir"`
	ConfDir         string `json:"conf_dir" yaml:"conf_dir"`
	DomainDir       string `json:"domain_dir" yaml:"domain_dir"`
	EntityDir       string `json:"entity_dir" yaml:"entity_dir"`
	DoDir           string `json:"do_dir" yaml:"do_dir"`
	ConvDoDir       string `json:"conv_do_dir" yaml:"conv_do_dir"`
	DaoDir          string `json:"dao_dir" yaml:"dao_dir"`
	PbDir           string `json:"pb_dir" yaml:"pb_dir"`
	ConvPbDir       string `json:"conv_pb_dir" yaml:"conv_pb_dir"`
	ServiceDir      string `json:"service_dir" yaml:"service_dir"`
	ExtendDir       string `json:"extend_dir" yaml:"extend_dir"`
	RepoDir         string `json:"repo_dir" yaml:"repo_dir"`
	RepoImplDir     string `json:"repo_impl_dir" yaml:"repo_impl_dir"`
	EventDir        string `json:"event_dir" yaml:"event_dir"`
	HandlerDir      string `json:"handler_dir" yaml:"handler_dir"`
	HandlerEntryDir string `json:"handler_entry_dir" yaml:"handler_entry_dir"`
	SQLDir          string `json:"sql_dir" yaml:"sql_dir"`
}

func (c *Config) GetRealTmpl(rootPath, appPath, appName string) *Tmpl {
	return &Tmpl{
		AppDir:          c.Tmpl.AppDir,
		ConfDir:         replacePath(c.Tmpl.ConfDir, rootPath, appPath, appName),
		DomainDir:       replacePath(c.Tmpl.DomainDir, rootPath, appPath, appName),
		EntityDir:       replacePath(c.Tmpl.EntityDir, rootPath, appPath, appName),
		DoDir:           replacePath(c.Tmpl.DoDir, rootPath, appPath, appName),
		ConvDoDir:       replacePath(c.Tmpl.ConvDoDir, rootPath, appPath, appName),
		DaoDir:          replacePath(c.Tmpl.DaoDir, rootPath, appPath, appName),
		PbDir:           replacePath(c.Tmpl.PbDir, rootPath, appPath, appName),
		ConvPbDir:       replacePath(c.Tmpl.ConvPbDir, rootPath, appPath, appName),
		ServiceDir:      replacePath(c.Tmpl.ServiceDir, rootPath, appPath, appName),
		ExtendDir:       replacePath(c.Tmpl.ExtendDir, rootPath, appPath, appName),
		RepoDir:         replacePath(c.Tmpl.RepoDir, rootPath, appPath, appName),
		RepoImplDir:     replacePath(c.Tmpl.RepoImplDir, rootPath, appPath, appName),
		EventDir:        replacePath(c.Tmpl.EventDir, rootPath, appPath, appName),
		HandlerDir:      replacePath(c.Tmpl.HandlerDir, rootPath, appPath, appName),
		HandlerEntryDir: replacePath(c.Tmpl.HandlerEntryDir, rootPath, appPath, appName),
		SQLDir:          replacePath(c.Tmpl.SQLDir, rootPath, appPath, appName),
	}
}

func replacePath(str string, rootPath, appPath, appName string) string {
	str = strings.Replace(str, "{rootPath}", rootPath, -1)
	str = strings.Replace(str, "{appPath}", appPath, -1)
	str = strings.Replace(str, "{appName}", appName, -1)
	return str
}
