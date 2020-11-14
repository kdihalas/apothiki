package ui

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/kdihalas/apothiki/pkg/utils"
	"github.com/spf13/viper"
	"sort"
)

type RepoController struct {
	LayoutController
}

func (this *RepoController) Get() {
	var tags []string
	// Get repo name
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	req := httplib.Get(fmt.Sprintf("http://%s:%d/v2/%s/tags/list", viper.GetString("addr"), viper.GetInt("port"), name))
	err := req.ToJSON(&tags)
	if err != nil {
		log.Error(err.Error())
	}
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))
	this.Data["repo"] = name
	this.Data["tags"] = tags
}