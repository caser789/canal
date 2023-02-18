package main

import (
	"fmt"
)

type syncer struct {
	host     string
	user     string
	password string
}

func main() {
	fmt.Println("start main")
	fmt.Println("ping")

	c, err := Connect("127.0.0.1:3306", "root", "test", "testdb")
	if err != nil {
		fmt.Printf("Connect error %s\n", err)
		return
	}

	err = c.Ping()
	if err != nil {
		fmt.Printf("Ping error %s\n", err)
		return
	}

	fmt.Println("end main")
}
