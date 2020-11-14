package cache

import (
	"fmt"
	"github.com/kdihalas/apothiki/pkg/utils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func IsAvailableUpstream(url string) (string, bool) {
	upstream := utils.GetUpstream()
	log.Info("Checking upstream server %s", upstream)
	resp, err := http.Head(fmt.Sprintf("%s%s", upstream, url))
	if err != nil {
		log.Error(err.Error())
	}
	if resp.StatusCode == 404 {
		return upstream, false
	}
	return upstream, true
}

func HeadUpstream(url string) (string, string, bool) {
	upstream := utils.GetUpstream()
	log.Info("Checking upstream server %s", upstream)
	resp, err := http.Head(fmt.Sprintf("%s%s", upstream, url))
	if err != nil {
		log.Error(err.Error())
	}
	if resp.StatusCode == 404 {
		return upstream, resp.Header.Get(""), false
	}
	return upstream,resp.Header.Get("Docker-Content-Digest"), true
}

func PullFromUpstream(url string) ([]byte, error){
	log.Info("Pulling file from upstream url %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err.Error())
	}
	if resp.StatusCode == 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return data, nil
	}
	return nil, err
}

func CacheLocal(upstream string, dir string, path string) ([]byte, error){
	AppFs := afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("data_path"))
	err := AppFs.MkdirAll(dir, 0755)
	if err != nil {
		log.Error(err.Error())
	}
	file, err := AppFs.Create(path)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	data, err := PullFromUpstream(fmt.Sprintf("%s/v2/%s", upstream, path))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	go func(){
		file.Write(data)
		file.Close()
	}()

	return data, nil
}