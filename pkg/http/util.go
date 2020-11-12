package http

import (
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
)

func getUpstream() string {
	upstreamServers := viper.GetStringSlice("upstream")
	size := len(upstreamServers)
	selected := upstreamServers[rand.Intn(size - (size-1))]
	if upstreamHealthCheck(selected){
		return selected
	}
	return getUpstream()
}

func upstreamHealthCheck(upstream string) bool {
	_, err := http.Get(fmt.Sprintf("http://%s/", upstream))
	if err != nil {
		return false
	}
	return true
}

func getContainerName(repository string, image string) string {
	name := fmt.Sprintf("%s/%s", repository, image)
	if repository == "" {
		name = fmt.Sprintf("%s", image)
	}
	return name
}