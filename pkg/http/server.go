package http

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	AppFs afero.Fs
	log = logs.NewLogger(10000)
)

func Run() {
	AppFs = afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("data_path"))

	// Set log level from config
	log.SetLogger("console")

	ns := beego.NewNamespace("/v2",
		beego.NSRouter("/", &VersionController{}, "get:Get"),
		beego.NSRouter("/_catalog", &CatalogController{}, "get:Get"),
		beego.NSRouter("/:name/tags/list", &TagController{}, "get:Get"),
		beego.NSRouter("/:repo/:name/tags/list", &TagController{}, "get:Get"),

		beego.NSRouter("/:name/manifests/:reference", &ManifestController{},"get:Get;head:Head;put:Put;delete:Delete"),
		beego.NSRouter("/:repo/:name/manifests/:reference", &ManifestController{},"get:Get;head:Head;put:Put;delete:Delete"),

		beego.NSRouter("/:name/blobs/uploads", &BlobUpload{},"post:Post"),
		beego.NSRouter("/:repo/:name/blobs/uploads", &BlobUpload{},"post:Post"),

		beego.NSRouter("/:name/blobs/uploads/:uuid", &BlobUploads{},"get:Get;patch:Patch;put:Put;delete:Delete"),
		beego.NSRouter("/:repo/:name/blobs/uploads/:uuid", &BlobUploads{},"get:Get;patch:Patch;put:Put;delete:Delete"),

		beego.NSRouter("/:name/blobs/:digest", &DigestController{}, "get:Get;head:Head;delete:Delete"),
		beego.NSRouter("/:repo/:name/blobs/:digest", &DigestController{}, "get:Get;head:Head;delete:Delete"),
	)

	// Add namespaced router to beego server
	beego.AddNamespace(ns)

	// Set server loglevel
	beego.BConfig.Log.AccessLogs = true

	// Enable Admin interface
	beego.BConfig.Listen.EnableAdmin = true
	// Read config from yaml file
	beego.BConfig.Listen.HTTPAddr = viper.GetString("addr")
	beego.BConfig.Listen.HTTPPort = viper.GetInt("port")


	// Start server
	beego.Run()
}
