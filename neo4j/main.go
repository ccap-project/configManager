/*
 *
 * Copyright (c) 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package neo4j

import (
	"fmt"
	"log"
	"os"

	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

type ConnPool driver.DriverPool
type Conn driver.Conn

func GetConnectionString() string {

	url := os.Getenv("CONFIGMANAGER_DB_HOST")

	if len(url) <= 0 {
		url = "bolt://192.168.20.54:7687"
	}

	return url
}

/*
func Connect(host string, port string, user string, passwd string) (driver.Conn, error) {

	var connStr string

	if len(host) <= 0 {
		port = "127.0.0.1"
	}
	if len(port) <= 0 {
		port = "7687"
	}

	if len(user) <= 0 {
		connStr = fmt.Sprintf("bolt://%s:%s", host, port)
	} else {
		connStr = fmt.Sprintf("bolt://%s:%s@%s:%s", user, passwd, host, port)
	}
*/

func Connect(connStr string) (driver.Conn, error) {

	log.Printf("Connecting to %s", connStr)

	db, err := driver.NewDriver().OpenNeo(connStr)

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, err
	}

	return db, err
}
func Pool(host string, port string, user string, passwd string, max_conn int) (driver.DriverPool, error) {

	var connStr string

	if len(host) <= 0 {
		port = "127.0.0.1"
	}
	if len(port) <= 0 {
		port = "7687"
	}

	if max_conn <= 0 {
		max_conn = 10
	}

	if len(user) <= 0 {
		connStr = fmt.Sprintf("bolt://%s:%s", host, port)
	} else {
		connStr = fmt.Sprintf("bolt://%s:%s@%s:%s", user, passwd, host, port)
	}

	log.Printf("Connecting to %s", connStr)

	pool, err := driver.NewDriverPool(connStr, max_conn)

	if err != nil {
		log.Println("error creating neo4j connection pool:", err)
		return nil, err
	}

	return pool, err
}
