CREATE DATABASE IF NOT EXISTS test_backup
USE test_backup

RESET MASTER

```
=== RotateEvent ===
Date: 1970-01-01 07:30:00
Log position: 0
Event size: 40
Position: 4
Next log name: binlog.000001

rotate to (binlog.000001, 4)
=== FormatDescriptionEvent ===
Date: 2023-02-19 21:11:48
Log position: 126
Event size: 122
Version: 4
Server version: 8.0.32                                            
Checksum algorithm: 1

=== PreviousGTIDsEvent ===
Date: 2023-02-19 21:11:48
Log position: 157
Event size: 31
Event data:
00000000  00 00 00 00 00 00 00 00                           |........|
```

StartBackup
