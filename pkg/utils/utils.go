package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
)

func GetUpstream() string {
	// Get upstream servers config
	upstreamServers := viper.Get("upstream").([]interface{})
	// Size them
	size := len(upstreamServers)
	// Return a random one and cast to map[string]string
	upstream := convertMapString(upstreamServers[rand.Intn(size-(size-1))].(map[interface{}]interface{}))

	// Generate selected URL
	return fmt.Sprintf(fmt.Sprintf("%s://%s", upstream["transport"], upstream["addr"]))
}

func upstreamHealthy(upstream string) bool {
	resp, _ := http.Get(fmt.Sprintf("%s/v2"))
	if resp.StatusCode != 200 {
		return false
	}
	return true
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