package main

import (
	"fmt"
)

type syncer struct {
	host     string
	user     string
	password string
}

func (s *syncer) Dump() {
}

func main() {
	fmt.Println("start main")
	host := "127.0.0.1:3306"
	user := "root"
	password := "test"
	d := syncer{
		host:     host,
		user:     user,
		password: password,
	}
	d.Dump()
	fmt.Println("end main")
}
