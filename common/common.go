package common

import (
	"fmt"
	_ "github.com/jinzhu/gorm"
	"log"
	"os"
	"strings"
)

var systemValues map[string]string

func Start() {
	systemValues = make(map[string]string)
	args := os.Args[1:]
	l := len(args)

	for i := 0; i < l; i++ {
		key := args[i][2:]
		if i+1 == l {
			systemValues[key] = "true"
			break
		}
		value := args[i+1]
		if value[0:2] == "--" {
			//some flag was set
			systemValues[key] = "true"
		} else {
			systemValues[key] = value
			i++
		}
	}
}

func GetAbsPath(path string) string {
	return SystemConfig.BaseDir + "/" + strings.TrimLeft(path, "/")
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
