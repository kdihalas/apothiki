package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
)

func GetUpstream() string {
	// Get upstream servers config
	upstreamServers := viper.Get("upstream").([]interface{})
	// Size them
	size := len(upstreamServers)
	// Return a random one and cast to map[string]string
	upstream := convertMapString(upstreamServers[rand.Intn(size-(size-1))].(map[interface{}]interface{}))

	// Generate selected URL
	selected :=  fmt.Sprintf(fmt.Sprintf("%s://%s", upstream["transport"], upstream["addr"]))

	return selected
}

func GetContainerName(repository string, image string) string {
	name := fmt.Sprintf("%s/%s", repository, image)
	if repository == "" {
		name = fmt.Sprintf("%s", image)
	}
	return name
}

func convertMapString(input map[interface{}]interface{}) (map[string]string) {
	var output = map[string]string{}
	for key, value := range input {
		output[key.(string)] = value.(string)
	}
	return output
}