# canal
## Setup MySQL
* docker run
  - [install mysql with docker](https://medium.com/@philipp.schmiedel/local-docker-mysql-macos-fa7ac14348c4)
  - [mycli](https://www.mycli.net/install)
```
mycli -h localhost --protocol=TCP -P 3306 -u root test_db
```
* set binlog format = row
## task to insert dummy data

## Features
1. connect db
2. ping
3. conn execute
4. statement execute
5. begin/rollback/commit transaction

## notes
### binlog_row_image
* full
* minimal
### FLUSH LOGS
相当于运行这些命令，close and reopen
* FLUSH BINARY LOGS
* FLUSH RELAY LOGS
* FLUSH ENGINE LOGS
* FLUSH ERROR LOGS
* FLUSH GENERAL LOGS
* FLUSH SLOW LOGS


## CMD

### SHOW MASTER STATUS

```
+---------------+----------+--------------+------------------+-------------------+
| File          | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+---------------+----------+--------------+------------------+-------------------+
| binlog.000004 | 12714    |              |                  |                   |
+---------------+----------+--------------+------------------+-------------------+
```
### SHOW SLAVE HOSTS

```
+-----------+------+------+-----------+------------+
| Server_id | Host | Port | Master_id | Slave_UUID |
+-----------+------+------+-----------+------------+
+-----------+------+------+-----------+------------+
```
### RESET MASTER

For a server where binary logging is enabled (log_bin is ON), RESET MASTER deletes all existing binary log files and resets the binary log index file, resetting the server to its state before binary logging was started. A new empty binary log file is created so that binary logging can be restarted.

### rpl_semi_sync_master_enabled

Set to ON to enable semi-synchronous replication primary. Disabled by default.

### rpl_semi_sync_slave

* Controls how the server should treat the plugin when the server starts up.

#### values
* OFF - Disables the plugin without removing it from the mysql.plugins table.
* ON - Enables the plugin. If the plugin cannot be initialized, then the server will still continue starting up, but the plugin will be disabled.
* FORCE - Enables the plugin. If the plugin cannot be initialized, then the server will fail to start with an error.
* FORCE_PLUS_PERMANENT - Enables the plugin. If the plugin cannot be initialized, then the server will fail to start with an error. In addition, the plugin cannot be uninstalled with UNINSTALL SONAME or UNINSTALL PLUGIN while the server is running.

### SET NAMES
* [doc](https://dev.mysql.com/doc/refman/8.0/en/set-names.html)
* This statement sets the three session system variables character_set_client, character_set_connection, and character_set_results to the given character set.