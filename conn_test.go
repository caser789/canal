package mysql

import "testing"

func newTestConn() *conn {
	c := new(conn)

	if err := c.Connect("127.0.0.1:3306", "root", "test", "testdb"); err != nil {
		panic(err)
	}

	return c
}

func TestConn_Connect(t *testing.T) {
	c := newTestConn()
	defer c.Close()
}
