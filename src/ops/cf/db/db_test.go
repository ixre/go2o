package db

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	DB *Connector
)

type User struct {
	User string `db:"user" pk:"yes" auto:"no"`
	Pwd  string `db:"password"`
	Host string `db:"host"`
}

func Test_to(t *testing.T) {
	fmt.Println("==== testing query =====")
	DB.Query("SELECT user,host FROM mysql.user", func(rows *sql.Rows) {
		fmt.Println(ConvRowsToMapForJson(rows))

		rows.Close()
	})
	t.Log("11")
}

func Test_model(t *testing.T) {
	fmt.Println("===== testing model =======")
	var usr User
	DB.ORM.Get(&usr, "root")
	fmt.Println("User:" + usr.User)
	fmt.Println("Pwd:" + usr.Pwd)
	fmt.Println("Host:" + usr.Host)
}

func Test_Select(t *testing.T) {
	fmt.Println("===== testing select model =======")
	for i := 0; i < 3; i++ {
		var usrs []User
		DB.ORM.Select(&usrs, User{User: "root"}, "")
		if i == 0 {
			fmt.Println(usrs)
		}
	}
}

func Test_insermodel(t *testing.T) {

	fmt.Println("\n===== testing insert model =======")
	i, i2, err := DB.ORM.Save(nil, User{Host: "localhost", User: "uu1", Pwd: "1233455"})
	fmt.Println(i, i2, err)

	var usr User
	DB.ORM.Get(&usr, "uu1")
	fmt.Println("Inserted :", usr)

}

func Test_savemodel(t *testing.T) {
	fmt.Println("===== testing save model =======")
	var usr User
	DB.ORM.Get(&usr, "uu1")
	usr.Host = "127.0.0.1"
	_, _, err := DB.ORM.Save(usr.User, usr)
	if err != nil {
		fmt.Println("happend error:", err.Error())
	} else {
		DB.ORM.Get(&usr, "uu1")
		fmt.Println("updated host:", usr.Host)
	}

}

func Test_delmodel(t *testing.T) {
	fmt.Println("===== testing deleting model =======")
	i, err := DB.ORM.Delete(User{User: "uu1"}, "")
	fmt.Println(i, "rows deleted")
	if err != nil {
		fmt.Println("happend error:", err.Error())
	}
}

func init() {
	log.SetOutput(os.Stdout)
	DB = NewConnector("mysql",
		"root:@tcp(localhost:3306)/mysql?charset=utf8")
	DB.ORM.SetTrace(true)
	DB.ORM.CreateTableMap(User{}, "user")
}
