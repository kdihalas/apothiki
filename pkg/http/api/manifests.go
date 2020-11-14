package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kdihalas/apothiki/pkg/utils"
	"io"
	"io/ioutil"
)

type ManifestController struct {
	beego.Controller
}

func (this *ManifestController) Get() {
	var doc map[string]interface{}
	reference := this.Ctx.Input.Param(":reference")

	// Get repo name
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	file, err := AppFs.Open(fmt.Sprintf("%s/manifests/%s", name, reference))
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err.Error())
	}
	json.Unmarshal(data, &doc)

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Error(err.Error())
	}

	this.Data["json"] = &doc
	this.Ctx.Output.Header("Docker-Content-Digest", fmt.Sprintf("sha256:%x", h.Sum(nil)))
	this.Ctx.Output.ContentType(fmt.Sprintf("%s", doc["mediaType"]))
	this.ServeJSON()
}

func (this *ManifestController) Head() {
	reference := this.Ctx.Input.Param(":reference")

	// Get repo name
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	file, err := AppFs.Open(fmt.Sprintf("%s/manifests/%s", name, reference))
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	this.Ctx.Output.Header("Docker-Content-Digest", fmt.Sprintf("sha256:%x", h.Sum(nil)))
	this.ServeJSON()
}

func (this *ManifestController) Put() {
	reference := this.Ctx.Input.Param(":reference")

	// Get repo name
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	err := AppFs.MkdirAll(fmt.Sprintf("%s/manifests", name), 0755)
	if err != nil {
		log.Error(err.Error())
	}
	file, err := AppFs.Create(fmt.Sprintf("%s/manifests/%s", name, reference))
	defer file.Close()
	if err != nil {
		log.Error(err.Error())
	}

	// Get data from request body
	chunk := this.Ctx.Input.RequestBody

	// Calculate sha256 hash for the file
	h := sha256.New()
	file.Write(chunk)
	if _, err := io.Copy(h, file); err != nil {
		log.Error(err.Error())
	}

	this.Ctx.Output.Header("Docker-Content-Digest", fmt.Sprintf("sha256:%x", h.Sum(nil)))
	this.Ctx.ResponseWriter.WriteHeader(201)
}

func (this *ManifestController) Delete() {
	reference := this.Ctx.Input.Param(":reference")

	// Get repo name
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	err := AppFs.Remove(fmt.Sprintf("%s/manifests/%s", name, reference))
	if err != nil {
		log.Error(err.Error())
	}

	this.Ctx.ResponseWriter.WriteHeader(202)
}
