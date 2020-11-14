package ui

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/spf13/viper"
	"sort"
)

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