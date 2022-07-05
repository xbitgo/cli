package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/xbitgo/core/tools/tool_file"
)

// InitProject 创建项目
func InitProject(pwd, projectName string) {
	dir := fmt.Sprintf("%s/%s", pwd, projectName)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		panic(err)
	}
	p := NewProject(dir)
	fs, err := p.ProjectTpl.ReadDir("template")
	if err != nil {
		panic(err)
	}
	for _, f := range fs {
		if !f.IsDir() {
			if f.Name() == "gomod.txt" {
				buf, _ := p.ProjectTpl.ReadFile("template/gomod.txt")
				fName := fmt.Sprintf("%s/%s", dir, "go.mod")
				str := strings.Replace(string(buf), "{{ProjectName}}", projectName, -1)
				err := tool_file.WriteFile(fName, []byte(str))
				if err != nil {
					panic(err)
				}
				continue
			}
			buf, _ := p.ProjectTpl.ReadFile("template/" + f.Name())
			fName := fmt.Sprintf("%s/%s", dir, f.Name())
			str := strings.Replace(string(buf), "{{ProjectName}}", projectName, -1)
			tool_file.WriteFile(fName, []byte(str))
		} else {
			p.writeTplDir(dir, "template/"+f.Name(), f.Name(), projectName)
		}
	}

}

func (p *Project) writeTplDir(dir string, tplDir string, name string, projectName string) {
	dir = dir + "/" + name
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		panic(err)
	}
	fs, err := p.ProjectTpl.ReadDir(tplDir)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {
		if !f.IsDir() {
			buf, _ := p.ProjectTpl.ReadFile(tplDir + "/" + f.Name())
			fName := fmt.Sprintf("%s/%s", dir, f.Name())
			str := strings.Replace(string(buf), "{{ProjectName}}", projectName, -1)
			tool_file.WriteFile(fName, []byte(str))
		} else {
			p.writeTplDir(dir, tplDir+"/"+f.Name(), f.Name(), projectName)
		}
	}
}
