package http

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"io/ioutil"
)

type BlobUpload struct {
	beego.Controller
}

type BlobUploads struct {
	beego.Controller
}

func (this *BlobUpload) Post() {
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))

	// Calculate Upload UUID
	uuid := uuid.NewV4()

	// Generate location path
	location := fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uuid)

	this.Ctx.Output.Header("Location", location)
	this.Ctx.ResponseWriter.WriteHeader(202)
}

func (this *BlobUploads) Get() {
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	uuid := this.Ctx.Input.Param(":uuid")

	fileInfo, err := AppFs.Stat(fmt.Sprintf("%s/uploads/%s", name, uuid))
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	location := fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uuid)

	this.Ctx.Output.Header("Range", fmt.Sprintf("0-%d", fileInfo.Size()))
	this.Ctx.Output.Header("Location", location)
	this.Ctx.Output.Header("Docker-Upload-UUID", uuid)
	this.Ctx.ResponseWriter.WriteHeader(204)

}

func (this *BlobUploads) Patch() {
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	uuid := this.Ctx.Input.Param(":uuid")
	err := AppFs.MkdirAll(fmt.Sprintf("%s/uploads", name), 0755)
	if err != nil {
		log.Error(err.Error())
	}
	file, err := AppFs.Create(fmt.Sprintf("%s/uploads/%s", name, uuid))
	if err != nil {
		log.Error(err.Error())
	}
	// Get data from request body
	chunk := this.Ctx.Input.RequestBody

	// Write chunk to file
	_, err = file.Write(chunk)
	if err != nil {
		log.Error(err.Error())
	}
	// Close file after chunk was appended
	file.Close()

	// Generate location path
	location := fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uuid)

	fileInfo, err := AppFs.Stat(fmt.Sprintf("%s/uploads/%s", name, uuid))
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}

	this.Ctx.Output.Header("Range", fmt.Sprintf("0-%d", fileInfo.Size()))
	this.Ctx.Output.Header("Location", location)
	this.Ctx.Output.Header("Docker-Upload-UUID", uuid)
	this.Ctx.ResponseWriter.WriteHeader(204)
}

func (this *BlobUploads) Put() {
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	uuid := this.Ctx.Input.Param(":uuid")
	digest := this.Ctx.Input.Query("digest")

	// Create Layers folder
	err := AppFs.MkdirAll(fmt.Sprintf("%s/layers", name), 0755)
	if err != nil {
		log.Error(err.Error())
	}

	// Create Layer file
	file, err := AppFs.Create(fmt.Sprintf("%s/layers/%s", name, digest))
	if err != nil {
		log.Error(err.Error())
	}
	tmpFile, err := AppFs.Open(fmt.Sprintf("%s/uploads/%s", name, uuid))
	defer tmpFile.Close()
	if err != nil {
		log.Error(err.Error())
	}
	// Read tmp File contents
	tmpFileContents, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		log.Error(err.Error())
	}
	// Write contents to layer file
	file.Write(tmpFileContents)
	file.Close()

	AppFs.Remove(fmt.Sprintf("%s/uploads/%s", name, uuid))
	location := fmt.Sprintf("/v2/%s/blobs/%s", name, digest)

	this.Ctx.Output.Header("Location", location)
	this.Ctx.Output.Header("Docker-Content-Digest", digest)
	this.Ctx.ResponseWriter.WriteHeader(201)
}

func (this *BlobUploads) Delete() {
	// Get repo name
	name := getContainerName(this.Ctx.Input.Param(":repo"), this.Ctx.Input.Param(":name"))
	uuid := this.Ctx.Input.Param(":uuid")
	err := AppFs.Remove(fmt.Sprintf("%s/uploads/%s", name, uuid))
	if err != nil {
		log.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(200)
}
