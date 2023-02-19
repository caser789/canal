package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/semver"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

func _TestMysqlGTIDSync(t *testing.T) {
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

	_, err = c.Execute("CREATE DATABASE IF NOT EXISTS test_backup")
	assert.NoError(t, err)

	_, err = c.Execute("USE test_backup")
	assert.NoError(t, err)

	r, err := c.Execute("SELECT @@gtid_mode")
	assert.NoError(t, err)
	modeOn, _ := r.GetString(0, 0)
	fmt.Println(r.Values)
	fmt.Printf("gtid_mode %s\n", modeOn)

	r, err = c.Execute("SHOW GLOBAL VARIABLES LIKE 'SERVER_UUID'")
	assert.NoError(t, err)
	fmt.Println(r.Values)
	s, _ := r.GetString(0, 1)
	assert.Less(t, 0, len(s))
	assert.NotEqual(t, s, "NONE")
	var masterUuid uuid.UUID
	masterUuid, err = uuid.Parse(s)
	assert.NoError(t, err)
	fmt.Printf("uuid is %s\n", masterUuid.String())

	set, _ := ParseMySQLGTIDSet(fmt.Sprintf("%s:%d-%d", masterUuid.String(), 1, 2))

	b := NewBinlogSyncer(cfg)
	var ss *BinlogStreamer
	ss, err = b.StartSyncGTID(set)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()

		for {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			e, err := ss.GetEvent(ctx)
			cancel()

			if err == context.DeadlineExceeded {
				return
			}

			assert.NoError(t, err)
			e.Dump(os.Stdout)
			os.Stdout.Sync()
		}
	}()

	testSync(t, c)
}

func _TestMysqlPositionSync(t *testing.T) {
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

	_, err = c.Execute("CREATE DATABASE IF NOT EXISTS test_backup")
	assert.NoError(t, err)

	_, err = c.Execute("USE test_backup")
	assert.NoError(t, err)

	r, err := c.Execute("SHOW MASTER STATUS")
	assert.NoError(t, err)
	binFile, _ := r.GetString(0, 0)
	binPos, _ := r.GetInt(0, 1)
	fmt.Printf("bin file %s\n", binFile)
	fmt.Printf("bin pos %d\n", binPos)

	b := NewBinlogSyncer(cfg)
	s, err := b.StartSync(Position{Name: binFile, Pos: uint32(binPos)})
	assert.NoError(t, err)

	r, err = c.Execute("SHOW SLAVE HOSTS")
	assert.NoError(t, err)
	fmt.Printf("values %v\n", r.Values)

	// Slave_UUID is empty for mysql 8.0.28+ (8.0.32 still broken)
	if semver.Compare(c.GetServerVersion(), "8.0.28") < 0 {
		// check we have set Slave_UUID
		slaveUUID, _ := r.GetString(0, 4)
		assert.Equal(t, 36, len(slaveUUID))
	}

	// Test re-sync.
	time.Sleep(100 * time.Millisecond)
	_ = b.c.SetReadDeadline(time.Now().Add(time.Millisecond))
	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()

		for {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			e, err := s.GetEvent(ctx)
			cancel()

			if err == context.DeadlineExceeded {
				return
			}

			assert.NoError(t, err)
			e.Dump(os.Stdout)
			os.Stdout.Sync()
		}
	}()

	testSync(t, c)
}

func _Test_BinlogSyncer_ParseFile(t *testing.T) {
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

	_, err = c.Execute("CREATE DATABASE IF NOT EXISTS test_backup")
	assert.NoError(t, err)

	_, err = c.Execute("USE test_backup")
	assert.NoError(t, err)

	b := NewBinlogSyncer(cfg)

	_, err = c.Execute("RESET MASTER")
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()

		testSync(t, c)

		_, err = c.Execute("FLUSH LOGS")
		assert.NoError(t, err)

		testSync(t, c)
	}()

	binlogDir := "./var"
	os.RemoveAll(binlogDir)

	err = b.StartBackup(binlogDir, Position{Name: "", Pos: uint32(0)}, time.Second*2)
	assert.NoError(t, err)

	p := NewBinlogParser()
	p.SetVerifyChecksum(true)

	f := func(e *BinlogEvent) error {
		e.Dump(os.Stdout)
		os.Stdout.Sync()
		return nil
	}

	dir, err := os.Open(binlogDir)
	assert.NoError(t, err)
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	assert.NoError(t, err)

	for _, file := range files {
		err = p.ParseFile(path.Join(binlogDir, file), 0, f)
		assert.NoError(t, err)
	}
}

