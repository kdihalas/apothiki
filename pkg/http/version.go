package http

import "github.com/astaxie/beego"

type VersionController struct {
	beego.Controller
}

func (this *VersionController) Get() {
	this.Ctx.ResponseWriter.WriteHeader(200)
	return
}
