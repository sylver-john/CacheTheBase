package main

import (
	"log"
	"net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "encoding/json"
    "os"
    "io"
    "io/ioutil"
)

type Data struct {
	Id int
	Text string
}

func main() {
	// here we use the cache file generated instead of doing a sql request
	http.HandleFunc("/", getData)
	log.Print("launch the server")
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("./cache/cache"); err == nil {
		r, wi := io.Pipe()
	    go func() {
	        defer wi.Close()
	        newFile, err := ioutil.ReadFile("./cache/cache")
	        if err != nil {
	            log.Fatal(err)
	        }
			json.NewEncoder(wi).Encode(newFile)
			json.NewEncoder(w).Encode(r)
	    }()
	    return
	}
	db, err := sql.Open("mysql","root:@tcp(127.0.0.1:3306)/base_test")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM info")
	if err != nil {
		log.Fatal(err)
	}
	var data []Data
    for rows.Next() {
    	row := new(Data)
        err = rows.Scan(&row.Id, &row.Text)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, *row)
    }
	json.NewEncoder(w).Encode(data)
}

