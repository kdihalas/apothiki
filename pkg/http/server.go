package http

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/kdihalas/apothiki/pkg/http/api"
	"github.com/kdihalas/apothiki/pkg/http/ui"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

func Run() {

	// Load API config
	api.Init()

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// Redirect / to /ui
	beego.Get("/", func(c *context.Context) {
		c.Redirect(302, "/ui")
	})
	
	// UI routes
	ui := beego.NewNamespace("/ui",
		beego.NSRouter("/", &ui.IndexController{}, "get:Get"),
		beego.NSRouter("/:name", &ui.RepoController{}, "get:Get"),
		beego.NSRouter("/:repo/:name", &ui.RepoController{}, "get:Get"),
	)

	// API routes
	api := beego.NewNamespace("/v2",
		beego.NSRouter("/", &api.VersionController{}, "get:Get"),
		beego.NSRouter("/_catalog", &api.CatalogController{}, "get:Get"),
		beego.NSRouter("/:name/tags/list", &api.TagController{}, "get:Get"),
		beego.NSRouter("/:repo/:name/tags/list", &api.TagController{}, "get:Get"),

		beego.NSRouter("/:name/manifests/:reference", &api.ManifestController{}, "get:Get;head:Head;put:Put;delete:Delete"),
		beego.NSRouter("/:repo/:name/manifests/:reference", &api.ManifestController{}, "get:Get;head:Head;put:Put;delete:Delete"),

		beego.NSRouter("/:name/blobs/uploads", &api.BlobUpload{}, "post:Post"),
		beego.NSRouter("/:repo/:name/blobs/uploads", &api.BlobUpload{}, "post:Post"),

		beego.NSRouter("/:name/blobs/uploads/:uuid", &api.BlobUploads{}, "get:Get;patch:Patch;put:Put;delete:Delete"),
		beego.NSRouter("/:repo/:name/blobs/uploads/:uuid", &api.BlobUploads{}, "get:Get;patch:Patch;put:Put;delete:Delete"),

		beego.NSRouter("/:name/blobs/:digest", &api.DigestController{}, "get:Get;head:Head;delete:Delete"),
		beego.NSRouter("/:repo/:name/blobs/:digest", &api.DigestController{}, "get:Get;head:Head;delete:Delete"),
	)

	// Add namespaced router to beego server
	beego.AddNamespace(ui, api)

	// Set template file path
	beego.BConfig.WebConfig.ViewsPath = "pkg/views"
	// Set server loglevel
	beego.BConfig.Log.AccessLogs = true
	// Enable copy request body
	beego.BConfig.CopyRequestBody = true
	// Enable Admin interface
	beego.BConfig.Listen.EnableAdmin = true
	// Read config from yaml file
	beego.BConfig.Listen.HTTPAddr = viper.GetString("addr")
	beego.BConfig.Listen.HTTPPort = viper.GetInt("port")

	// Start server
	beego.Run()
}