func _Test_BinlogSyncer_StartBackup(t *testing.T) {
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

	_, err = c.Execute("CREATE DATABASE IF NOT EXISTS test_backup")
	assert.NoError(t, err)

	_, err = c.Execute("USE test_backup")
	assert.NoError(t, err)

	b := NewBinlogSyncer(cfg)

	_, err = c.Execute("RESET MASTER")
	assert.NoError(t, err)

	for times := 1; times <= 2; times++ {
		testSync(t, c)
		_, err = c.Execute("FLUSH LOGS")
		assert.NoError(t, err)
	}

	binlogDir := "./var"
	os.RemoveAll(binlogDir)
	timeout := 2 * time.Second
	done := make(chan bool)

	go func() {
		err = b.StartBackup(binlogDir, Position{Name: "", Pos: uint32(0)}, timeout)
		assert.NoError(t, err)
		done <- true
	}()
	failTimeout := 5 * timeout
	ctx, _ := context.WithTimeout(context.Background(), failTimeout)
	select {
	case <-done:
		return
	case <-ctx.Done():
		assert.NoError(t, errors.New("time out error"))
	}
}

func testSync(t *testing.T, c *Conn) {
	var err error
	var str string
	// use mixed format

	_, err = c.Execute("SET SESSION binlog_format = 'MIXED'")
	assert.NoError(t, err)

	_, err = c.Execute("DROP TABLE IF EXISTS test_replication")
	assert.NoError(t, err)
	str = `CREATE TABLE test_replication (
			id BIGINT(64) UNSIGNED  NOT NULL AUTO_INCREMENT,
			str VARCHAR(256),
			f FLOAT,
			d DOUBLE,
			de DECIMAL(10,2),
			i INT,
			bi BIGINT,
			e enum ("e1", "e2"),
			b BIT(8),
			y YEAR,
			da DATE,
			ts TIMESTAMP,
			dt DATETIME,
			tm TIME,
			t TEXT,
			bb BLOB,
			se SET('a', 'b', 'c'),
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`
	_, err = c.Execute(str)
	assert.NoError(t, err)
	// use row format
	_, err = c.Execute("SET SESSION binlog_format = 'ROW'")
	assert.NoError(t, err)
	str = `INSERT INTO test_replication (str, f, i, e, b, y, da, ts, dt, tm, de, t, bb, se)
		VALUES ("3", -3.14, 10, "e1", 0b0011, 1985,
		"2012-05-07", "2012-05-07 14:01:01", "2012-05-07 14:01:01",
		"14:01:01", -45363.64, "abc", "12345", "a,b")`
	_, err = c.Execute(str)
	assert.NoError(t, err)
	id := 100
	_, err = c.Execute("SET SESSION binlog_row_image = 'MINIMAL'")
	assert.NoError(t, err)
	_, err = c.Execute(fmt.Sprintf(`INSERT INTO test_replication (id, str, f, i, bb, de) VALUES (%d, "4", -3.14, 100, "abc", -45635.64)`, id))
	assert.NoError(t, err)
	_, err = c.Execute(fmt.Sprintf(`UPDATE test_replication SET f = -12.14, de = 555.34 WHERE id = %d`, id))
	assert.NoError(t, err)
	_, err = c.Execute(fmt.Sprintf(`DELETE FROM test_replication WHERE id = %d`, id))
	assert.NoError(t, err)

	// check whether we can create the table including the json field
	str = `DROP TABLE IF EXISTS test_json`
	_, err = c.Execute(str)
	assert.NoError(t, err)
	str = `CREATE TABLE test_json (
			id BIGINT(64) UNSIGNED  NOT NULL AUTO_INCREMENT,
			c1 JSON,
			c2 DECIMAL(10, 0),
			PRIMARY KEY (id)
			) ENGINE=InnoDB`
	_, err = c.Execute(str)
	assert.NoError(t, err)
	_, err = c.Execute(`INSERT INTO test_json (c2) VALUES (1)`)
	assert.NoError(t, err)
	_, err = c.Execute(`INSERT INTO test_json (c1, c2) VALUES ('{"key1": "value1", "key2": "value2"}', 1)`)
	assert.NoError(t, err)

	_, err = c.Execute(`DROP TABLE IF EXISTS test_json_v2`)
	assert.NoError(t, err)
	str = `CREATE TABLE test_json_v2 (
			id INT, 
			c JSON, 
			PRIMARY KEY (id)
			) ENGINE=InnoDB`
	_, err = c.Execute(str)
	assert.NoError(t, err)

	tbls := []string{
		// Refer: https://github.com/shyiko/mysql-binlog-connector-java/blob/c8e81c879710dc19941d952f9031b0a98f8b7c02/src/test/java/com/github/shyiko/mysql/binlog/event/deserialization/json/JsonBinaryValueIntegrationTest.java#L84
		// License: https://github.com/shyiko/mysql-binlog-connector-java#license
		`INSERT INTO test_json_v2 VALUES (0, NULL)`,
		`INSERT INTO test_json_v2 VALUES (1, '{\"a\": 2}')`,
		`INSERT INTO test_json_v2 VALUES (2, '[1,2]')`,
		`INSERT INTO test_json_v2 VALUES (3, '{\"a\":\"b\", \"c\":\"d\",\"ab\":\"abc\", \"bc\": [\"x\", \"y\"]}')`,
		`INSERT INTO test_json_v2 VALUES (4, '[\"here\", [\"I\", \"am\"], \"!!!\"]')`,
		`INSERT INTO test_json_v2 VALUES (5, '\"scalar string\"')`,
		`INSERT INTO test_json_v2 VALUES (6, 'true')`,
		`INSERT INTO test_json_v2 VALUES (7, 'false')`,
		`INSERT INTO test_json_v2 VALUES (8, 'null')`,
		`INSERT INTO test_json_v2 VALUES (9, '-1')`,
		`INSERT INTO test_json_v2 VALUES (10, CAST(CAST(1 AS UNSIGNED) AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (11, '32767')`,
		`INSERT INTO test_json_v2 VALUES (12, '32768')`,
		`INSERT INTO test_json_v2 VALUES (13, '-32768')`,
		`INSERT INTO test_json_v2 VALUES (14, '-32769')`,
		`INSERT INTO test_json_v2 VALUES (15, '2147483647')`,
		`INSERT INTO test_json_v2 VALUES (16, '2147483648')`,
		`INSERT INTO test_json_v2 VALUES (17, '-2147483648')`,
		`INSERT INTO test_json_v2 VALUES (18, '-2147483649')`,
		`INSERT INTO test_json_v2 VALUES (19, '18446744073709551615')`,
		`INSERT INTO test_json_v2 VALUES (20, '18446744073709551616')`,
		`INSERT INTO test_json_v2 VALUES (21, '3.14')`,
		`INSERT INTO test_json_v2 VALUES (22, '{}')`,
		`INSERT INTO test_json_v2 VALUES (23, '[]')`,
		`INSERT INTO test_json_v2 VALUES (24, CAST(CAST('2015-01-15 23:24:25' AS DATETIME) AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (25, CAST(CAST('23:24:25' AS TIME) AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (125, CAST(CAST('23:24:25.12' AS TIME(3)) AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (225, CAST(CAST('23:24:25.0237' AS TIME(3)) AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (26, CAST(CAST('2015-01-15' AS DATE) AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (27, CAST(TIMESTAMP'2015-01-15 23:24:25' AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (127, CAST(TIMESTAMP'2015-01-15 23:24:25.12' AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (227, CAST(TIMESTAMP'2015-01-15 23:24:25.0237' AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (327, CAST(UNIX_TIMESTAMP('2015-01-15 23:24:25') AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (28, CAST(ST_GeomFromText('POINT(1 1)') AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (29, CAST('[]' AS CHAR CHARACTER SET 'ascii'))`,
		// TODO: 30 and 31 are BIT type from JSON_TYPE, may support later.
		`INSERT INTO test_json_v2 VALUES (30, CAST(x'cafe' AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (31, CAST(x'cafebabe' AS JSON))`,
		`INSERT INTO test_json_v2 VALUES (100, CONCAT('{\"', REPEAT('a', 64 * 1024 - 1), '\":123}'))`,
	}
	for _, query := range tbls {
		_, err = c.Execute(query)
		assert.NoError(t, err)
	}

	// If MySQL supports JSON, it must supports GEOMETRY.
	_, err = c.Execute("DROP TABLE IF EXISTS test_geo")
	assert.NoError(t, err)
	str = `CREATE TABLE test_geo (g GEOMETRY)`
	_, err = c.Execute(str)
	assert.NoError(t, err)
	tbls = []string{
		`INSERT INTO test_geo VALUES (POINT(1, 1))`,
		`INSERT INTO test_geo VALUES (LINESTRING(POINT(0,0), POINT(1,1), POINT(2,2)))`,
		// TODO: add more geometry tests
	}
	for _, query := range tbls {
		_, err = c.Execute(query)
		assert.NoError(t, err)
	}

	str = `DROP TABLE IF EXISTS test_parse_time`
	_, err = c.Execute(str)
	assert.NoError(t, err)
	// Must allow zero time.
	_, err = c.Execute(`SET sql_mode=''`)
	assert.NoError(t, err)
	str = `CREATE TABLE test_parse_time (
			a1 DATETIME, 
			a2 DATETIME(3), 
			a3 DATETIME(6), 
			b1 TIMESTAMP, 
			b2 TIMESTAMP(3) , 
			b3 TIMESTAMP(6))`
	_, err = c.Execute(str)
	assert.NoError(t, err)
	str = `INSERT INTO test_parse_time VALUES
		("2014-09-08 17:51:04.123456", "2014-09-08 17:51:04.123456", "2014-09-08 17:51:04.123456", 
		"2014-09-08 17:51:04.123456","2014-09-08 17:51:04.123456","2014-09-08 17:51:04.123456"),
		("0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000",
		"0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000"),
		("2014-09-08 17:51:04.000456", "2014-09-08 17:51:04.000456", "2014-09-08 17:51:04.000456", 
		"2014-09-08 17:51:04.000456","2014-09-08 17:51:04.000456","2014-09-08 17:51:04.000456")`
	_, err = c.Execute(str)
	assert.NoError(t, err)
}
