package api

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kdihalas/apothiki/pkg/utils"
	"github.com/spf13/afero"
	"sort"
)

var protectedDirs = []string{"manifests", "layers", "uploads"}

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
		singleLevel := true

		if !dir.IsDir() {
			continue
		}

		manifestsExists, _ := afero.DirExists(AppFs, fmt.Sprintf("./%s/manifests", dir.Name()))
		layersExists, _ := afero.DirExists(AppFs, fmt.Sprintf("./%s/layers", dir.Name()))

		subDirInfo, err := afero.ReadDir(AppFs, fmt.Sprintf("./%s", dir.Name()))
		if err != nil {
			log.Error(err.Error())
		}

		for _, subDir := range subDirInfo {
			if subDir.IsDir() && !contains(protectedDirs, subDir.Name()) {
				singleLevel = false
				repos = append(repos, fmt.Sprintf("%s/%s", dir.Name(), subDir.Name()))
			}
		}

		if (manifestsExists && layersExists) || singleLevel {
			repos = append(repos, dir.Name())
		}

	}
	sort.Strings(repos)
	this.Data["json"] = &repos
	this.ServeJSON()
}

func (this *TagController) Get() {
	var tags = []string{}
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
