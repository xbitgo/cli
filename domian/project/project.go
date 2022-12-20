package project

import (
	"embed"
	"log"
	"strings"
	"xbit/domian/gen"
	"xbit/utils"

	"xbit/conf"
	"xbit/domian/app"

	"github.com/xbitgo/core/di"
)

type Project struct {
	Pwd        string
	AppDir     string
	AppList    []string
	activeApps []*app.App
	ProjectTpl embed.FS `di:"project_tpl"`
}

func NewProject(pwd string) *Project {
	p := &Project{
		Pwd:        pwd,
		AppDir:     conf.Global.Tmpl.AppDir,
		AppList:    make([]string, 0),
		activeApps: make([]*app.App, 0),
	}
	p.LoadApps()
	// di注册project
	di.Register("project", p)
	di.MustBind(p)
	return p
}

func (p *Project) SetActiveApps(apps ...string) int {
	appMap := map[string]string{}
	for _, s := range p.AppList {
		appName := strings.TrimPrefix(s, p.Pwd+"/"+p.AppDir+"/")
		appMap[appName] = s
	}
	if len(apps) >= 0 {
		for _, s := range apps {
			if appPath, ok := appMap[s]; ok {
				actApp := app.NewApp(p.Pwd, appPath, s)
				p.activeApps = append(p.activeApps, actApp)
			}
		}
	} else {
		for s, appPath := range appMap {
			actApp := app.NewApp(p.Pwd, appPath, s)
			p.activeApps = append(p.activeApps, actApp)
		}
	}
	return len(p.activeApps)
}

// List 应用列表
func (p *Project) List() []string {
	apps := make([]string, len(p.AppList))
	for i, s := range p.AppList {
		apps[i] = strings.TrimPrefix(s, p.Pwd+"/"+p.AppDir+"/")
	}
	return apps
}

func (p *Project) RootPath() string {
	return p.Pwd
}

// Create 创建新应用
func (p *Project) Create(name string) {
	gm := gen.NewManager(conf.Global.Tmpl, name, "")
	_ = gm.App(p.Pwd, name)
}

// Protoc 生成Pb 文件及grpc客户端
func (p *Project) Protoc() {
	gm := gen.NewManager(conf.Global.Tmpl, "", "")
	// 生成pb文件和grpc client
	dirs := []string{
		p.RootPath() + "/proto/base",
		p.RootPath() + "/proto/apps",
	}
	for _, dir := range dirs {
		pbFiles := utils.ScanPbDir(dir)
		for _, pbFile := range pbFiles {
			gm.Client(pbFile)
			pbFile = strings.Replace(pbFile, p.RootPath()+"/proto/", "", 1)
			_ = gm.Protoc(pbFile)
		}
	}

}

// Generate 生成所有GO代码
func (p *Project) Generate() {
	for _, actApp := range p.activeApps {
		if err := actApp.Generate(); err != nil {
			log.Panicf("Watch app[%s] err: %v", actApp.Name, err)
		}
	}
	// pb生成
	//p.Protoc()
}

// Impl 生成
func (p *Project) Impl() {
	for _, actApp := range p.activeApps {
		if err := actApp.Impl(); err != nil {
			log.Panicf("Impl app[%s] err: %v", actApp.Name, err)
		}
	}
}

// CRepo 生成
func (p *Project) CRepo() {
	for _, actApp := range p.activeApps {
		if err := actApp.CRepo(); err != nil {
			log.Panicf("Impl app[%s] err: %v", actApp.Name, err)
		}
	}
}

// CService 生成
func (p *Project) CService() {
	for _, actApp := range p.activeApps {
		if err := actApp.CService(); err != nil {
			log.Panicf("Impl app[%s] err: %v", actApp.Name, err)
		}
	}
}

// Handler 生成
func (p *Project) Handler() {
	for _, actApp := range p.activeApps {
		if err := actApp.Handler(); err != nil {
			log.Panicf("Handler app[%s] err: %v", actApp.Name, err)
		}
	}
}

// Tests 生成
func (p *Project) Tests() {
	for _, actApp := range p.activeApps {
		if err := actApp.Tests(); err != nil {
			log.Panicf("Tests app[%s] err: %v", actApp.Name, err)
		}
	}
}

// Do2Sql 生成
func (p *Project) Do2Sql(dsn string) {
	for _, actApp := range p.activeApps {
		if err := actApp.Do2Sql(dsn); err != nil {
			log.Panicf("Do2Sql app[%s] err: %v", actApp.Name, err)
		}
	}
}

// Sql2Entity 生成
func (p *Project) Sql2Entity() {
	for _, actApp := range p.activeApps {
		if err := actApp.Sql2Entity(); err != nil {
			log.Panicf("Sql2Entity app[%s] err: %v", actApp.Name, err)
		}
	}
}

// Dao 生成
func (p *Project) Dao() {
	for _, actApp := range p.activeApps {
		if err := actApp.Dao(); err != nil {
			log.Panicf("Dao app[%s] err: %v", actApp.Name, err)
		}
	}
}
