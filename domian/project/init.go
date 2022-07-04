package project

import (
	"fmt"
	"os"
)

// InitProject 创建项目
func InitProject(pwd, name string) {
	dir := fmt.Sprintf("%s/%s", pwd, name)
	err := os.MkdirAll(dir, 0664)
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
				p.ProjectTpl.ReadFile("template/gomod.txt")
			}
			fmt.Println(f.Name())
		}
	}

}
