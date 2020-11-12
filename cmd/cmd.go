package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/kdihalas/apothiki/pkg/cache"
	"github.com/kdihalas/apothiki/pkg/http"
	"github.com/kdihalas/apothiki/pkg/sync"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile string
	cmd = &cobra.Command{
		Use: "apothiki",
		Short: "A docker registry",
		Long: `This is a simple docker registry with some advanced capabilities`,
		Run: func(cmd *cobra.Command, args []string){
			if viper.GetString("mode") == "cache" {
				go cache.ExpireCache()
			}
			if viper.GetString("mode") == "replica" {
				go sync.Sync()
			}
			http.Run()
		},
	}
	log = logs.NewLogger(10000)
)

func init() {
	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
}

func er(msg string) {
	log.Error("Error: %s", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err.Error())
		}

		// Search config in home directory with name ".docker-repo" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config %s", viper.ConfigFileUsed())
	}
}