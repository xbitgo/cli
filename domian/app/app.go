package app

import (
	"xbit/conf"
	"xbit/domian/gen"
)

type App struct {
	Pwd      string
	Name     string
	RootPath string
	Tmpl     *conf.Tmpl
}

func NewApp(rootPath string, appPath string, appName string) *App {
	return &App{
		Pwd:      appPath,
		Name:     appName,
		RootPath: rootPath,
		Tmpl:     conf.Global.GetRealTmpl(rootPath, appPath, appName),
	}
}

func (a *App) parseDirs() (dirs []string) {
	dirs = []string{
		a.Tmpl.ConfDir,
		a.Tmpl.EntityDir,
		a.Tmpl.ServiceDir,
		a.Tmpl.RepoImplDir,
	}
	dirTmpUniMap := map[string]struct{}{}
	for _, dir := range dirs {
		dirTmpUniMap[dir] = struct{}{}
	}
	items := a.ScanDir(a.Tmpl.DomainDir)
	for _, item := range items {
		if _, ok := dirTmpUniMap[item]; ok {
			continue
		}
		dirTmpUniMap[item] = struct{}{}
		dirs = append(dirs, item)
	}
	return dirs
}

func (a *App) Impl() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	dirs := a.parseDirs()
	return gm.Impl(dirs)
}

func (a *App) CRepo() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.CRepo()
}

func (a *App) CService() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.CService()
}

func (a *App) CHandler() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.CHandler()
}

func (a *App) Handler() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.Handler()
}

func (a *App) Tests() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.TestHandlers()
}

func (a *App) Do2Sql(dsn string) error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.Do2Sql(dsn)
}

func (a *App) Sql2Entity() error {
	return nil
}

func (a *App) Dao() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.Dao()
}

func (a *App) Conv() error {
	// todo
	return nil
}

func (a *App) Run() error {
	return nil
}
