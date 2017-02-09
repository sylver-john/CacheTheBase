# CacheTheBase
A Cache file generator for mysql

## First we launch the api with :
``
go run api.go
``

which start a web server on the port 1234, will listen to / with :
```go
	http.HandleFunc("/", getData)
	log.Print("launch the server")
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
```

execute a sql request in the handler with :
```go 
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
```

and return the result with :
```go
	json.NewEncoder(w).Encode(data)
```

## Then we run the cache file generator with :
``
go run fileGenerator.go
`` 
(which could be a cron job).It will execute the sql request and generate a file with the sql result as content with :
```go
	var data []Data
    	for rows.Next() {
    		row := new(Data)
        	err = rows.Scan(&row.Id, &row.Text)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, *row)
    	}
	dataJson, _ := json.Marshal(data)
 	err = ioutil.WriteFile("./cache/cache.json", dataJson, 0644)
```

## And finally in the Api Handler we use the json instead of doing a SQL request with :
```go
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
```
