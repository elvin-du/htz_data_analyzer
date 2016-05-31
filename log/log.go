package log

import (
	"github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//func init() {
////	log.SetFormatter(new(prefixed.TextFormatter))
//	log.SetLevel(log.DebugLevel)
//}

var log = logrus.New()

func init() {
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel
}

func Infoln(v ...interface{}) {
	log.Infoln(v...)
}

func Debugln(v ...interface{}) {
	log.Debugln(v...)
}

func Warningln(v ...interface{}) {
	log.Warnln(v...)
}

func Errorln(v ...interface{}) {
	log.Errorln(v...)
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}
