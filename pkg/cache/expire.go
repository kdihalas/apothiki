package cache

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"time"
)

var (
	log   = logs.NewLogger(10000)
)

func ExpireCache() {
	AppFs := afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("data_path"))
	// Get expiration time
	expiration := viper.GetDuration("expire")
	log.Info("Expiration time set to %s", expiration)

	for {
		afero.Walk(AppFs, ".", func(path string, info os.FileInfo, err error) error {
			timeElapsed := time.Now().Sub(info.ModTime())
			if !info.IsDir() {
				if timeElapsed > expiration {
					log.Info(fmt.Sprintf("Deleting file %s from cache", path))
					AppFs.Remove(path)
				}
			}
			return nil
		})

		time.Sleep(time.Duration(10 * time.Second))
	}
}
