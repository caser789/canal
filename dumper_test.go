package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func Test_unescapeString(t *testing.T) {
	tests := []struct {
		escaped  string
		expected string
	}{
		{`\\n`, `\n`},
		{`\\t`, `\t`},
		{`\\"`, `\"`},
		{`\\'`, `\'`},
		{`\\0`, `\0`},
		{`\\b`, `\b`},
		{`\\Z`, `\Z`},
		{`\\r`, `\r`},
		{`abc`, `abc`},
		{`abc\`, `abc`},
		{`ab\c`, `abc`},
		{`\abc`, `abc`},
	}

	for _, tt := range tests {
		t.Run("unescapeString", func(t *testing.T) {
			got := unescapeString(tt.escaped)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func Test_Dump(t *testing.T) {
	cfg := BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   MySQLFlavor,
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "test",
	}

	c, err := Connect(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.User, cfg.Password, "")
	assert.NoError(t, err)

	defer func() {
		_, err := c.Execute("DROP DATABASE IF EXISTS test1")
		assert.NoError(t, err)

		_, err = c.Execute("DROP DATABASE IF EXISTS test2")
		assert.NoError(t, err)

		c.Close()
	}()

	exe := "/usr/local/opt/mysql-client/bin/mysqldump"
	d, err := NewDumper(exe, fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.User, cfg.Password)
	assert.NoError(t, err)

	d.SetCharset("utf8")
	d.SetErrOut(os.Stderr)

	_, err = c.Execute("CREATE DATABASE IF NOT EXISTS test1")
	assert.NoError(t, err)

	_, err = c.Execute("CREATE DATABASE IF NOT EXISTS test2")
	assert.NoError(t, err)

	str := `CREATE TABLE IF NOT EXISTS test%d.t%d (
			id int AUTO_INCREMENT,
			name varchar(256),
			PRIMARY KEY(id)
			) ENGINE=INNODB`
	_, err = c.Execute(fmt.Sprintf(str, 1, 1))
	assert.NoError(t, err)

	_, err = c.Execute(fmt.Sprintf(str, 2, 1))
	assert.NoError(t, err)

	_, err = c.Execute(fmt.Sprintf(str, 1, 2))
	assert.NoError(t, err)

	_, err = c.Execute(fmt.Sprintf(str, 2, 2))
	assert.NoError(t, err)

	str = `INSERT INTO test%d.t%d (name) VALUES ("a"), ("b"), ("\\"), ("''")`

	_, err = c.Execute(fmt.Sprintf(str, 1, 1))
	assert.NoError(t, err)

	_, err = c.Execute(fmt.Sprintf(str, 2, 1))
	assert.NoError(t, err)

	_, err = c.Execute(fmt.Sprintf(str, 1, 2))
	assert.NoError(t, err)

	_, err = c.Execute(fmt.Sprintf(str, 2, 2))
	assert.NoError(t, err)

	d.AddDatabases("test1", "test2")
	d.AddIgnoreTables("test1", "t2")
	err = d.Dump(os.Stdout)
	assert.NoError(t, err)

	d.AddTables("test1", "t1")
	err = d.Dump(io.Discard)
	assert.NoError(t, err)

	var buf bytes.Buffer
	d.Reset()
	d.AddDatabases("test1", "test2")
	err = d.Dump(&buf)
	assert.NoError(t, err)

	err = Parse(&buf, new(testParseHandler), true)
	assert.NoError(t, err)
}

type testParseHandler struct {
	gset GTIDSet
}

func (h *testParseHandler) BinLog(name string, pos uint64) error {
	fmt.Printf("binlog intput %s %d\n", name, pos)
	return nil
}

func (h *testParseHandler) GtidSet(gtidsets string) (err error) {
	if h.gset != nil {
		err = h.gset.Update(gtidsets)
	} else {
		h.gset, err = ParseGTIDSet("mysql", gtidsets)
	}

	fmt.Printf("gset %v\n", h.gset)
	return err
}

func (h *testParseHandler) Data(schema string, table string, values []string) error {
	fmt.Printf("data %s %s %s\n", schema, table, values)
	return nil
}
