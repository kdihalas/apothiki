package http

import (
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
)

type DigestController struct {
	beego.Controller
}

func (this *DigestController) Get(){
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	// Get digest
	digest := this.Ctx.Input.Param(":digest")

	file, err := AppFs.Open(fmt.Sprintf("%s/layers/%s", name, digest))
	defer file.Close()
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	data, err := ioutil.ReadAll(file)
	this.Ctx.Output.Header("Docker-Content-Digest", fmt.Sprintf("%s", digest))
	this.Ctx.ResponseWriter.Write(data)
	this.Ctx.ResponseWriter.WriteHeader(200)
}
func (this *DigestController) Head(){
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	// Get digest
	digest := this.Ctx.Input.Param(":digest")

	fileInfo, err := AppFs.Stat(fmt.Sprintf("%s/layers/%s", name, digest))
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	this.Ctx.Output.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	this.Ctx.Output.Header("Docker-Content-Digest", fmt.Sprintf("%s", digest))
	this.Ctx.ResponseWriter.WriteHeader(200)

}

func (this *DigestController) Delete(){
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	// Get digest
	digest := this.Ctx.Input.Param(":digest")

	err := AppFs.Remove(fmt.Sprintf("%s/layers/%s", name, digest))
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}

	this.Ctx.ResponseWriter.WriteHeader(200)

}