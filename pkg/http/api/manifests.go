package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kdihalas/apothiki/pkg/cache"
	"github.com/kdihalas/apothiki/pkg/utils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
)

type ManifestController struct {
	beego.Controller
}

func (this *ManifestController) Get() {
	var doc map[string]interface{}
	var file afero.File
	var fromUpstream = false
	var data []byte
	reference := this.Ctx.Input.Param(":reference")

	// Get repo name
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	// Path
	path := fmt.Sprintf("%s/manifests/%s", name, reference)
	dir := fmt.Sprintf("%s/manifests", name)

	// Open file if it exists
	file, err := AppFs.Open(path)
	if err != nil {
		// Run only when mode is cache
		if viper.GetString("mode") == "cache" {
			// Check if upstream server has the file in question
			if upstream, exists := cache.IsAvailableUpstream(this.Ctx.Request.RequestURI); exists {
				// Cache the file locally
				data, err = cache.CacheLocal(upstream, dir, path)
				if err != nil {
					log.Error(err.Error())
					this.Ctx.ResponseWriter.WriteHeader(500)
					return
				}
				// Flag that the file was fetched from upstream
				fromUpstream = true
			} else {
				log.Error(err.Error())
				this.Ctx.ResponseWriter.WriteHeader(404)
				return
			}
		} else {
			log.Error(err.Error())
			this.Ctx.ResponseWriter.WriteHeader(404)
			return
		}
	} else {
		defer file.Close()
	}

	// If the file was not fetched from upstream read data from local
	if !fromUpstream {
		data, err = ioutil.ReadAll(file)
		if err != nil {
			log.Error(err.Error())
		}
	}

	// Parse manifest file
	json.Unmarshal(data, &doc)

	// Get the file digest from the file
	config := doc["config"].(map[string]interface{})
	digest := config["digest"].(string)

	// Return the file
	this.Data["json"] = &doc
	this.Ctx.Output.Header("Docker-Content-Digest", fmt.Sprintf("%s", digest))
	this.Ctx.Output.ContentType(fmt.Sprintf("%s", doc["mediaType"]))
	this.ServeJSON()
}

func (this *ManifestController) Head() {
	var doc map[string]interface{}

	reference := this.Ctx.Input.Param(":reference")

	// Get repo name
	name := utils.GetContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	file, err := AppFs.Open(fmt.Sprintf("%s/manifests/%s", name, reference))
	if err != nil {
		// Run only when mode is cache
		if viper.GetString("mode") == "cache" {
			// Check if the file exists upstream and return the digest
			if _, digest, exists := cache.HeadUpstream(this.Ctx.Request.RequestURI); exists {
				this.Ctx.Output.Header("Docker-Content-Digest", digest)
				this.ServeJSON()
				return
			} else {
				log.Error(err.Error())
				this.Ctx.ResponseWriter.WriteHeader(404)
				return
			}
		} else {
			log.Error(err.Error())
			this.Ctx.ResponseWriter.WriteHeader(404)
			return
		}
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err.Error())
	}

	// Parse manifest file
	json.Unmarshal(data, &doc)

	// Get the file digest from the file
	config := doc["config"].(map[string]interface{})
	digest := config["digest"].(string)

	this.Ctx.Output.Header("Docker-Content-Digest", digest)
	this.Ctx.ResponseWriter.WriteHeader(200)
	return
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
	if err != nil {
		log.Error(err.Error())
	}
	defer file.Close()

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
