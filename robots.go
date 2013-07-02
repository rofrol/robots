package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"runtime"
	"time"
)

type Location struct {
	Lat  float64
	Lng  float64
	Name string
}

func r_6043(c chan string) {
	for {
		fmt.Println(<-c)
	}
}

func dispatcher(c chan string) {
	db := openConn()
	defer db.Close()
	stmt, err := db.Prepare("select lat, lng, name from t_6043")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	location := Location{}
	for rows.Next() {
		rows.Scan(&location.Lat, &location.Lng, &location.Name)
		c <- fmt.Sprintln(location)
		time.Sleep(time.Second * 1)
	}
}

func mintime(table string) time.Time {
	db := openConn()
	defer db.Close()
	ts := time.Time{}
	err := db.QueryRow("select min(ts) from " + table).Scan(&ts)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return ts
}

func openConn() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres dbname=gisdb sslmode=disable password=droot")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	db.SetMaxIdleConns(100)
	return db
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	now := time.Now()
	min_6043 := mintime("t_6043")
	min_5937 := mintime("t_5937")
	min := min_6043
	if min_6043.After(min_5937) {
		min = min_5937
	}
	fmt.Println("now", now)
	fmt.Println("min_6043", min_6043)
	fmt.Println("min_5937", min_5937)
	fmt.Println("min     ", min)
	dur := now.Sub(min)
	fmt.Println("dur", dur)
	c := make(chan string)
	go dispatcher(c)
	go r_6043(c)
	var input string
	fmt.Scanln(&input)
}
