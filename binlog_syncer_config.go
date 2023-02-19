package main

import (
	"crypto/tls"
	"time"
)

// BinlogSyncerConfig is the configuration for BinlogSyncer.
type BinlogSyncerConfig struct {
	// ServerID is the unique ID in cluster.
	ServerID uint32
	// Flavor is "mysql" or "mariadb", if not set, use "mysql" default.
	Flavor string
	// Charset is for MySQL client character set
	Charset string
	// RecvBufferSize sets the size in bytes of the operating system's receive buffer associated with the connection.
	RecvBufferSize int

	// Host is for MySQL server host.
	Host string
	// Port is for MySQL server port.
	Port uint16
	// User is for MySQL user.
	User string
	// Password is for MySQL password.
	Password string
	// Localhost is local hostname if register salve.
	// If not set, use os.Hostname() instead.
	Localhost string

	DiscardGTIDSet bool

	// SemiSyncEnabled enables semi-sync or not.
	SemiSyncEnabled bool
	// DumpCommandFlag is used to send binglog dump command. Default 0, aka BINLOG_DUMP_NEVER_STOP.
	// For MySQL, BINLOG_DUMP_NEVER_STOP and BINLOG_DUMP_NON_BLOCK are available.
	// https://dev.mysql.com/doc/internals/en/com-binlog-dump.html#binlog-dump-non-block
	// For MariaDB, BINLOG_DUMP_NEVER_STOP, BINLOG_DUMP_NON_BLOCK and BINLOG_SEND_ANNOTATE_ROWS_EVENT are available.
	// https://mariadb.com/kb/en/library/com_binlog_dump/
	// https://mariadb.com/kb/en/library/annotate_rows_event/

	DumpCommandFlag uint16

	// whether disable re-sync for broken connection
	DisableRetrySync bool

	// maximum number of attempts to re-establish a broken connection, zero or negative number means infinite retry.
	// this configuration will not work if DisableRetrySync is true
	MaxReconnectAttempts int

	// read timeout
	ReadTimeout time.Duration

	// master heartbeat period
	HeartbeatPeriod time.Duration

	//Option function is used to set outside of BinlogSyncerConfig， between mysql connection and COM_REGISTER_SLAVE
	//For MariaDB: slave_gtid_ignore_duplicates、skip_replication、slave_until_gtid
	Option func(*Conn) error
	// Set Dialer
	Dialer Dialer
	// If not nil, use the provided tls.Config to connect to the database using TLS/SSL.
	TLSConfig *tls.Config
}
