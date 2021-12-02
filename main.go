package main

import (
	"parsers/common"
	"parsers/database"
	"parsers/worker"
)

func main() {
	common.ReadConfig()
	database.Start()
	worker.Start()
}
