package api

import (
	"github.com/astaxie/beego/logs"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	AppFs afero.Fs
	log = logs.NewLogger(10000)
)

func Init() {
	// Load FS jailed inside the path of the config
	AppFs = afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("data_path"))
}