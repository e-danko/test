package common

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SystemConfiguration struct {
	MySql       string
	BaseDir     string
	ParsingMode string
}

var SystemConfig SystemConfiguration

func ReadConfig() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if strings.HasPrefix(dir, "/var") || strings.HasPrefix(dir, "/tmp") {
		dir, err = os.Getwd()
	}
	if err != nil {
		log.Fatal(err)
	}
	SystemConfig.BaseDir = dir

	if err := godotenv.Load(GetAbsPath(".env")); err != nil {
		log.Print("No .env file found")
	}

	dbName, dbNameExists := os.LookupEnv("DB_NAME")
	dbUser, dbUserExists := os.LookupEnv("DB_USER")
	dbPass, dbPassExists := os.LookupEnv("DB_PASS")
	dbHost, dbHostExists := os.LookupEnv("DB_HOST")
	if dbNameExists && dbUserExists && dbPassExists && dbHostExists {
		SystemConfig.MySql = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbName)
	}

	SystemConfig.ParsingMode = "on"
}
