package database

import (
	"fmt"
	"parsers/common"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Connection *gorm.DB
var err error
var DateFormat = "2006-01-02"
var DateTimeFormat = "2006-01-02 15:04:05"

func Start() bool {
	if url := common.SystemConfig.MySql; url != "" {
		Connection, err = gorm.Open("mysql", url)
		if err != nil {
			common.FailOnError(err, "Failed to connect to database")
			return false
		}

		fmt.Println("Connected to database!")
		return true
	} else {
		fmt.Println("No database mode is ON")
		return false
	}
}