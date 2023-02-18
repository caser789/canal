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