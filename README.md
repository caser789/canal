# canal
## Setup MySQL
* docker run
  - [install mysql with docker](https://medium.com/@philipp.schmiedel/local-docker-mysql-macos-fa7ac14348c4)
  - [mycli](https://www.mycli.net/install)
```
mycli -h localhost --protocol=TCP -P 3303 -u root test_db
```
* set binlog format = row
## task to insert dummy data
