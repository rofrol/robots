package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

type Location struct {
	Lat  float64
	Lng  float64
	Name string
}

var traffic = [...]string{"HEAVY", "LIGHT", "MODERATE"}

var db *sql.DB

func robot(table string, c1 chan Location) {
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

func within(r Location, radius float64) {
	q := fmt.Sprintf(`
SELECT ST_Y(t.geom_4326), ST_X(t.geom_4326), t.name, 
ST_Distance_Sphere(ST_SetSRID(ST_Point(%v, %v),4326), t.geom_4326)
FROM tube t join t_6043 r
on ST_DWithin(ST_SetSRID(ST_Point(%v, %v),4326), t.geom_4326, %v, true)
`, r.Lng, r.Lat, r.Lng, r.Lat, radius)

	//fmt.Println(q)
	loc := Location{}
	var dist float64
	err := db.QueryRow(q).Scan(&loc.Lat, &loc.Lng, &loc.Name, &dist)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
		os.Exit(1)
	} else if err != sql.ErrNoRows {
		fmt.Println(time.Now().Format("15:04:05"), r.Name, loc.Name, traffic[rand.Intn(3)], int(dist))
	}
}

func openConn() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=gisdb sslmode=disable password=droot")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	db.SetMaxIdleConns(100)

	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	openConn()
	defer db.Close()

	c_6043 := make(chan Location, 10)
	c_5937 := make(chan Location, 10)
	go robot("t_6043", c_6043)
	go robot("t_5937", c_6043)
	go func() {
		for {
			select {
			case msg1 := <-c_6043:
				within(msg1, float64(350))
			case msg2 := <-c_5937:
				within(msg2, float64(350))
			}
		}
	}()
	time.Sleep(time.Second * 500)
}
