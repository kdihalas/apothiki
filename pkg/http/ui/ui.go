package ui

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	log   = logs.NewLogger(10000)
)

type LayoutController struct {
	beego.Controller
}

func (this *LayoutController) Prepare() {
	this.Layout = "layout.tpl"
}
