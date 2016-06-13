package main

import (
	"htz_data_analyzer/log"

	"github.com/BurntSushi/toml"
)

type config struct {
	ZhanZhuang    ZhanZhuang
	JingZuo       JingZuo
	KuanLiangMiao KuanLiangMiao
}

type ZhanZhuang struct{
    Pre []string
}

type JingZuo struct{
    Pre []string
}

type KuanLiangMiao struct{
    Pre []string
}

var (
	Conf config
)

func init() {
	_, err := toml.DecodeFile("./conf.toml", &Conf)
	if nil != err {
		log.Fatalln(err)
	}

//	log.Debugln("config meta:", meta)
	log.Debugln("config data:", Conf)
}
