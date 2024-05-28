package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ProxyConfig struct {
	Name          string `yaml:"name"`
	ListenAddress string `yaml:"listenAddress"`
	Auth          string `yaml:"auth"`
	Session       struct {
		Name   string `yaml:"name"`
		Secret string `yaml:"secret"`
	} `yaml:"session"`
	UnAuthedResponse string `yaml:"unAuthedResponse"`
	Token            struct {
		Valid struct {
			URL      string `yaml:"url"`
			Format   string `yaml:"format"`
			JSONPath string `yaml:"jsonPath"`
			XMLPath  string `yaml:"xmlPath"`
		} `yaml:"valid"`
	} `yaml:"token"`
	Jwt struct {
		Valid struct {
			Secret   string `yaml:"secret"`
			JSONPath string `yaml:"jsonPath"`
		} `yaml:"valid"`
	} `yaml:"jwt"`
	Cas struct {
		EndPoint string `yaml:"endPoint"`
		XMLPath  string `yaml:"xmlPath"`
	} `yaml:"cas"`
	Reverse []struct {
		Type    string `yaml:"type"`
		Backend string `yaml:"backend"`
		Code    int    `yaml:"code"`
		To      string `yaml:"to"`
		Rewrite struct {
			Header []struct {
				Name string `yaml:"name"`
				From string `yaml:"from"`
				To   string `yaml:"to"`
			} `yaml:"header"`
		} `yaml:"rewrite"`
		URL []string `yaml:"url"`
	} `yaml:"reverse"`
}

var instances struct {
	Instances []ProxyConfig `yaml:"instances"`
}

var Config []ProxyConfig

func rebuildConfig(v *viper.Viper) {
	if err := v.Unmarshal(&instances); err != nil {
		log.Fatalln(err)
	}
	Config = instances.Instances
}

func Init(configFile string) {
	v := viper.New()
	v.SetConfigFile("./" + configFile)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln("viper read:", err)
	}
	rebuildConfig(v)

	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		rebuildConfig(v)
	})
	viper.WatchConfig()
}
