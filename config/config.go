package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config struct {
	ListenAddress string `yaml:"listenAddress"`
	Proxy         struct {
		Name    string `yaml:"name"`
		Auth    string `yaml:"auth"`
		Session struct {
			Name   string `yaml:"name"`
			Secret string `yaml:"secret"`
		} `yaml:"session"`
		UnAuthedResponse string      `yaml:"unAuthedResponse"`
		Jwt              interface{} `yaml:"jwt"`
		Token            struct {
			Valid struct {
				URL      string `yaml:"url"`
				Format   string `yaml:"format"`
				JSONPath string `yaml:"jsonPath"`
				XMLPath  string `yaml:"xmlPath"`
			} `yaml:"valid"`
		} `yaml:"token"`
		Cas struct {
			EndPoint string `yaml:"endPoint"`
			XMLPath  string `yaml:"id"`
		} `yaml:"cas"`
		Reverse struct {
			Backend string `yaml:"backend"`
			Rewrite struct {
				Header []struct {
					Name string `yaml:"name"`
					From string `yaml:"from"`
					To   string `yaml:"to"`
				} `yaml:"header"`
			} `yaml:"rewrite"`
			URL []string `yaml:"url"`
		} `yaml:"reverse"`
	} `yaml:"proxy"`
}

func rebuildConfig(v *viper.Viper) {
	if err := v.Unmarshal(&Config); err != nil {
		log.Fatalln(err)
	}
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
