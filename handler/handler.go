package handler

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"xbit/conf"
	"xbit/domian/project"
)

var pwd, _ = os.Getwd()

func CmdList() []*cobra.Command {
	return []*cobra.Command{
		// ------- project -----
		{
			Use:   "init",
			Short: "创建xbitgo项目",
			Long:  "创建xbitgo项目; 参数 {name}: name为项目名称 必须",
			Run:   initProject,
		},
		{
			Use:   "list",
			Short: "显示项目中应用列表",
			Long:  "显示项目中应用列表",
			Run:   list,
		},
		{
			Use:   "create",
			Short: "创建新的应用",
			Long:  "创建新的应用; 参数 {app}:app为应用名称 必须",
			Run:   create,
		},
		// ------- generate core go code -------
		{
			Use:   "handler",
			Short: "根据pb定义生成handler",
			Long:  "根据pb定义生成handler; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   handler,
		},
		{
			Use:   "impl",
			Short: "生成标注的interface的对应实现中 未实现的方法",
			Long:  "生成标注的interface的对应实现中 未实现的方法; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   impl,
		},
		{
			Use:   "conv",
			Short: "生成指定对象到指定对象的转换方法",
			Long:  "生成在converter内自定义的指定对象到指定对象的转换方法",
			Run:   impl,
		},
		{
			Use:   "dao",
			Short: "根据Do生成对应通用操作数据库方法",
			Long:  "根据Do生成对应通用操作数据库方法; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   dao,
		},
		{
			Use:   "protoc",
			Short: "pb生成包含grpc客户端",
			Long:  "pb生成包含grpc客户端",
			Run:   protoc,
		},
		{
			Use:   "generate",
			Short: "生成所有go代码 依次最新 generate > dao  > handler > impl > conv",
			Long:  "生成标注的interface的对应实现中 未实现的方法; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   generate,
		},

		// ------- other code generate -------
		{
			Use:   "c.repo",
			Short: "根据Entity生成基本的仓库接口及实现",
			Long:  "根据Entity生成基本的仓库接口及实现; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   cRepo,
		},
		{
			Use:   "c.service",
			Short: "根据Entity生成基本的Service实现",
			Long:  "根据Entity生成基本的Service实现; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   cService,
		},
		{
			Use:   "c.handler",
			Short: "根据Entity生成基本的接口定义及实现",
			Long:  "根据Entity生成基本的接口定义及实现; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   cHandler,
		},
		{
			Use:   "do2Sql",
			Short: "根据Do生成对应数据表创建/新增SQL",
			Long:  "根据Do生成对应数据表创建/新增SQL; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   do2Sql,
		},
		{
			Use:   "sql2Entity",
			Short: "根据数据库表结构生成Entity",
			Long:  "根据数据库表结构生成Entity; 参数 {app}; app为应用名称 支持多个 不传默认全部应用",
			Run:   sql2Entity,
		},
		{
			Use:   "d.ts",
			Short: "根据pb定义生成对应ts文件(使用json)",
			Long:  "根据pb定义生成对应ts文件; 参数 {app}; app为应用名称 必须",
			Run:   d2Ts,
		},
		{
			Use:   "d.dart",
			Short: "根据pb定义生成对应dart文件(使用json)",
			Long:  "根据pb定义生成对应dart文件; 参数 {app}; app为应用名称 必须",
			Run:   d2Dart,
		},
		{
			Use:   "tests",
			Short: "生成接口单元测试用例",
			Long:  "生成接口单元测试用例; 参数 {app}; app为应用名称 必须",
			Run:   tests,
		},
	}
}

func initProject(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("请输入项目名! ")
	}
	name := args[0]
	project.InitProject(pwd, name)
}

func list(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	for _, appName := range p.List() {
		fmt.Printf("---- %s \n", appName)
	}
}

func create(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	if len(args) == 0 {
		log.Fatalf("请输入应用名! ")
	}
	name := args[0]
	p.Create(name)
}

func generate(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Generate()
}

func handler(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Handler()
}

func impl(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Impl()
}

func cRepo(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.CRepo()
}

func cService(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.CService()
}
func cHandler(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.CHandler()
}

func do2Sql(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	var dsn = ""
	for _, arg := range args {
		if strings.HasPrefix(arg, "dsn=") {
			dsn = strings.TrimPrefix(arg, "dsn=")
		} else if strings.HasPrefix(arg, "db=") {
			db := strings.TrimPrefix(arg, "db=")
			if v, ok := conf.Global.DBList[db]; ok {
				dsn = v
			}
		}
	}
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Do2Sql(dsn)
}

func sql2Entity(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Sql2Entity()
}

func dao(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Dao()
}

func protoc(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	p.Protoc()
}

func d2Ts(cmd *cobra.Command, args []string) {
	fmt.Println("wip ....")
}

func d2Dart(cmd *cobra.Command, args []string) {
	fmt.Println("wip ....")
}

func tests(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num != 1 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Tests()
}
