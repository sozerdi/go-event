package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	"github.com/BurntSushi/toml"
)

func main() {
	config := ReadConfig()
	db, err := sql.Open("mysql", config.Dblink)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	stmt, err := db.Prepare("CREATE TABLE event (id int NOT NULL AUTO_INCREMENT, api_key varchar(15), body json, PRIMARY KEY (id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Event Table successfully migrated....")
	}
}

// Config ...
type Config struct {
	Dblink string
}

// Reads info from config file
func ReadConfig() Config {
	var configfile = "properties.ini"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return config
}
