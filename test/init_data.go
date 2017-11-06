package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon/complex"
	"github.com/m0cchi/gfalcon/model"
	"os"
)

func main() {
	var dbhost string
	flag.StringVar(&dbhost, "dbhost", "", "gfalcon's DB")
	flag.Parse()

	if dbhost == "" {
		fmt.Println("required --dbhost [host]")
		os.Exit(1)
	}

	db, err := sqlx.Connect("mysql", dbhost)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	user, _ := model.GetUser(db, 1, "gfadmin")
	fmt.Println(user.UpdatePassword(db, "secret"))
	session, err := complex.AuthenticateWithPassword(db, user, "secret")
	if session.Validate() == nil {
		fmt.Println("update password: success")
	} else {
		fmt.Println("update password: fail")
	}
}
