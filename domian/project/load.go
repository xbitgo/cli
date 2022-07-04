package project

import (
	"io/ioutil"
)

func (p *Project) LoadApps() {
	apps := make([]string, 0)
	appsPwd := p.Pwd + "/" + p.AppDir
	apps = p.findApps(appsPwd)
	p.AppList = apps
}

func (p *Project) isApp(pwd string) bool {
	fileInfos, err := ioutil.ReadDir(pwd)
	if err != nil {
		return false
	}
	for _, fi := range fileInfos {
		if fi.Name() == "main.go" && !fi.IsDir() {
			return true
		}
	}
	return false
}

func (p *Project) findApps(pwd string) []string {
	apps := make([]string, 0)
	fileInfos, err := ioutil.ReadDir(pwd)
	if err != nil {
		return apps
	}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			iPwd := pwd + "/" + fi.Name()
			if p.isApp(iPwd) {
				apps = append(apps, iPwd)
			} else { //多层结构
				apps = append(apps, p.findApps(iPwd)...)
			}
		}
	}
	return apps
}
