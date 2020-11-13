package http

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/httplib"
	"github.com/spf13/viper"
	"sort"
	"time"
)

var (
	mc cache.Cache
	cacheExpire = 60 * time.Second
)

type LayoutController struct {
	beego.Controller
}

func (this *LayoutController) Prepare() {
	this.Layout = "layout.tpl"
}

type IndexController struct {
	LayoutController
}

func (this *IndexController) Get() {
	var repos []string
	req := httplib.Get(fmt.Sprintf("http://%s:%d/v2/_catalog", viper.GetString("addr"), viper.GetInt("port")))
	err := req.ToJSON(&repos)
	if err != nil {
		log.Error(err.Error())
	}
	sort.Strings(repos)

	this.Data["repos"] = repos
}

type RepoController struct {
	LayoutController
}

func (this *RepoController) Get() {
	var tags []string
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	req := httplib.Get(fmt.Sprintf("http://%s:%d/v2/%s/tags/list", viper.GetString("addr"), viper.GetInt("port"), name))
	err := req.ToJSON(&tags)
	if err != nil {
		log.Error(err.Error())
	}
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))
	this.Data["repo"] = name
	this.Data["tags"] = tags
}
