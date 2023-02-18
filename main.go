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

	var (
		r *Result
	)
	r, err = c.Execute("CREATE DATABASE IF NOT EXISTS a")
	fmt.Println("create database")
	fmt.Println(err)
	fmt.Println(r)

	r, err = c.Execute("use a")
	fmt.Println("use database")
	fmt.Println(err)
	fmt.Println(r)

	r, err = c.Execute("CREATE TABLE IF NOT EXISTS a_tab (id bigint unsigned not null auto_increment, primary key (id))")
	fmt.Println("create table")
	fmt.Println(err)
	fmt.Println(r)

	str := `CREATE TABLE IF NOT EXISTS mixer_test_conn (
          id BIGINT(64) UNSIGNED  NOT NULL,
          str VARCHAR(256),
          f DOUBLE,
          e enum("test1", "test2"),
          u tinyint unsigned,
          i tinyint,
          j json,
          PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8`

	r, err = c.Execute(str)
	fmt.Println("create table")
	fmt.Println(err)
	fmt.Println(r)

	str = `insert into mixer_test_conn (id, str, f, e) values(1, "a", 3.14, "test1")`
	r, err = c.Execute(str)
	fmt.Println("insert")
	fmt.Println(err)
	fmt.Println(r)

	str = `select str, f, e from mixer_test_conn where id = 1`
	r, err = c.Execute(str)
	fmt.Println("select")
	fmt.Println(err)
	fmt.Println(r.Fields)
	fmt.Println(r.Values)

	str = fmt.Sprintf(`insert into mixer_test_conn (id, str) values(5, "%s")`, "abc")
	r, err = c.Execute(str)
	str = `select str from mixer_test_conn where id = ?`
	r, err = c.Execute(str, 5)
	fmt.Println("select with args")
	fmt.Println(err)
	fmt.Println(r.GetString(0, 0))

	str = `insert into mixer_test_conn (id, i) values (?, ?)`
	stmt, err := c.Prepare(str)
	fmt.Println("prepare")
	fmt.Println(err)
	fmt.Println(stmt)
	defer stmt.Close()

	_, err = stmt.Execute(4, 127)
	_, err = stmt.Execute(uint64(18446744073709551516), int8(-128))

	fmt.Println("end main")

}
