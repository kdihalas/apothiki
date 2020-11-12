package http

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/spf13/afero"
	"sort"
)

type CatalogController struct {
	beego.Controller
}

type TagController struct {
	beego.Controller
}

func (this *CatalogController) Get() {
	var repos = []string{}
	dirInfo, err := afero.ReadDir(AppFs, ".")
	if err != nil {
		log.Error(err.Error())
	}
	for _, dir := range dirInfo {
		if dir.IsDir() {
			repos = append(repos, dir.Name())
		}

	}
	sort.Strings(repos)
	this.Data["json"] = &repos
	this.ServeJSON()
}

func (this *TagController) Get() {
	var tags = []string{}
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	fileInfo, err := afero.ReadDir(AppFs, fmt.Sprintf("%s/manifests", name))
	if err != nil {
		log.Error(err.Error())
	}
	for _, file := range fileInfo {
		tags = append(tags, file.Name())
	}
	sort.Strings(tags)
	this.Data["json"] = &tags
	this.ServeJSON()
}
