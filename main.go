package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/go-ozzo/ozzo-validation"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	. "time"
	"google.golang.org/grpc/benchmark/stats"
	"math/rand"
	"os"
	"github.com/BurntSushi/toml"
)

var Histogram *stats.Histogram

func init () {
	var opts stats.HistogramOptions
	opts.NumBuckets = 7
	opts.MinValue = 0
	opts.BaseBucketSize = 1
	opts.GrowthFactor = 4
	Histogram = stats.NewHistogram(opts)

	Histogram.Buckets[0].LowBound = 0
	Histogram.Buckets[1].LowBound = 1
	Histogram.Buckets[2].LowBound = 5
	Histogram.Buckets[3].LowBound = 10
	Histogram.Buckets[4].LowBound = 20
	Histogram.Buckets[5].LowBound = 50
	Histogram.Buckets[6].LowBound = 100

}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", Event).Methods("POST")
	router.HandleFunc("/histogram", GetHistogram)
	log.Fatal(http.ListenAndServe(":8080", router))

}


type Model struct {
	API_KEY   string `json:"API_KEY"`
	USER_ID   string `json:"USER_ID"`
	TIMESTAMP string `json:"TIMESTAMP"`
}

func (m Model) Validate() error {
	return validation.ValidateStruct(&m,
		// API_KEY cannot be empty, and should be either 1, 2 or 3
		validation.Field(&m.API_KEY, validation.Required, validation.In("1", "2", "3")),
		// USER_ID cannot be empty
		validation.Field(&m.USER_ID, validation.Required),
		// TIMESTAMP cannot be empty
		validation.Field(&m.TIMESTAMP, validation.Required),
	)
}

func Event(w http.ResponseWriter, req *http.Request) {
	start := Now()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer req.Body.Close()
	// Unmarshal
	var msg Model
	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	errs := msg.Validate()
	if errs != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		config := ReadConfig2()
		db, err := sql.Open("mysql", config.Dblink)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer db.Close()
		// make sure connection is available
		err = db.Ping()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		stmt, err := db.Prepare("insert into event (api_key, body) values(?,?);")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_, err = stmt.Exec(msg.API_KEY, body)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintln(w, "Event is succesfully stored!")
		r := rand.Intn(100)
		Sleep(Millisecond * Duration(r))
		duration := Now().Sub(start).Seconds()
		Histogram.Add(int64(duration*1000))
	}

}

func GetHistogram(w http.ResponseWriter, req *http.Request) {
	Histogram.Print(w)
}

// Config ...
type Config struct {
	Dblink string
}

// Reads info from config file
func ReadConfig2() Config {
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
