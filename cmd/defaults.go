package main

import "github.com/spf13/viper"

func init() {
	// Service defaults
	viper.SetDefault("mode", "repo")

	// Web server defaults
	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("port", 8080)

	// Storage defaults
	viper.SetDefault("data_path", "/opt/apothiki/data")

}
