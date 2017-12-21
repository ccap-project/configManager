package neo4j

import (
	"log"
	"os"

	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

type Conn driver.Conn

func GetConnectionString() string {

	url := os.Getenv("FOO")

	if len(url) <= 0 {
		url = "bolt://192.168.20.54:7687"
	}

	return url
}

func Connect(connStr string) (driver.Conn, error) {

	if len(connStr) <= 0 {
		connStr = GetConnectionString()
	}

	db, err := driver.NewDriver().OpenNeo(connStr)

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, err
	}

	return db, err
}
