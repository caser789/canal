```
mysqldump: [Warning] Using a password on the command line interface can be insecure.                                                                                                                                [118/1822]
WARNING: --master-data is deprecated and will be removed in a future version. Use --source-data instead.
Warning: A partial dump from a server that has GTIDs will by default include the GTIDs of all transactions, even those that changed suppressed parts of the database. If you don't want to restore GTIDs, pass --set-gtid-purg
ed=OFF. To make a complete dump, pass --all-databases --triggers --routines --events.
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;
SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '6ebab33e-abef-11ed-be38-0242ac120002:1-200';
CHANGE MASTER TO MASTER_LOG_FILE='binlog.000006', MASTER_LOG_POS=258629;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `test1` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `test1`;
INSERT INTO `t1` VALUES (1,'a');
INSERT INTO `t1` VALUES (2,'b');
INSERT INTO `t1` VALUES (3,'\\');
INSERT INTO `t1` VALUES (4,'\'\'');

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `test2` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `test2`;
INSERT INTO `t1` VALUES (1,'a');
INSERT INTO `t1` VALUES (2,'b');
INSERT INTO `t1` VALUES (3,'\\');
INSERT INTO `t1` VALUES (4,'\'\'');
INSERT INTO `t2` VALUES (1,'a');
INSERT INTO `t2` VALUES (2,'b');
INSERT INTO `t2` VALUES (3,'\\');
INSERT INTO `t2` VALUES (4,'\'\'');
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
```
