package main

import(
	"log"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "os"
    "io"
)

type Data struct {
	Id int
	Text string
}

func main() {
	// here we make the sql request and put the result in a file
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
    file, err := os.Create("./cache/cache")
    if err != nil {
        log.Fatal(err)
    }
    _, err = io.Copy(file, data)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()
}