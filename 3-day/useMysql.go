package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

package main

import (
"database/sql"
"fmt"
"io/ioutil"
"os"
"strings"

_ "github.com/go-sql-driver/mysql"
)

func readDbConfig() (id string, pw string, host string, dbName string, port string) {
	temp, err := ioutil.ReadFile("./db.cnf")
	if err != nil {
		os.Exit(1)
	}
	valBuffer := strings.Split(string(temp), "\n")
	id = valBuffer[0]
	pw = valBuffer[1]
	host = valBuffer[2]
	dbName = valBuffer[3]
	port = valBuffer[4]
	return
}
func main() {
	id, pw, host, dbName, _ := readDbConfig()
	value := id + ":" + pw + "@tcp(" + host + ")/" + dbName
	db, err := sql.Open("mysql", value)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result string
	rows, err := db.Query("show databases")
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result)
	}
	defer db.Close()

}
