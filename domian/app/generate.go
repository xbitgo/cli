package app

import (
	"log"
	"xbit/domian/gen"
	"xbit/domian/parser"
)

type Pkg struct {
	Dir      string
	IsEntity bool
	Parser   *parser.IParser
}

func (a *App) Generate() error {
	var (
		dirs    = a.parseDirs()
		pkgList = map[string]*Pkg{}
	)
	for _, dir := range dirs {
		pkg := &Pkg{
			Dir: dir,
		}
		if dir == a.Tmpl.EntityDir {
			pkg.IsEntity = true
		}
		pr, err := parser.Scan(dir, parser.ParseTypeWatch)
		if err != nil {
			log.Printf("generate parser pkg[%s] err: %v \n", dir, err)
			return err
		}
		pkg.Parser = pr
		pkgList[dir] = pkg
	}

	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	entityPkg := pkgList[a.Tmpl.EntityDir]
	confPkg := pkgList[a.Tmpl.ConfDir]
	// entity 生成
	for _, xst := range entityPkg.Parser.StructList {
		err := gm.Do(xst)
		if err != nil {
			log.Printf("gen do err: %v \n", err)
			return err
		}
	}
	gm.EntityTypeDef()
	// do 的TypeDef
	gm.DoTypeDef()
	// 生成pb
	err := gm.Pb(entityPkg.Parser.StructList, entityPkg.Parser.OtherStruct)
	if err != nil {
		log.Printf("gen pb err: %v \n", err)
		return err
	}
	// SDI 生成
	if xst, ok := confPkg.Parser.StructList["Config"]; ok {
		err = gm.SDI(confPkg.Parser.Package, xst)
		if err != nil {
			log.Printf("gen sdi err: %v \n", err)
			return err
		}
	}
	// 其他的 DI生成
	for s, pkg := range pkgList {
		if s == a.Tmpl.EntityDir {
			continue
		}
		if s == a.Tmpl.ConfDir {
			continue
		}
		if s == a.Tmpl.RepoDir {
			continue
		}
		if s == a.Tmpl.EventDir {
			continue
		}
		err = gm.DI(pkg.Parser)
		if err != nil {
			log.Printf("gen di err: %v \n", err)
			return err
		}
	}

	// 其他步骤
	//> dao > impl > handler > conv
	if err := a.Dao(); err != nil {
		log.Printf("Dao app[%s] err: %v", a.Name, err)
		return err
	}
	if err := a.Impl(); err != nil {
		log.Printf("Impl app[%s] err: %v", a.Name, err)
		return err
	}
	if err := a.Handler(); err != nil {
		log.Printf("Handler app[%s] err: %v", a.Name, err)
		return err
	}
	if err := a.Conv(); err != nil {
		log.Printf("Conv app[%s] err: %v", a.Name, err)
		return err
	}
	return nil
}
