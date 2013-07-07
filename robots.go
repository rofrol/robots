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

func robot(table string, c1 chan Location) {
	db := openConn()
	defer db.Close()
	stmt, err := db.Prepare("select lat, lng, name from " + table)
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
	for i := 0; rows.Next(); i++ {
		rows.Scan(&location.Lat, &location.Lng, &location.Name)
		c1 <- location
		if i%10 == 0 {
			time.Sleep(time.Second * 1)
		}
	}
}

func within(r1_lat float64, r1_lng float64, radius float64) {
	db := openConn()
	defer db.Close()

	var tube_name string
	q := fmt.Sprintf("SELECT name FROM tube WHERE ST_DWithin(ST_SetSRID(ST_MakePoint(%v, %v),4326, true), geom_4326,%v)", r1_lng, r1_lat, radius)
	err := db.QueryRow(q).Scan(&tube_name)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(r1_lat, r1_lng, tube_name)
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

	c_6043 := make(chan Location, 10)
	c_5937 := make(chan Location, 10)
	go robot("t_6043", c_6043)
	go robot("t_5937", c_6043)
	go func() {
		for {
			select {
			case msg1 := <-c_6043:
				fmt.Println(msg1)
				within(msg1.Lat, msg1.Lng, float64(150))
			case msg2 := <-c_5937:
				fmt.Println(msg2)
				within(msg2.Lat, msg2.Lng, float64(150))
			}
		}
	}()
	time.Sleep(time.Second * 5)
}
