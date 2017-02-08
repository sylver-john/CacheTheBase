package main

import (
	"log"
	"net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "encoding/json"
    "os"
    "io/ioutil"
)

type Data struct {
	Id int `json:Id`
	Text string `json:Text`
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
	if _, err := os.Stat("./cache/cache.json"); err == nil {
		log.Print("use cache file")
		file, err := ioutil.ReadFile("./cache/cache.json")
		if err != nil {
			log.Fatal(err)
		}
		var data *[]Data
   		err = json.Unmarshal(file, &data)
		json.NewEncoder(w).Encode(data)
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

