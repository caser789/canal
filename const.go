package main

const (
	OK_HEADER          byte = 0x00
	MORE_DATE_HEADER   byte = 0x01
	ERR_HEADER         byte = 0xff
	EOF_HEADER         byte = 0xfe
	LocalInFile_HEADER byte = 0xfb

	CACHE_SHA2_FAST_AUTH byte = 0x03
	CACHE_SHA2_FULL_AUTH byte = 0x04
)

const (
	COM_SLEEP byte = iota
	COM_QUIT
	COM_INIT_DB
	COM_QUERY
	COM_FIELD_LIST
	COM_CREATE_DB
	COM_DROP_DB
	COM_REFRESH
	COM_SHUTDOWN
	COM_STATISTICS
	COM_PROCESS_INFO
	COM_CONNECT
	COM_PROCESS_KILL
	COM_DEBUG
	COM_PING
	COM_TIME
	COM_DELAYED_INSERT
	COM_CHANGE_USER
	COM_BINLOG_DUMP
	COM_TABLE_DUMP
	COM_CONNECT_OUT
	COM_REGISTER_SLAVE
	COM_STMT_PREPARE
	COM_STMT_EXECUTE
	COM_STMT_SEND_LONG_DATA
	COM_STMT_CLOSE
	COM_STMT_RESET
	COM_SET_OPTION
	COM_STMT_FETCH
	COM_DAEMON
	COM_BINLOG_DUMP_GTID
	COM_RESET_CONNECTION
)

const (
	DEFAULT_MYSQL_STATE = "HY000"
)

var MySQLState = map[uint16]string{
	ER_DUP_KEY:                                  "23000",
	ER_OUTOFMEMORY:                              "HY001",
	ER_OUT_OF_SORTMEMORY:                        "HY001",
	ER_CON_COUNT_ERROR:                          "08004",
	ER_BAD_HOST_ERROR:                           "08S01",
	ER_HANDSHAKE_ERROR:                          "08S01",
	ER_DBACCESS_DENIED_ERROR:                    "42000",
	ER_ACCESS_DENIED_ERROR:                      "28000",
	ER_NO_DB_ERROR:                              "3D000",
	ER_UNKNOWN_COM_ERROR:                        "08S01",
	ER_BAD_NULL_ERROR:                           "23000",
	ER_BAD_DB_ERROR:                             "42000",
	ER_TABLE_EXISTS_ERROR:                       "42S01",
	ER_BAD_TABLE_ERROR:                          "42S02",
	ER_NON_UNIQ_ERROR:                           "23000",
	ER_SERVER_SHUTDOWN:                          "08S01",
	ER_BAD_FIELD_ERROR:                          "42S22",
	ER_WRONG_FIELD_WITH_GROUP:                   "42000",
	ER_WRONG_SUM_SELECT:                         "42000",
	ER_WRONG_GROUP_FIELD:                        "42000",
	ER_WRONG_VALUE_COUNT:                        "21S01",
	ER_TOO_LONG_IDENT:                           "42000",
	ER_DUP_FIELDNAME:                            "42S21",
	ER_DUP_KEYNAME:                              "42000",
	ER_DUP_ENTRY:                                "23000",
	ER_WRONG_FIELD_SPEC:                         "42000",
	ER_PARSE_ERROR:                              "42000",
	ER_EMPTY_QUERY:                              "42000",
	ER_NONUNIQ_TABLE:                            "42000",
	ER_INVALID_DEFAULT:                          "42000",
	ER_MULTIPLE_PRI_KEY:                         "42000",
	ER_TOO_MANY_KEYS:                            "42000",
	ER_TOO_MANY_KEY_PARTS:                       "42000",
	ER_TOO_LONG_KEY:                             "42000",
	ER_KEY_COLUMN_DOES_NOT_EXITS:                "42000",
	ER_BLOB_USED_AS_KEY:                         "42000",
	ER_TOO_BIG_FIELDLENGTH:                      "42000",
	ER_WRONG_AUTO_KEY:                           "42000",
	ER_FORCING_CLOSE:                            "08S01",
	ER_IPSOCK_ERROR:                             "08S01",
	ER_NO_SUCH_INDEX:                            "42S12",
	ER_WRONG_FIELD_TERMINATORS:                  "42000",
	ER_BLOBS_AND_NO_TERMINATED:                  "42000",
	ER_CANT_REMOVE_ALL_FIELDS:                   "42000",
	ER_CANT_DROP_FIELD_OR_KEY:                   "42000",
	ER_BLOB_CANT_HAVE_DEFAULT:                   "42000",
	ER_WRONG_DB_NAME:                            "42000",
	ER_WRONG_TABLE_NAME:                         "42000",
	ER_TOO_BIG_SELECT:                           "42000",
	ER_UNKNOWN_PROCEDURE:                        "42000",
	ER_WRONG_PARAMCOUNT_TO_PROCEDURE:            "42000",
	ER_UNKNOWN_TABLE:                            "42S02",
	ER_FIELD_SPECIFIED_TWICE:                    "42000",
	ER_UNSUPPORTED_EXTENSION:                    "42000",
	ER_TABLE_MUST_HAVE_COLUMNS:                  "42000",
	ER_UNKNOWN_CHARACTER_SET:                    "42000",
	ER_TOO_BIG_ROWSIZE:                          "42000",
	ER_WRONG_OUTER_JOIN:                         "42000",
	ER_NULL_COLUMN_IN_INDEX:                     "42000",
	ER_PASSWORD_ANONYMOUS_USER:                  "42000",
	ER_PASSWORD_NOT_ALLOWED:                     "42000",
	ER_PASSWORD_NO_MATCH:                        "42000",
	ER_WRONG_VALUE_COUNT_ON_ROW:                 "21S01",
	ER_INVALID_USE_OF_NULL:                      "22004",
	ER_REGEXP_ERROR:                             "42000",
	ER_MIX_OF_GROUP_FUNC_AND_FIELDS:             "42000",
	ER_NONEXISTING_GRANT:                        "42000",
	ER_TABLEACCESS_DENIED_ERROR:                 "42000",
	ER_COLUMNACCESS_DENIED_ERROR:                "42000",
	ER_ILLEGAL_GRANT_FOR_TABLE:                  "42000",
	ER_GRANT_WRONG_HOST_OR_USER:                 "42000",
	ER_NO_SUCH_TABLE:                            "42S02",
	ER_NONEXISTING_TABLE_GRANT:                  "42000",
	ER_NOT_ALLOWED_COMMAND:                      "42000",
	ER_SYNTAX_ERROR:                             "42000",
	ER_ABORTING_CONNECTION:                      "08S01",
	ER_NET_PACKET_TOO_LARGE:                     "08S01",
	ER_NET_READ_ERROR_FROM_PIPE:                 "08S01",
	ER_NET_FCNTL_ERROR:                          "08S01",
	ER_NET_PACKETS_OUT_OF_ORDER:                 "08S01",
	ER_NET_UNCOMPRESS_ERROR:                     "08S01",
	ER_NET_READ_ERROR:                           "08S01",
	ER_NET_READ_INTERRUPTED:                     "08S01",
	ER_NET_ERROR_ON_WRITE:                       "08S01",
	ER_NET_WRITE_INTERRUPTED:                    "08S01",
	ER_TOO_LONG_STRING:                          "42000",
	ER_TABLE_CANT_HANDLE_BLOB:                   "42000",
	ER_TABLE_CANT_HANDLE_AUTO_INCREMENT:         "42000",
	ER_WRONG_COLUMN_NAME:                        "42000",
	ER_WRONG_KEY_COLUMN:                         "42000",
	ER_DUP_UNIQUE:                               "23000",
	ER_BLOB_KEY_WITHOUT_LENGTH:                  "42000",
	ER_PRIMARY_CANT_HAVE_NULL:                   "42000",
	ER_TOO_MANY_ROWS:                            "42000",
	ER_REQUIRES_PRIMARY_KEY:                     "42000",
	ER_KEY_DOES_NOT_EXITS:                       "42000",
	ER_CHECK_NO_SUCH_TABLE:                      "42000",
	ER_CHECK_NOT_IMPLEMENTED:                    "42000",
	ER_CANT_DO_THIS_DURING_AN_TRANSACTION:       "25000",
	ER_NEW_ABORTING_CONNECTION:                  "08S01",
	ER_MASTER_NET_READ:                          "08S01",
	ER_MASTER_NET_WRITE:                         "08S01",
	ER_TOO_MANY_USER_CONNECTIONS:                "42000",
	ER_READ_ONLY_TRANSACTION:                    "25000",
	ER_NO_PERMISSION_TO_CREATE_USER:             "42000",
	ER_LOCK_DEADLOCK:                            "40001",
	ER_NO_REFERENCED_ROW:                        "23000",
	ER_ROW_IS_REFERENCED:                        "23000",
	ER_CONNECT_TO_MASTER:                        "08S01",
	ER_WRONG_NUMBER_OF_COLUMNS_IN_SELECT:        "21000",
	ER_USER_LIMIT_REACHED:                       "42000",
	ER_SPECIFIC_ACCESS_DENIED_ERROR:             "42000",
	ER_NO_DEFAULT:                               "42000",
	ER_WRONG_VALUE_FOR_VAR:                      "42000",
	ER_WRONG_TYPE_FOR_VAR:                       "42000",
	ER_CANT_USE_OPTION_HERE:                     "42000",
	ER_NOT_SUPPORTED_YET:                        "42000",
	ER_WRONG_FK_DEF:                             "42000",
	ER_OPERAND_COLUMNS:                          "21000",
	ER_SUBQUERY_NO_1_ROW:                        "21000",
	ER_ILLEGAL_REFERENCE:                        "42S22",
	ER_DERIVED_MUST_HAVE_ALIAS:                  "42000",
	ER_SELECT_REDUCED:                           "01000",
	ER_TABLENAME_NOT_ALLOWED_HERE:               "42000",
	ER_NOT_SUPPORTED_AUTH_MODE:                  "08004",
	ER_SPATIAL_CANT_HAVE_NULL:                   "42000",
	ER_COLLATION_CHARSET_MISMATCH:               "42000",
	ER_WARN_TOO_FEW_RECORDS:                     "01000",
	ER_WARN_TOO_MANY_RECORDS:                    "01000",
	ER_WARN_NULL_TO_NOTNULL:                     "22004",
	ER_WARN_DATA_OUT_OF_RANGE:                   "22003",
	WARN_DATA_TRUNCATED:                         "01000",
	ER_WRONG_NAME_FOR_INDEX:                     "42000",
	ER_WRONG_NAME_FOR_CATALOG:                   "42000",
	ER_UNKNOWN_STORAGE_ENGINE:                   "42000",
	ER_TRUNCATED_WRONG_VALUE:                    "22007",
	ER_SP_NO_RECURSIVE_CREATE:                   "2F003",
	ER_SP_ALREADY_EXISTS:                        "42000",
	ER_SP_DOES_NOT_EXIST:                        "42000",
	ER_SP_LILABEL_MISMATCH:                      "42000",
	ER_SP_LABEL_REDEFINE:                        "42000",
	ER_SP_LABEL_MISMATCH:                        "42000",
	ER_SP_UNINIT_VAR:                            "01000",
	ER_SP_BADSELECT:                             "0A000",
	ER_SP_BADRETURN:                             "42000",
	ER_SP_BADSTATEMENT:                          "0A000",
	ER_UPDATE_LOG_DEPRECATED_IGNORED:            "42000",
	ER_UPDATE_LOG_DEPRECATED_TRANSLATED:         "42000",
	ER_QUERY_INTERRUPTED:                        "70100",
	ER_SP_WRONG_NO_OF_ARGS:                      "42000",
	ER_SP_COND_MISMATCH:                         "42000",
	ER_SP_NORETURN:                              "42000",
	ER_SP_NORETURNEND:                           "2F005",
	ER_SP_BAD_CURSOR_QUERY:                      "42000",
	ER_SP_BAD_CURSOR_SELECT:                     "42000",
	ER_SP_CURSOR_MISMATCH:                       "42000",
	ER_SP_CURSOR_ALREADY_OPEN:                   "24000",
	ER_SP_CURSOR_NOT_OPEN:                       "24000",
	ER_SP_UNDECLARED_VAR:                        "42000",
	ER_SP_FETCH_NO_DATA:                         "02000",
	ER_SP_DUP_PARAM:                             "42000",
	ER_SP_DUP_VAR:                               "42000",
	ER_SP_DUP_COND:                              "42000",
	ER_SP_DUP_CURS:                              "42000",
	ER_SP_SUBSELECT_NYI:                         "0A000",
	ER_STMT_NOT_ALLOWED_IN_SF_OR_TRG:            "0A000",
	ER_SP_VARCOND_AFTER_CURSHNDLR:               "42000",
	ER_SP_CURSOR_AFTER_HANDLER:                  "42000",
	ER_SP_CASE_NOT_FOUND:                        "20000",
	ER_DIVISION_BY_ZERO:                         "22012",
	ER_ILLEGAL_VALUE_FOR_TYPE:                   "22007",
	ER_PROCACCESS_DENIED_ERROR:                  "42000",
	ER_XAER_NOTA:                                "XAE04",
	ER_XAER_INVAL:                               "XAE05",
	ER_XAER_RMFAIL:                              "XAE07",
	ER_XAER_OUTSIDE:                             "XAE09",
	ER_XAER_RMERR:                               "XAE03",
	ER_XA_RBROLLBACK:                            "XA100",
	ER_NONEXISTING_PROC_GRANT:                   "42000",
	ER_DATA_TOO_LONG:                            "22001",
	ER_SP_BAD_SQLSTATE:                          "42000",
	ER_CANT_CREATE_USER_WITH_GRANT:              "42000",
	ER_SP_DUP_HANDLER:                           "42000",
	ER_SP_NOT_VAR_ARG:                           "42000",
	ER_SP_NO_RETSET:                             "0A000",
	ER_CANT_CREATE_GEOMETRY_OBJECT:              "22003",
	ER_TOO_BIG_SCALE:                            "42000",
	ER_TOO_BIG_PRECISION:                        "42000",
	ER_M_BIGGER_THAN_D:                          "42000",
	ER_TOO_LONG_BODY:                            "42000",
	ER_TOO_BIG_DISPLAYWIDTH:                     "42000",
	ER_XAER_DUPID:                               "XAE08",
	ER_DATETIME_FUNCTION_OVERFLOW:               "22008",
	ER_ROW_IS_REFERENCED_2:                      "23000",
	ER_NO_REFERENCED_ROW_2:                      "23000",
	ER_SP_BAD_VAR_SHADOW:                        "42000",
	ER_SP_WRONG_NAME:                            "42000",
	ER_SP_NO_AGGREGATE:                          "42000",
	ER_MAX_PREPARED_STMT_COUNT_REACHED:          "42000",
	ER_NON_GROUPING_FIELD_USED:                  "42000",
	ER_FOREIGN_DUPLICATE_KEY_OLD_UNUSED:         "23000",
	ER_CANT_CHANGE_TX_CHARACTERISTICS:           "25001",
	ER_WRONG_PARAMCOUNT_TO_NATIVE_FCT:           "42000",
	ER_WRONG_PARAMETERS_TO_NATIVE_FCT:           "42000",
	ER_WRONG_PARAMETERS_TO_STORED_FCT:           "42000",
	ER_DUP_ENTRY_WITH_KEY_NAME:                  "23000",
	ER_XA_RBTIMEOUT:                             "XA106",
	ER_XA_RBDEADLOCK:                            "XA102",
	ER_FUNC_INEXISTENT_NAME_COLLISION:           "42000",
	ER_DUP_SIGNAL_SET:                           "42000",
	ER_SIGNAL_WARN:                              "01000",
	ER_SIGNAL_NOT_FOUND:                         "02000",
	ER_SIGNAL_EXCEPTION:                         "HY000",
	ER_RESIGNAL_WITHOUT_ACTIVE_HANDLER:          "0K000",
	ER_SPATIAL_MUST_HAVE_GEOM_COL:               "42000",
	ER_DATA_OUT_OF_RANGE:                        "22003",
	ER_ACCESS_DENIED_NO_PASSWORD_ERROR:          "28000",
	ER_TRUNCATE_ILLEGAL_FK:                      "42000",
	ER_DA_INVALID_CONDITION_NUMBER:              "35000",
	ER_FOREIGN_DUPLICATE_KEY_WITH_CHILD_INFO:    "23000",
	ER_FOREIGN_DUPLICATE_KEY_WITHOUT_CHILD_INFO: "23000",
	ER_CANT_EXECUTE_IN_READ_ONLY_TRANSACTION:    "25006",
	ER_ALTER_OPERATION_NOT_SUPPORTED:            "0A000",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON:     "0A000",
	ER_DUP_UNKNOWN_IN_INDEX:                     "23000",
}

const (
	ER_ERROR_FIRST                                                      = 1000
	ER_HASHCHK                                                          = 1000
	ER_NISAMCHK                                                         = 1001
	ER_NO                                                               = 1002
	ER_YES                                                              = 1003
	ER_CANT_CREATE_FILE                                                 = 1004
	ER_CANT_CREATE_TABLE                                                = 1005
	ER_CANT_CREATE_DB                                                   = 1006
	ER_DB_CREATE_EXISTS                                                 = 1007
	ER_DB_DROP_EXISTS                                                   = 1008
	ER_DB_DROP_DELETE                                                   = 1009
	ER_DB_DROP_RMDIR                                                    = 1010
	ER_CANT_DELETE_FILE                                                 = 1011
	ER_CANT_FIND_SYSTEM_REC                                             = 1012
	ER_CANT_GET_STAT                                                    = 1013
	ER_CANT_GET_WD                                                      = 1014
	ER_CANT_LOCK                                                        = 1015
	ER_CANT_OPEN_FILE                                                   = 1016
	ER_FILE_NOT_FOUND                                                   = 1017
	ER_CANT_READ_DIR                                                    = 1018
	ER_CANT_SET_WD                                                      = 1019
	ER_CHECKREAD                                                        = 1020
	ER_DISK_FULL                                                        = 1021
	ER_DUP_KEY                                                          = 1022
	ER_ERROR_ON_CLOSE                                                   = 1023
	ER_ERROR_ON_READ                                                    = 1024
	ER_ERROR_ON_RENAME                                                  = 1025
	ER_ERROR_ON_WRITE                                                   = 1026
	ER_FILE_USED                                                        = 1027
	ER_FILSORT_ABORT                                                    = 1028
	ER_FORM_NOT_FOUND                                                   = 1029
	ER_GET_ERRNO                                                        = 1030
	ER_ILLEGAL_HA                                                       = 1031
	ER_KEY_NOT_FOUND                                                    = 1032
	ER_NOT_FORM_FILE                                                    = 1033
	ER_NOT_KEYFILE                                                      = 1034
	ER_OLD_KEYFILE                                                      = 1035
	ER_OPEN_AS_READONLY                                                 = 1036
	ER_OUTOFMEMORY                                                      = 1037
	ER_OUT_OF_SORTMEMORY                                                = 1038
	ER_UNEXPECTED_EOF                                                   = 1039
	ER_CON_COUNT_ERROR                                                  = 1040
	ER_OUT_OF_RESOURCES                                                 = 1041
	ER_BAD_HOST_ERROR                                                   = 1042
	ER_HANDSHAKE_ERROR                                                  = 1043
	ER_DBACCESS_DENIED_ERROR                                            = 1044
	ER_ACCESS_DENIED_ERROR                                              = 1045
	ER_NO_DB_ERROR                                                      = 1046
	ER_UNKNOWN_COM_ERROR                                                = 1047
	ER_BAD_NULL_ERROR                                                   = 1048
	ER_BAD_DB_ERROR                                                     = 1049
	ER_TABLE_EXISTS_ERROR                                               = 1050
	ER_BAD_TABLE_ERROR                                                  = 1051
	ER_NON_UNIQ_ERROR                                                   = 1052
	ER_SERVER_SHUTDOWN                                                  = 1053
	ER_BAD_FIELD_ERROR                                                  = 1054
	ER_WRONG_FIELD_WITH_GROUP                                           = 1055
	ER_WRONG_GROUP_FIELD                                                = 1056
	ER_WRONG_SUM_SELECT                                                 = 1057
	ER_WRONG_VALUE_COUNT                                                = 1058
	ER_TOO_LONG_IDENT                                                   = 1059
	ER_DUP_FIELDNAME                                                    = 1060
	ER_DUP_KEYNAME                                                      = 1061
	ER_DUP_ENTRY                                                        = 1062
	ER_WRONG_FIELD_SPEC                                                 = 1063
	ER_PARSE_ERROR                                                      = 1064
	ER_EMPTY_QUERY                                                      = 1065
	ER_NONUNIQ_TABLE                                                    = 1066
	ER_INVALID_DEFAULT                                                  = 1067
	ER_MULTIPLE_PRI_KEY                                                 = 1068
	ER_TOO_MANY_KEYS                                                    = 1069
	ER_TOO_MANY_KEY_PARTS                                               = 1070
	ER_TOO_LONG_KEY                                                     = 1071
	ER_KEY_COLUMN_DOES_NOT_EXITS                                        = 1072
	ER_BLOB_USED_AS_KEY                                                 = 1073
	ER_TOO_BIG_FIELDLENGTH                                              = 1074
	ER_WRONG_AUTO_KEY                                                   = 1075
	ER_READY                                                            = 1076
	ER_NORMAL_SHUTDOWN                                                  = 1077
	ER_GOT_SIGNAL                                                       = 1078
	ER_SHUTDOWN_COMPLETE                                                = 1079
	ER_FORCING_CLOSE                                                    = 1080
	ER_IPSOCK_ERROR                                                     = 1081
	ER_NO_SUCH_INDEX                                                    = 1082
	ER_WRONG_FIELD_TERMINATORS                                          = 1083
	ER_BLOBS_AND_NO_TERMINATED                                          = 1084
	ER_TEXTFILE_NOT_READABLE                                            = 1085
	ER_FILE_EXISTS_ERROR                                                = 1086
	ER_LOAD_INFO                                                        = 1087
	ER_ALTER_INFO                                                       = 1088
	ER_WRONG_SUB_KEY                                                    = 1089
	ER_CANT_REMOVE_ALL_FIELDS                                           = 1090
	ER_CANT_DROP_FIELD_OR_KEY                                           = 1091
	ER_INSERT_INFO                                                      = 1092
	ER_UPDATE_TABLE_USED                                                = 1093
	ER_NO_SUCH_THREAD                                                   = 1094
	ER_KILL_DENIED_ERROR                                                = 1095
	ER_NO_TABLES_USED                                                   = 1096
	ER_TOO_BIG_SET                                                      = 1097
	ER_NO_UNIQUE_LOGFILE                                                = 1098
	ER_TABLE_NOT_LOCKED_FOR_WRITE                                       = 1099
	ER_TABLE_NOT_LOCKED                                                 = 1100
	ER_BLOB_CANT_HAVE_DEFAULT                                           = 1101
	ER_WRONG_DB_NAME                                                    = 1102
	ER_WRONG_TABLE_NAME                                                 = 1103
	ER_TOO_BIG_SELECT                                                   = 1104
	ER_UNKNOWN_ERROR                                                    = 1105
	ER_UNKNOWN_PROCEDURE                                                = 1106
	ER_WRONG_PARAMCOUNT_TO_PROCEDURE                                    = 1107
	ER_WRONG_PARAMETERS_TO_PROCEDURE                                    = 1108
	ER_UNKNOWN_TABLE                                                    = 1109
	ER_FIELD_SPECIFIED_TWICE                                            = 1110
	ER_INVALID_GROUP_FUNC_USE                                           = 1111
	ER_UNSUPPORTED_EXTENSION                                            = 1112
	ER_TABLE_MUST_HAVE_COLUMNS                                          = 1113
	ER_RECORD_FILE_FULL                                                 = 1114
	ER_UNKNOWN_CHARACTER_SET                                            = 1115
	ER_TOO_MANY_TABLES                                                  = 1116
	ER_TOO_MANY_FIELDS                                                  = 1117
	ER_TOO_BIG_ROWSIZE                                                  = 1118
	ER_STACK_OVERRUN                                                    = 1119
	ER_WRONG_OUTER_JOIN                                                 = 1120
	ER_NULL_COLUMN_IN_INDEX                                             = 1121
	ER_CANT_FIND_UDF                                                    = 1122
	ER_CANT_INITIALIZE_UDF                                              = 1123
	ER_UDF_NO_PATHS                                                     = 1124
	ER_UDF_EXISTS                                                       = 1125
	ER_CANT_OPEN_LIBRARY                                                = 1126
	ER_CANT_FIND_DL_ENTRY                                               = 1127
	ER_FUNCTION_NOT_DEFINED                                             = 1128
	ER_HOST_IS_BLOCKED                                                  = 1129
	ER_HOST_NOT_PRIVILEGED                                              = 1130
	ER_PASSWORD_ANONYMOUS_USER                                          = 1131
	ER_PASSWORD_NOT_ALLOWED                                             = 1132
	ER_PASSWORD_NO_MATCH                                                = 1133
	ER_UPDATE_INFO                                                      = 1134
	ER_CANT_CREATE_THREAD                                               = 1135
	ER_WRONG_VALUE_COUNT_ON_ROW                                         = 1136
	ER_CANT_REOPEN_TABLE                                                = 1137
	ER_INVALID_USE_OF_NULL                                              = 1138
	ER_REGEXP_ERROR                                                     = 1139
	ER_MIX_OF_GROUP_FUNC_AND_FIELDS                                     = 1140
	ER_NONEXISTING_GRANT                                                = 1141
	ER_TABLEACCESS_DENIED_ERROR                                         = 1142
	ER_COLUMNACCESS_DENIED_ERROR                                        = 1143
	ER_ILLEGAL_GRANT_FOR_TABLE                                          = 1144
	ER_GRANT_WRONG_HOST_OR_USER                                         = 1145
	ER_NO_SUCH_TABLE                                                    = 1146
	ER_NONEXISTING_TABLE_GRANT                                          = 1147
	ER_NOT_ALLOWED_COMMAND                                              = 1148
	ER_SYNTAX_ERROR                                                     = 1149
	ER_DELAYED_CANT_CHANGE_LOCK                                         = 1150
	ER_TOO_MANY_DELAYED_THREADS                                         = 1151
	ER_ABORTING_CONNECTION                                              = 1152
	ER_NET_PACKET_TOO_LARGE                                             = 1153
	ER_NET_READ_ERROR_FROM_PIPE                                         = 1154
	ER_NET_FCNTL_ERROR                                                  = 1155
	ER_NET_PACKETS_OUT_OF_ORDER                                         = 1156
	ER_NET_UNCOMPRESS_ERROR                                             = 1157
	ER_NET_READ_ERROR                                                   = 1158
	ER_NET_READ_INTERRUPTED                                             = 1159
	ER_NET_ERROR_ON_WRITE                                               = 1160
	ER_NET_WRITE_INTERRUPTED                                            = 1161
	ER_TOO_LONG_STRING                                                  = 1162
	ER_TABLE_CANT_HANDLE_BLOB                                           = 1163
	ER_TABLE_CANT_HANDLE_AUTO_INCREMENT                                 = 1164
	ER_DELAYED_INSERT_TABLE_LOCKED                                      = 1165
	ER_WRONG_COLUMN_NAME                                                = 1166
	ER_WRONG_KEY_COLUMN                                                 = 1167
	ER_WRONG_MRG_TABLE                                                  = 1168
	ER_DUP_UNIQUE                                                       = 1169
	ER_BLOB_KEY_WITHOUT_LENGTH                                          = 1170
	ER_PRIMARY_CANT_HAVE_NULL                                           = 1171
	ER_TOO_MANY_ROWS                                                    = 1172
	ER_REQUIRES_PRIMARY_KEY                                             = 1173
	ER_NO_RAID_COMPILED                                                 = 1174
	ER_UPDATE_WITHOUT_KEY_IN_SAFE_MODE                                  = 1175
	ER_KEY_DOES_NOT_EXITS                                               = 1176
	ER_CHECK_NO_SUCH_TABLE                                              = 1177
	ER_CHECK_NOT_IMPLEMENTED                                            = 1178
	ER_CANT_DO_THIS_DURING_AN_TRANSACTION                               = 1179
	ER_ERROR_DURING_COMMIT                                              = 1180
	ER_ERROR_DURING_ROLLBACK                                            = 1181
	ER_ERROR_DURING_FLUSH_LOGS                                          = 1182
	ER_ERROR_DURING_CHECKPOINT                                          = 1183
	ER_NEW_ABORTING_CONNECTION                                          = 1184
	ER_DUMP_NOT_IMPLEMENTED                                             = 1185
	ER_FLUSH_MASTER_BINLOG_CLOSED                                       = 1186
	ER_INDEX_REBUILD                                                    = 1187
	ER_MASTER                                                           = 1188
	ER_MASTER_NET_READ                                                  = 1189
	ER_MASTER_NET_WRITE                                                 = 1190
	ER_FT_MATCHING_KEY_NOT_FOUND                                        = 1191
	ER_LOCK_OR_ACTIVE_TRANSACTION                                       = 1192
	ER_UNKNOWN_SYSTEM_VARIABLE                                          = 1193
	ER_CRASHED_ON_USAGE                                                 = 1194
	ER_CRASHED_ON_REPAIR                                                = 1195
	ER_WARNING_NOT_COMPLETE_ROLLBACK                                    = 1196
	ER_TRANS_CACHE_FULL                                                 = 1197
	ER_SLAVE_MUST_STOP                                                  = 1198
	ER_SLAVE_NOT_RUNNING                                                = 1199
	ER_BAD_SLAVE                                                        = 1200
	ER_MASTER_INFO                                                      = 1201
	ER_SLAVE_THREAD                                                     = 1202
	ER_TOO_MANY_USER_CONNECTIONS                                        = 1203
	ER_SET_CONSTANTS_ONLY                                               = 1204
	ER_LOCK_WAIT_TIMEOUT                                                = 1205
	ER_LOCK_TABLE_FULL                                                  = 1206
	ER_READ_ONLY_TRANSACTION                                            = 1207
	ER_DROP_DB_WITH_READ_LOCK                                           = 1208
	ER_CREATE_DB_WITH_READ_LOCK                                         = 1209
	ER_WRONG_ARGUMENTS                                                  = 1210
	ER_NO_PERMISSION_TO_CREATE_USER                                     = 1211
	ER_UNION_TABLES_IN_DIFFERENT_DIR                                    = 1212
	ER_LOCK_DEADLOCK                                                    = 1213
	ER_TABLE_CANT_HANDLE_FT                                             = 1214
	ER_CANNOT_ADD_FOREIGN                                               = 1215
	ER_NO_REFERENCED_ROW                                                = 1216
	ER_ROW_IS_REFERENCED                                                = 1217
	ER_CONNECT_TO_MASTER                                                = 1218
	ER_QUERY_ON_MASTER                                                  = 1219
	ER_ERROR_WHEN_EXECUTING_COMMAND                                     = 1220
	ER_WRONG_USAGE                                                      = 1221
	ER_WRONG_NUMBER_OF_COLUMNS_IN_SELECT                                = 1222
	ER_CANT_UPDATE_WITH_READLOCK                                        = 1223
	ER_MIXING_NOT_ALLOWED                                               = 1224
	ER_DUP_ARGUMENT                                                     = 1225
	ER_USER_LIMIT_REACHED                                               = 1226
	ER_SPECIFIC_ACCESS_DENIED_ERROR                                     = 1227
	ER_LOCAL_VARIABLE                                                   = 1228
	ER_GLOBAL_VARIABLE                                                  = 1229
	ER_NO_DEFAULT                                                       = 1230
	ER_WRONG_VALUE_FOR_VAR                                              = 1231
	ER_WRONG_TYPE_FOR_VAR                                               = 1232
	ER_VAR_CANT_BE_READ                                                 = 1233
	ER_CANT_USE_OPTION_HERE                                             = 1234
	ER_NOT_SUPPORTED_YET                                                = 1235
	ER_MASTER_FATAL_ERROR_READING_BINLOG                                = 1236
	ER_SLAVE_IGNORED_TABLE                                              = 1237
	ER_INCORRECT_GLOBAL_LOCAL_VAR                                       = 1238
	ER_WRONG_FK_DEF                                                     = 1239
	ER_KEY_REF_DO_NOT_MATCH_TABLE_REF                                   = 1240
	ER_OPERAND_COLUMNS                                                  = 1241
	ER_SUBQUERY_NO_1_ROW                                                = 1242
	ER_UNKNOWN_STMT_HANDLER                                             = 1243
	ER_CORRUPT_HELP_DB                                                  = 1244
	ER_CYCLIC_REFERENCE                                                 = 1245
	ER_AUTO_CONVERT                                                     = 1246
	ER_ILLEGAL_REFERENCE                                                = 1247
	ER_DERIVED_MUST_HAVE_ALIAS                                          = 1248
	ER_SELECT_REDUCED                                                   = 1249
	ER_TABLENAME_NOT_ALLOWED_HERE                                       = 1250
	ER_NOT_SUPPORTED_AUTH_MODE                                          = 1251
	ER_SPATIAL_CANT_HAVE_NULL                                           = 1252
	ER_COLLATION_CHARSET_MISMATCH                                       = 1253
	ER_SLAVE_WAS_RUNNING                                                = 1254
	ER_SLAVE_WAS_NOT_RUNNING                                            = 1255
	ER_TOO_BIG_FOR_UNCOMPRESS                                           = 1256
	ER_ZLIB_Z_MEM_ERROR                                                 = 1257
	ER_ZLIB_Z_BUF_ERROR                                                 = 1258
	ER_ZLIB_Z_DATA_ERROR                                                = 1259
	ER_CUT_VALUE_GROUP_CONCAT                                           = 1260
	ER_WARN_TOO_FEW_RECORDS                                             = 1261
	ER_WARN_TOO_MANY_RECORDS                                            = 1262
	ER_WARN_NULL_TO_NOTNULL                                             = 1263
	ER_WARN_DATA_OUT_OF_RANGE                                           = 1264
	WARN_DATA_TRUNCATED                                                 = 1265
	ER_WARN_USING_OTHER_HANDLER                                         = 1266
	ER_CANT_AGGREGATE_2COLLATIONS                                       = 1267
	ER_DROP_USER                                                        = 1268
	ER_REVOKE_GRANTS                                                    = 1269
	ER_CANT_AGGREGATE_3COLLATIONS                                       = 1270
	ER_CANT_AGGREGATE_NCOLLATIONS                                       = 1271
	ER_VARIABLE_IS_NOT_STRUCT                                           = 1272
	ER_UNKNOWN_COLLATION                                                = 1273
	ER_SLAVE_IGNORED_SSL_PARAMS                                         = 1274
	ER_SERVER_IS_IN_SECURE_AUTH_MODE                                    = 1275
	ER_WARN_FIELD_RESOLVED                                              = 1276
	ER_BAD_SLAVE_UNTIL_COND                                             = 1277
	ER_MISSING_SKIP_SLAVE                                               = 1278
	ER_UNTIL_COND_IGNORED                                               = 1279
	ER_WRONG_NAME_FOR_INDEX                                             = 1280
	ER_WRONG_NAME_FOR_CATALOG                                           = 1281
	ER_WARN_QC_RESIZE                                                   = 1282
	ER_BAD_FT_COLUMN                                                    = 1283
	ER_UNKNOWN_KEY_CACHE                                                = 1284
	ER_WARN_HOSTNAME_WONT_WORK                                          = 1285
	ER_UNKNOWN_STORAGE_ENGINE                                           = 1286
	ER_WARN_DEPRECATED_SYNTAX                                           = 1287
	ER_NON_UPDATABLE_TABLE                                              = 1288
	ER_FEATURE_DISABLED                                                 = 1289
	ER_OPTION_PREVENTS_STATEMENT                                        = 1290
	ER_DUPLICATED_VALUE_IN_TYPE                                         = 1291
	ER_TRUNCATED_WRONG_VALUE                                            = 1292
	ER_TOO_MUCH_AUTO_TIMESTAMP_COLS                                     = 1293
	ER_INVALID_ON_UPDATE                                                = 1294
	ER_UNSUPPORTED_PS                                                   = 1295
	ER_GET_ERRMSG                                                       = 1296
	ER_GET_TEMPORARY_ERRMSG                                             = 1297
	ER_UNKNOWN_TIME_ZONE                                                = 1298
	ER_WARN_INVALID_TIMESTAMP                                           = 1299
	ER_INVALID_CHARACTER_STRING                                         = 1300
	ER_WARN_ALLOWED_PACKET_OVERFLOWED                                   = 1301
	ER_CONFLICTING_DECLARATIONS                                         = 1302
	ER_SP_NO_RECURSIVE_CREATE                                           = 1303
	ER_SP_ALREADY_EXISTS                                                = 1304
	ER_SP_DOES_NOT_EXIST                                                = 1305
	ER_SP_DROP_FAILED                                                   = 1306
	ER_SP_STORE_FAILED                                                  = 1307
	ER_SP_LILABEL_MISMATCH                                              = 1308
	ER_SP_LABEL_REDEFINE                                                = 1309
	ER_SP_LABEL_MISMATCH                                                = 1310
	ER_SP_UNINIT_VAR                                                    = 1311
	ER_SP_BADSELECT                                                     = 1312
	ER_SP_BADRETURN                                                     = 1313
	ER_SP_BADSTATEMENT                                                  = 1314
	ER_UPDATE_LOG_DEPRECATED_IGNORED                                    = 1315
	ER_UPDATE_LOG_DEPRECATED_TRANSLATED                                 = 1316
	ER_QUERY_INTERRUPTED                                                = 1317
	ER_SP_WRONG_NO_OF_ARGS                                              = 1318
	ER_SP_COND_MISMATCH                                                 = 1319
	ER_SP_NORETURN                                                      = 1320
	ER_SP_NORETURNEND                                                   = 1321
	ER_SP_BAD_CURSOR_QUERY                                              = 1322
	ER_SP_BAD_CURSOR_SELECT                                             = 1323
	ER_SP_CURSOR_MISMATCH                                               = 1324
	ER_SP_CURSOR_ALREADY_OPEN                                           = 1325
	ER_SP_CURSOR_NOT_OPEN                                               = 1326
	ER_SP_UNDECLARED_VAR                                                = 1327
	ER_SP_WRONG_NO_OF_FETCH_ARGS                                        = 1328
	ER_SP_FETCH_NO_DATA                                                 = 1329
	ER_SP_DUP_PARAM                                                     = 1330
	ER_SP_DUP_VAR                                                       = 1331
	ER_SP_DUP_COND                                                      = 1332
	ER_SP_DUP_CURS                                                      = 1333
	ER_SP_CANT_ALTER                                                    = 1334
	ER_SP_SUBSELECT_NYI                                                 = 1335
	ER_STMT_NOT_ALLOWED_IN_SF_OR_TRG                                    = 1336
	ER_SP_VARCOND_AFTER_CURSHNDLR                                       = 1337
	ER_SP_CURSOR_AFTER_HANDLER                                          = 1338
	ER_SP_CASE_NOT_FOUND                                                = 1339
	ER_FPARSER_TOO_BIG_FILE                                             = 1340
	ER_FPARSER_BAD_HEADER                                               = 1341
	ER_FPARSER_EOF_IN_COMMENT                                           = 1342
	ER_FPARSER_ERROR_IN_PARAMETER                                       = 1343
	ER_FPARSER_EOF_IN_UNKNOWN_PARAMETER                                 = 1344
	ER_VIEW_NO_EXPLAIN                                                  = 1345
	ER_FRM_UNKNOWN_TYPE                                                 = 1346
	ER_WRONG_OBJECT                                                     = 1347
	ER_NONUPDATEABLE_COLUMN                                             = 1348
	ER_VIEW_SELECT_DERIVED                                              = 1349
	ER_VIEW_SELECT_CLAUSE                                               = 1350
	ER_VIEW_SELECT_VARIABLE                                             = 1351
	ER_VIEW_SELECT_TMPTABLE                                             = 1352
	ER_VIEW_WRONG_LIST                                                  = 1353
	ER_WARN_VIEW_MERGE                                                  = 1354
	ER_WARN_VIEW_WITHOUT_KEY                                            = 1355
	ER_VIEW_INVALID                                                     = 1356
	ER_SP_NO_DROP_SP                                                    = 1357
	ER_SP_GOTO_IN_HNDLR                                                 = 1358
	ER_TRG_ALREADY_EXISTS                                               = 1359
	ER_TRG_DOES_NOT_EXIST                                               = 1360
	ER_TRG_ON_VIEW_OR_TEMP_TABLE                                        = 1361
	ER_TRG_CANT_CHANGE_ROW                                              = 1362
	ER_TRG_NO_SUCH_ROW_IN_TRG                                           = 1363
	ER_NO_DEFAULT_FOR_FIELD                                             = 1364
	ER_DIVISION_BY_ZERO                                                 = 1365
	ER_TRUNCATED_WRONG_VALUE_FOR_FIELD                                  = 1366
	ER_ILLEGAL_VALUE_FOR_TYPE                                           = 1367
	ER_VIEW_NONUPD_CHECK                                                = 1368
	ER_VIEW_CHECK_FAILED                                                = 1369
	ER_PROCACCESS_DENIED_ERROR                                          = 1370
	ER_RELAY_LOG_FAIL                                                   = 1371
	ER_PASSWD_LENGTH                                                    = 1372
	ER_UNKNOWN_TARGET_BINLOG                                            = 1373
	ER_IO_ERR_LOG_INDEX_READ                                            = 1374
	ER_BINLOG_PURGE_PROHIBITED                                          = 1375
	ER_FSEEK_FAIL                                                       = 1376
	ER_BINLOG_PURGE_FATAL_ERR                                           = 1377
	ER_LOG_IN_USE                                                       = 1378
	ER_LOG_PURGE_UNKNOWN_ERR                                            = 1379
	ER_RELAY_LOG_INIT                                                   = 1380
	ER_NO_BINARY_LOGGING                                                = 1381
	ER_RESERVED_SYNTAX                                                  = 1382
	ER_WSAS_FAILED                                                      = 1383
	ER_DIFF_GROUPS_PROC                                                 = 1384
	ER_NO_GROUP_FOR_PROC                                                = 1385
	ER_ORDER_WITH_PROC                                                  = 1386
	ER_LOGGING_PROHIBIT_CHANGING_OF                                     = 1387
	ER_NO_FILE_MAPPING                                                  = 1388
	ER_WRONG_MAGIC                                                      = 1389
	ER_PS_MANY_PARAM                                                    = 1390
	ER_KEY_PART_0                                                       = 1391
	ER_VIEW_CHECKSUM                                                    = 1392
	ER_VIEW_MULTIUPDATE                                                 = 1393
	ER_VIEW_NO_INSERT_FIELD_LIST                                        = 1394
	ER_VIEW_DELETE_MERGE_VIEW                                           = 1395
	ER_CANNOT_USER                                                      = 1396
	ER_XAER_NOTA                                                        = 1397
	ER_XAER_INVAL                                                       = 1398
	ER_XAER_RMFAIL                                                      = 1399
	ER_XAER_OUTSIDE                                                     = 1400
	ER_XAER_RMERR                                                       = 1401
	ER_XA_RBROLLBACK                                                    = 1402
	ER_NONEXISTING_PROC_GRANT                                           = 1403
	ER_PROC_AUTO_GRANT_FAIL                                             = 1404
	ER_PROC_AUTO_REVOKE_FAIL                                            = 1405
	ER_DATA_TOO_LONG                                                    = 1406
	ER_SP_BAD_SQLSTATE                                                  = 1407
	ER_STARTUP                                                          = 1408
	ER_LOAD_FROM_FIXED_SIZE_ROWS_TO_VAR                                 = 1409
	ER_CANT_CREATE_USER_WITH_GRANT                                      = 1410
	ER_WRONG_VALUE_FOR_TYPE                                             = 1411
	ER_TABLE_DEF_CHANGED                                                = 1412
	ER_SP_DUP_HANDLER                                                   = 1413
	ER_SP_NOT_VAR_ARG                                                   = 1414
	ER_SP_NO_RETSET                                                     = 1415
	ER_CANT_CREATE_GEOMETRY_OBJECT                                      = 1416
	ER_FAILED_ROUTINE_BREAK_BINLOG                                      = 1417
	ER_BINLOG_UNSAFE_ROUTINE                                            = 1418
	ER_BINLOG_CREATE_ROUTINE_NEED_SUPER                                 = 1419
	ER_EXEC_STMT_WITH_OPEN_CURSOR                                       = 1420
	ER_STMT_HAS_NO_OPEN_CURSOR                                          = 1421
	ER_COMMIT_NOT_ALLOWED_IN_SF_OR_TRG                                  = 1422
	ER_NO_DEFAULT_FOR_VIEW_FIELD                                        = 1423
	ER_SP_NO_RECURSION                                                  = 1424
	ER_TOO_BIG_SCALE                                                    = 1425
	ER_TOO_BIG_PRECISION                                                = 1426
	ER_M_BIGGER_THAN_D                                                  = 1427
	ER_WRONG_LOCK_OF_SYSTEM_TABLE                                       = 1428
	ER_CONNECT_TO_FOREIGN_DATA_SOURCE                                   = 1429
	ER_QUERY_ON_FOREIGN_DATA_SOURCE                                     = 1430
	ER_FOREIGN_DATA_SOURCE_DOESNT_EXIST                                 = 1431
	ER_FOREIGN_DATA_STRING_INVALID_CANT_CREATE                          = 1432
	ER_FOREIGN_DATA_STRING_INVALID                                      = 1433
	ER_CANT_CREATE_FEDERATED_TABLE                                      = 1434
	ER_TRG_IN_WRONG_SCHEMA                                              = 1435
	ER_STACK_OVERRUN_NEED_MORE                                          = 1436
	ER_TOO_LONG_BODY                                                    = 1437
	ER_WARN_CANT_DROP_DEFAULT_KEYCACHE                                  = 1438
	ER_TOO_BIG_DISPLAYWIDTH                                             = 1439
	ER_XAER_DUPID                                                       = 1440
	ER_DATETIME_FUNCTION_OVERFLOW                                       = 1441
	ER_CANT_UPDATE_USED_TABLE_IN_SF_OR_TRG                              = 1442
	ER_VIEW_PREVENT_UPDATE                                              = 1443
	ER_PS_NO_RECURSION                                                  = 1444
	ER_SP_CANT_SET_AUTOCOMMIT                                           = 1445
	ER_MALFORMED_DEFINER                                                = 1446
	ER_VIEW_FRM_NO_USER                                                 = 1447
	ER_VIEW_OTHER_USER                                                  = 1448
	ER_NO_SUCH_USER                                                     = 1449
	ER_FORBID_SCHEMA_CHANGE                                             = 1450
	ER_ROW_IS_REFERENCED_2                                              = 1451
	ER_NO_REFERENCED_ROW_2                                              = 1452
	ER_SP_BAD_VAR_SHADOW                                                = 1453
	ER_TRG_NO_DEFINER                                                   = 1454
	ER_OLD_FILE_FORMAT                                                  = 1455
	ER_SP_RECURSION_LIMIT                                               = 1456
	ER_SP_PROC_TABLE_CORRUPT                                            = 1457
	ER_SP_WRONG_NAME                                                    = 1458
	ER_TABLE_NEEDS_UPGRADE                                              = 1459
	ER_SP_NO_AGGREGATE                                                  = 1460
	ER_MAX_PREPARED_STMT_COUNT_REACHED                                  = 1461
	ER_VIEW_RECURSIVE                                                   = 1462
	ER_NON_GROUPING_FIELD_USED                                          = 1463
	ER_TABLE_CANT_HANDLE_SPKEYS                                         = 1464
	ER_NO_TRIGGERS_ON_SYSTEM_SCHEMA                                     = 1465
	ER_REMOVED_SPACES                                                   = 1466
	ER_AUTOINC_READ_FAILED                                              = 1467
	ER_USERNAME                                                         = 1468
	ER_HOSTNAME                                                         = 1469
	ER_WRONG_STRING_LENGTH                                              = 1470
	ER_NON_INSERTABLE_TABLE                                             = 1471
	ER_ADMIN_WRONG_MRG_TABLE                                            = 1472
	ER_TOO_HIGH_LEVEL_OF_NESTING_FOR_SELECT                             = 1473
	ER_NAME_BECOMES_EMPTY                                               = 1474
	ER_AMBIGUOUS_FIELD_TERM                                             = 1475
	ER_FOREIGN_SERVER_EXISTS                                            = 1476
	ER_FOREIGN_SERVER_DOESNT_EXIST                                      = 1477
	ER_ILLEGAL_HA_CREATE_OPTION                                         = 1478
	ER_PARTITION_REQUIRES_VALUES_ERROR                                  = 1479
	ER_PARTITION_WRONG_VALUES_ERROR                                     = 1480
	ER_PARTITION_MAXVALUE_ERROR                                         = 1481
	ER_PARTITION_SUBPARTITION_ERROR                                     = 1482
	ER_PARTITION_SUBPART_MIX_ERROR                                      = 1483
	ER_PARTITION_WRONG_NO_PART_ERROR                                    = 1484
	ER_PARTITION_WRONG_NO_SUBPART_ERROR                                 = 1485
	ER_WRONG_EXPR_IN_PARTITION_FUNC_ERROR                               = 1486
	ER_NO_CONST_EXPR_IN_RANGE_OR_LIST_ERROR                             = 1487
	ER_FIELD_NOT_FOUND_PART_ERROR                                       = 1488
	ER_LIST_OF_FIELDS_ONLY_IN_HASH_ERROR                                = 1489
	ER_INCONSISTENT_PARTITION_INFO_ERROR                                = 1490
	ER_PARTITION_FUNC_NOT_ALLOWED_ERROR                                 = 1491
	ER_PARTITIONS_MUST_BE_DEFINED_ERROR                                 = 1492
	ER_RANGE_NOT_INCREASING_ERROR                                       = 1493
	ER_INCONSISTENT_TYPE_OF_FUNCTIONS_ERROR                             = 1494
	ER_MULTIPLE_DEF_CONST_IN_LIST_PART_ERROR                            = 1495
	ER_PARTITION_ENTRY_ERROR                                            = 1496
	ER_MIX_HANDLER_ERROR                                                = 1497
	ER_PARTITION_NOT_DEFINED_ERROR                                      = 1498
	ER_TOO_MANY_PARTITIONS_ERROR                                        = 1499
	ER_SUBPARTITION_ERROR                                               = 1500
	ER_CANT_CREATE_HANDLER_FILE                                         = 1501
	ER_BLOB_FIELD_IN_PART_FUNC_ERROR                                    = 1502
	ER_UNIQUE_KEY_NEED_ALL_FIELDS_IN_PF                                 = 1503
	ER_NO_PARTS_ERROR                                                   = 1504
	ER_PARTITION_MGMT_ON_NONPARTITIONED                                 = 1505
	ER_FOREIGN_KEY_ON_PARTITIONED                                       = 1506
	ER_DROP_PARTITION_NON_EXISTENT                                      = 1507
	ER_DROP_LAST_PARTITION                                              = 1508
	ER_COALESCE_ONLY_ON_HASH_PARTITION                                  = 1509
	ER_REORG_HASH_ONLY_ON_SAME_NO                                       = 1510
	ER_REORG_NO_PARAM_ERROR                                             = 1511
	ER_ONLY_ON_RANGE_LIST_PARTITION                                     = 1512
	ER_ADD_PARTITION_SUBPART_ERROR                                      = 1513
	ER_ADD_PARTITION_NO_NEW_PARTITION                                   = 1514
	ER_COALESCE_PARTITION_NO_PARTITION                                  = 1515
	ER_REORG_PARTITION_NOT_EXIST                                        = 1516
	ER_SAME_NAME_PARTITION                                              = 1517
	ER_NO_BINLOG_ERROR                                                  = 1518
	ER_CONSECUTIVE_REORG_PARTITIONS                                     = 1519
	ER_REORG_OUTSIDE_RANGE                                              = 1520
	ER_PARTITION_FUNCTION_FAILURE                                       = 1521
	ER_PART_STATE_ERROR                                                 = 1522
	ER_LIMITED_PART_RANGE                                               = 1523
	ER_PLUGIN_IS_NOT_LOADED                                             = 1524
	ER_WRONG_VALUE                                                      = 1525
	ER_NO_PARTITION_FOR_GIVEN_VALUE                                     = 1526
	ER_FILEGROUP_OPTION_ONLY_ONCE                                       = 1527
	ER_CREATE_FILEGROUP_FAILED                                          = 1528
	ER_DROP_FILEGROUP_FAILED                                            = 1529
	ER_TABLESPACE_AUTO_EXTEND_ERROR                                     = 1530
	ER_WRONG_SIZE_NUMBER                                                = 1531
	ER_SIZE_OVERFLOW_ERROR                                              = 1532
	ER_ALTER_FILEGROUP_FAILED                                           = 1533
	ER_BINLOG_ROW_LOGGING_FAILED                                        = 1534
	ER_BINLOG_ROW_WRONG_TABLE_DEF                                       = 1535
	ER_BINLOG_ROW_RBR_TO_SBR                                            = 1536
	ER_EVENT_ALREADY_EXISTS                                             = 1537
	ER_EVENT_STORE_FAILED                                               = 1538
	ER_EVENT_DOES_NOT_EXIST                                             = 1539
	ER_EVENT_CANT_ALTER                                                 = 1540
	ER_EVENT_DROP_FAILED                                                = 1541
	ER_EVENT_INTERVAL_NOT_POSITIVE_OR_TOO_BIG                           = 1542
	ER_EVENT_ENDS_BEFORE_STARTS                                         = 1543
	ER_EVENT_EXEC_TIME_IN_THE_PAST                                      = 1544
	ER_EVENT_OPEN_TABLE_FAILED                                          = 1545
	ER_EVENT_NEITHER_M_EXPR_NOR_M_AT                                    = 1546
	ER_OBSOLETE_COL_COUNT_DOESNT_MATCH_CORRUPTED                        = 1547
	ER_OBSOLETE_CANNOT_LOAD_FROM_TABLE                                  = 1548
	ER_EVENT_CANNOT_DELETE                                              = 1549
	ER_EVENT_COMPILE_ERROR                                              = 1550
	ER_EVENT_SAME_NAME                                                  = 1551
	ER_EVENT_DATA_TOO_LONG                                              = 1552
	ER_DROP_INDEX_FK                                                    = 1553
	ER_WARN_DEPRECATED_SYNTAX_WITH_VER                                  = 1554
	ER_CANT_WRITE_LOCK_LOG_TABLE                                        = 1555
	ER_CANT_LOCK_LOG_TABLE                                              = 1556
	ER_FOREIGN_DUPLICATE_KEY_OLD_UNUSED                                 = 1557
	ER_COL_COUNT_DOESNT_MATCH_PLEASE_UPDATE                             = 1558
	ER_TEMP_TABLE_PREVENTS_SWITCH_OUT_OF_RBR                            = 1559
	ER_STORED_FUNCTION_PREVENTS_SWITCH_BINLOG_FORMAT                    = 1560
	ER_NDB_CANT_SWITCH_BINLOG_FORMAT                                    = 1561
	ER_PARTITION_NO_TEMPORARY                                           = 1562
	ER_PARTITION_CONST_DOMAIN_ERROR                                     = 1563
	ER_PARTITION_FUNCTION_IS_NOT_ALLOWED                                = 1564
	ER_DDL_LOG_ERROR                                                    = 1565
	ER_NULL_IN_VALUES_LESS_THAN                                         = 1566
	ER_WRONG_PARTITION_NAME                                             = 1567
	ER_CANT_CHANGE_TX_CHARACTERISTICS                                   = 1568
	ER_DUP_ENTRY_AUTOINCREMENT_CASE                                     = 1569
	ER_EVENT_MODIFY_QUEUE_ERROR                                         = 1570
	ER_EVENT_SET_VAR_ERROR                                              = 1571
	ER_PARTITION_MERGE_ERROR                                            = 1572
	ER_CANT_ACTIVATE_LOG                                                = 1573
	ER_RBR_NOT_AVAILABLE                                                = 1574
	ER_BASE64_DECODE_ERROR                                              = 1575
	ER_EVENT_RECURSION_FORBIDDEN                                        = 1576
	ER_EVENTS_DB_ERROR                                                  = 1577
	ER_ONLY_INTEGERS_ALLOWED                                            = 1578
	ER_UNSUPORTED_LOG_ENGINE                                            = 1579
	ER_BAD_LOG_STATEMENT                                                = 1580
	ER_CANT_RENAME_LOG_TABLE                                            = 1581
	ER_WRONG_PARAMCOUNT_TO_NATIVE_FCT                                   = 1582
	ER_WRONG_PARAMETERS_TO_NATIVE_FCT                                   = 1583
	ER_WRONG_PARAMETERS_TO_STORED_FCT                                   = 1584
	ER_NATIVE_FCT_NAME_COLLISION                                        = 1585
	ER_DUP_ENTRY_WITH_KEY_NAME                                          = 1586
	ER_BINLOG_PURGE_EMFILE                                              = 1587
	ER_EVENT_CANNOT_CREATE_IN_THE_PAST                                  = 1588
	ER_EVENT_CANNOT_ALTER_IN_THE_PAST                                   = 1589
	ER_SLAVE_INCIDENT                                                   = 1590
	ER_NO_PARTITION_FOR_GIVEN_VALUE_SILENT                              = 1591
	ER_BINLOG_UNSAFE_STATEMENT                                          = 1592
	ER_SLAVE_FATAL_ERROR                                                = 1593
	ER_SLAVE_RELAY_LOG_READ_FAILURE                                     = 1594
	ER_SLAVE_RELAY_LOG_WRITE_FAILURE                                    = 1595
	ER_SLAVE_CREATE_EVENT_FAILURE                                       = 1596
	ER_SLAVE_MASTER_COM_FAILURE                                         = 1597
	ER_BINLOG_LOGGING_IMPOSSIBLE                                        = 1598
	ER_VIEW_NO_CREATION_CTX                                             = 1599
	ER_VIEW_INVALID_CREATION_CTX                                        = 1600
	ER_SR_INVALID_CREATION_CTX                                          = 1601
	ER_TRG_CORRUPTED_FILE                                               = 1602
	ER_TRG_NO_CREATION_CTX                                              = 1603
	ER_TRG_INVALID_CREATION_CTX                                         = 1604
	ER_EVENT_INVALID_CREATION_CTX                                       = 1605
	ER_TRG_CANT_OPEN_TABLE                                              = 1606
	ER_CANT_CREATE_SROUTINE                                             = 1607
	ER_NEVER_USED                                                       = 1608
	ER_NO_FORMAT_DESCRIPTION_EVENT_BEFORE_BINLOG_STATEMENT              = 1609
	ER_SLAVE_CORRUPT_EVENT                                              = 1610
	ER_LOAD_DATA_INVALID_COLUMN                                         = 1611
	ER_LOG_PURGE_NO_FILE                                                = 1612
	ER_XA_RBTIMEOUT                                                     = 1613
	ER_XA_RBDEADLOCK                                                    = 1614
	ER_NEED_REPREPARE                                                   = 1615
	ER_DELAYED_NOT_SUPPORTED                                            = 1616
	WARN_NO_MASTER_INFO                                                 = 1617
	WARN_OPTION_IGNORED                                                 = 1618
	WARN_PLUGIN_DELETE_BUILTIN                                          = 1619
	WARN_PLUGIN_BUSY                                                    = 1620
	ER_VARIABLE_IS_READONLY                                             = 1621
	ER_WARN_ENGINE_TRANSACTION_ROLLBACK                                 = 1622
	ER_SLAVE_HEARTBEAT_FAILURE                                          = 1623
	ER_SLAVE_HEARTBEAT_VALUE_OUT_OF_RANGE                               = 1624
	ER_NDB_REPLICATION_SCHEMA_ERROR                                     = 1625
	ER_CONFLICT_FN_PARSE_ERROR                                          = 1626
	ER_EXCEPTIONS_WRITE_ERROR                                           = 1627
	ER_TOO_LONG_TABLE_COMMENT                                           = 1628
	ER_TOO_LONG_FIELD_COMMENT                                           = 1629
	ER_FUNC_INEXISTENT_NAME_COLLISION                                   = 1630
	ER_DATABASE_NAME                                                    = 1631
	ER_TABLE_NAME                                                       = 1632
	ER_PARTITION_NAME                                                   = 1633
	ER_SUBPARTITION_NAME                                                = 1634
	ER_TEMPORARY_NAME                                                   = 1635
	ER_RENAMED_NAME                                                     = 1636
	ER_TOO_MANY_CONCURRENT_TRXS                                         = 1637
	WARN_NON_ASCII_SEPARATOR_NOT_IMPLEMENTED                            = 1638
	ER_DEBUG_SYNC_TIMEOUT                                               = 1639
	ER_DEBUG_SYNC_HIT_LIMIT                                             = 1640
	ER_DUP_SIGNAL_SET                                                   = 1641
	ER_SIGNAL_WARN                                                      = 1642
	ER_SIGNAL_NOT_FOUND                                                 = 1643
	ER_SIGNAL_EXCEPTION                                                 = 1644
	ER_RESIGNAL_WITHOUT_ACTIVE_HANDLER                                  = 1645
	ER_SIGNAL_BAD_CONDITION_TYPE                                        = 1646
	WARN_COND_ITEM_TRUNCATED                                            = 1647
	ER_COND_ITEM_TOO_LONG                                               = 1648
	ER_UNKNOWN_LOCALE                                                   = 1649
	ER_SLAVE_IGNORE_SERVER_IDS                                          = 1650
	ER_QUERY_CACHE_DISABLED                                             = 1651
	ER_SAME_NAME_PARTITION_FIELD                                        = 1652
	ER_PARTITION_COLUMN_LIST_ERROR                                      = 1653
	ER_WRONG_TYPE_COLUMN_VALUE_ERROR                                    = 1654
	ER_TOO_MANY_PARTITION_FUNC_FIELDS_ERROR                             = 1655
	ER_MAXVALUE_IN_VALUES_IN                                            = 1656
	ER_TOO_MANY_VALUES_ERROR                                            = 1657
	ER_ROW_SINGLE_PARTITION_FIELD_ERROR                                 = 1658
	ER_FIELD_TYPE_NOT_ALLOWED_AS_PARTITION_FIELD                        = 1659
	ER_PARTITION_FIELDS_TOO_LONG                                        = 1660
	ER_BINLOG_ROW_ENGINE_AND_STMT_ENGINE                                = 1661
	ER_BINLOG_ROW_MODE_AND_STMT_ENGINE                                  = 1662
	ER_BINLOG_UNSAFE_AND_STMT_ENGINE                                    = 1663
	ER_BINLOG_ROW_INJECTION_AND_STMT_ENGINE                             = 1664
	ER_BINLOG_STMT_MODE_AND_ROW_ENGINE                                  = 1665
	ER_BINLOG_ROW_INJECTION_AND_STMT_MODE                               = 1666
	ER_BINLOG_MULTIPLE_ENGINES_AND_SELF_LOGGING_ENGINE                  = 1667
	ER_BINLOG_UNSAFE_LIMIT                                              = 1668
	ER_BINLOG_UNSAFE_INSERT_DELAYED                                     = 1669
	ER_BINLOG_UNSAFE_SYSTEM_TABLE                                       = 1670
	ER_BINLOG_UNSAFE_AUTOINC_COLUMNS                                    = 1671
	ER_BINLOG_UNSAFE_UDF                                                = 1672
	ER_BINLOG_UNSAFE_SYSTEM_VARIABLE                                    = 1673
	ER_BINLOG_UNSAFE_SYSTEM_FUNCTION                                    = 1674
	ER_BINLOG_UNSAFE_NONTRANS_AFTER_TRANS                               = 1675
	ER_MESSAGE_AND_STATEMENT                                            = 1676
	ER_SLAVE_CONVERSION_FAILED                                          = 1677
	ER_SLAVE_CANT_CREATE_CONVERSION                                     = 1678
	ER_INSIDE_TRANSACTION_PREVENTS_SWITCH_BINLOG_FORMAT                 = 1679
	ER_PATH_LENGTH                                                      = 1680
	ER_WARN_DEPRECATED_SYNTAX_NO_REPLACEMENT                            = 1681
	ER_WRONG_NATIVE_TABLE_STRUCTURE                                     = 1682
	ER_WRONG_PERFSCHEMA_USAGE                                           = 1683
	ER_WARN_I_S_SKIPPED_TABLE                                           = 1684
	ER_INSIDE_TRANSACTION_PREVENTS_SWITCH_BINLOG_DIRECT                 = 1685
	ER_STORED_FUNCTION_PREVENTS_SWITCH_BINLOG_DIRECT                    = 1686
	ER_SPATIAL_MUST_HAVE_GEOM_COL                                       = 1687
	ER_TOO_LONG_INDEX_COMMENT                                           = 1688
	ER_LOCK_ABORTED                                                     = 1689
	ER_DATA_OUT_OF_RANGE                                                = 1690
	ER_WRONG_SPVAR_TYPE_IN_LIMIT                                        = 1691
	ER_BINLOG_UNSAFE_MULTIPLE_ENGINES_AND_SELF_LOGGING_ENGINE           = 1692
	ER_BINLOG_UNSAFE_MIXED_STATEMENT                                    = 1693
	ER_INSIDE_TRANSACTION_PREVENTS_SWITCH_SQL_LOG_BIN                   = 1694
	ER_STORED_FUNCTION_PREVENTS_SWITCH_SQL_LOG_BIN                      = 1695
	ER_FAILED_READ_FROM_PAR_FILE                                        = 1696
	ER_VALUES_IS_NOT_INT_TYPE_ERROR                                     = 1697
	ER_ACCESS_DENIED_NO_PASSWORD_ERROR                                  = 1698
	ER_SET_PASSWORD_AUTH_PLUGIN                                         = 1699
	ER_GRANT_PLUGIN_USER_EXISTS                                         = 1700
	ER_TRUNCATE_ILLEGAL_FK                                              = 1701
	ER_PLUGIN_IS_PERMANENT                                              = 1702
	ER_SLAVE_HEARTBEAT_VALUE_OUT_OF_RANGE_MIN                           = 1703
	ER_SLAVE_HEARTBEAT_VALUE_OUT_OF_RANGE_MAX                           = 1704
	ER_STMT_CACHE_FULL                                                  = 1705
	ER_MULTI_UPDATE_KEY_CONFLICT                                        = 1706
	ER_TABLE_NEEDS_REBUILD                                              = 1707
	WARN_OPTION_BELOW_LIMIT                                             = 1708
	ER_INDEX_COLUMN_TOO_LONG                                            = 1709
	ER_ERROR_IN_TRIGGER_BODY                                            = 1710
	ER_ERROR_IN_UNKNOWN_TRIGGER_BODY                                    = 1711
	ER_INDEX_CORRUPT                                                    = 1712
	ER_UNDO_RECORD_TOO_BIG                                              = 1713
	ER_BINLOG_UNSAFE_INSERT_IGNORE_SELECT                               = 1714
	ER_BINLOG_UNSAFE_INSERT_SELECT_UPDATE                               = 1715
	ER_BINLOG_UNSAFE_REPLACE_SELECT                                     = 1716
	ER_BINLOG_UNSAFE_CREATE_IGNORE_SELECT                               = 1717
	ER_BINLOG_UNSAFE_CREATE_REPLACE_SELECT                              = 1718
	ER_BINLOG_UNSAFE_UPDATE_IGNORE                                      = 1719
	ER_PLUGIN_NO_UNINSTALL                                              = 1720
	ER_PLUGIN_NO_INSTALL                                                = 1721
	ER_BINLOG_UNSAFE_WRITE_AUTOINC_SELECT                               = 1722
	ER_BINLOG_UNSAFE_CREATE_SELECT_AUTOINC                              = 1723
	ER_BINLOG_UNSAFE_INSERT_TWO_KEYS                                    = 1724
	ER_TABLE_IN_FK_CHECK                                                = 1725
	ER_UNSUPPORTED_ENGINE                                               = 1726
	ER_BINLOG_UNSAFE_AUTOINC_NOT_FIRST                                  = 1727
	ER_CANNOT_LOAD_FROM_TABLE_V2                                        = 1728
	ER_MASTER_DELAY_VALUE_OUT_OF_RANGE                                  = 1729
	ER_ONLY_FD_AND_RBR_EVENTS_ALLOWED_IN_BINLOG_STATEMENT               = 1730
	ER_PARTITION_EXCHANGE_DIFFERENT_OPTION                              = 1731
	ER_PARTITION_EXCHANGE_PART_TABLE                                    = 1732
	ER_PARTITION_EXCHANGE_TEMP_TABLE                                    = 1733
	ER_PARTITION_INSTEAD_OF_SUBPARTITION                                = 1734
	ER_UNKNOWN_PARTITION                                                = 1735
	ER_TABLES_DIFFERENT_METADATA                                        = 1736
	ER_ROW_DOES_NOT_MATCH_PARTITION                                     = 1737
	ER_BINLOG_CACHE_SIZE_GREATER_THAN_MAX                               = 1738
	ER_WARN_INDEX_NOT_APPLICABLE                                        = 1739
	ER_PARTITION_EXCHANGE_FOREIGN_KEY                                   = 1740
	ER_NO_SUCH_KEY_VALUE                                                = 1741
	ER_RPL_INFO_DATA_TOO_LONG                                           = 1742
	ER_NETWORK_READ_EVENT_CHECKSUM_FAILURE                              = 1743
	ER_BINLOG_READ_EVENT_CHECKSUM_FAILURE                               = 1744
	ER_BINLOG_STMT_CACHE_SIZE_GREATER_THAN_MAX                          = 1745
	ER_CANT_UPDATE_TABLE_IN_CREATE_TABLE_SELECT                         = 1746
	ER_PARTITION_CLAUSE_ON_NONPARTITIONED                               = 1747
	ER_ROW_DOES_NOT_MATCH_GIVEN_PARTITION_SET                           = 1748
	ER_NO_SUCH_PARTITION__UNUSED                                        = 1749
	ER_CHANGE_RPL_INFO_REPOSITORY_FAILURE                               = 1750
	ER_WARNING_NOT_COMPLETE_ROLLBACK_WITH_CREATED_TEMP_TABLE            = 1751
	ER_WARNING_NOT_COMPLETE_ROLLBACK_WITH_DROPPED_TEMP_TABLE            = 1752
	ER_MTS_FEATURE_IS_NOT_SUPPORTED                                     = 1753
	ER_MTS_UPDATED_DBS_GREATER_MAX                                      = 1754
	ER_MTS_CANT_PARALLEL                                                = 1755
	ER_MTS_INCONSISTENT_DATA                                            = 1756
	ER_FULLTEXT_NOT_SUPPORTED_WITH_PARTITIONING                         = 1757
	ER_DA_INVALID_CONDITION_NUMBER                                      = 1758
	ER_INSECURE_PLAIN_TEXT                                              = 1759
	ER_INSECURE_CHANGE_MASTER                                           = 1760
	ER_FOREIGN_DUPLICATE_KEY_WITH_CHILD_INFO                            = 1761
	ER_FOREIGN_DUPLICATE_KEY_WITHOUT_CHILD_INFO                         = 1762
	ER_SQLTHREAD_WITH_SECURE_SLAVE                                      = 1763
	ER_TABLE_HAS_NO_FT                                                  = 1764
	ER_VARIABLE_NOT_SETTABLE_IN_SF_OR_TRIGGER                           = 1765
	ER_VARIABLE_NOT_SETTABLE_IN_TRANSACTION                             = 1766
	ER_GTID_NEXT_IS_NOT_IN_GTID_NEXT_LIST                               = 1767
	ER_CANT_CHANGE_GTID_NEXT_IN_TRANSACTION_WHEN_GTID_NEXT_LIST_IS_NULL = 1768
	ER_SET_STATEMENT_CANNOT_INVOKE_FUNCTION                             = 1769
	ER_GTID_NEXT_CANT_BE_AUTOMATIC_IF_GTID_NEXT_LIST_IS_NON_NULL        = 1770
	ER_SKIPPING_LOGGED_TRANSACTION                                      = 1771
	ER_MALFORMED_GTID_SET_SPECIFICATION                                 = 1772
	ER_MALFORMED_GTID_SET_ENCODING                                      = 1773
	ER_MALFORMED_GTID_SPECIFICATION                                     = 1774
	ER_GNO_EXHAUSTED                                                    = 1775
	ER_BAD_SLAVE_AUTO_POSITION                                          = 1776
	ER_AUTO_POSITION_REQUIRES_GTID_MODE_ON                              = 1777
	ER_CANT_DO_IMPLICIT_COMMIT_IN_TRX_WHEN_GTID_NEXT_IS_SET             = 1778
	ER_GTID_MODE_2_OR_3_REQUIRES_ENFORCE_GTID_CONSISTENCY_ON            = 1779
	ER_GTID_MODE_REQUIRES_BINLOG                                        = 1780
	ER_CANT_SET_GTID_NEXT_TO_GTID_WHEN_GTID_MODE_IS_OFF                 = 1781
	ER_CANT_SET_GTID_NEXT_TO_ANONYMOUS_WHEN_GTID_MODE_IS_ON             = 1782
	ER_CANT_SET_GTID_NEXT_LIST_TO_NON_NULL_WHEN_GTID_MODE_IS_OFF        = 1783
	ER_FOUND_GTID_EVENT_WHEN_GTID_MODE_IS_OFF                           = 1784
	ER_GTID_UNSAFE_NON_TRANSACTIONAL_TABLE                              = 1785
	ER_GTID_UNSAFE_CREATE_SELECT                                        = 1786
	ER_GTID_UNSAFE_CREATE_DROP_TEMPORARY_TABLE_IN_TRANSACTION           = 1787
	ER_GTID_MODE_CAN_ONLY_CHANGE_ONE_STEP_AT_A_TIME                     = 1788
	ER_MASTER_HAS_PURGED_REQUIRED_GTIDS                                 = 1789
	ER_CANT_SET_GTID_NEXT_WHEN_OWNING_GTID                              = 1790
	ER_UNKNOWN_EXPLAIN_FORMAT                                           = 1791
	ER_CANT_EXECUTE_IN_READ_ONLY_TRANSACTION                            = 1792
	ER_TOO_LONG_TABLE_PARTITION_COMMENT                                 = 1793
	ER_SLAVE_CONFIGURATION                                              = 1794
	ER_INNODB_FT_LIMIT                                                  = 1795
	ER_INNODB_NO_FT_TEMP_TABLE                                          = 1796
	ER_INNODB_FT_WRONG_DOCID_COLUMN                                     = 1797
	ER_INNODB_FT_WRONG_DOCID_INDEX                                      = 1798
	ER_INNODB_ONLINE_LOG_TOO_BIG                                        = 1799
	ER_UNKNOWN_ALTER_ALGORITHM                                          = 1800
	ER_UNKNOWN_ALTER_LOCK                                               = 1801
	ER_MTS_CHANGE_MASTER_CANT_RUN_WITH_GAPS                             = 1802
	ER_MTS_RECOVERY_FAILURE                                             = 1803
	ER_MTS_RESET_WORKERS                                                = 1804
	ER_COL_COUNT_DOESNT_MATCH_CORRUPTED_V2                              = 1805
	ER_SLAVE_SILENT_RETRY_TRANSACTION                                   = 1806
	ER_DISCARD_FK_CHECKS_RUNNING                                        = 1807
	ER_TABLE_SCHEMA_MISMATCH                                            = 1808
	ER_TABLE_IN_SYSTEM_TABLESPACE                                       = 1809
	ER_IO_READ_ERROR                                                    = 1810
	ER_IO_WRITE_ERROR                                                   = 1811
	ER_TABLESPACE_MISSING                                               = 1812
	ER_TABLESPACE_EXISTS                                                = 1813
	ER_TABLESPACE_DISCARDED                                             = 1814
	ER_INTERNAL_ERROR                                                   = 1815
	ER_INNODB_IMPORT_ERROR                                              = 1816
	ER_INNODB_INDEX_CORRUPT                                             = 1817
	ER_INVALID_YEAR_COLUMN_LENGTH                                       = 1818
	ER_NOT_VALID_PASSWORD                                               = 1819
	ER_MUST_CHANGE_PASSWORD                                             = 1820
	ER_FK_NO_INDEX_CHILD                                                = 1821
	ER_FK_NO_INDEX_PARENT                                               = 1822
	ER_FK_FAIL_ADD_SYSTEM                                               = 1823
	ER_FK_CANNOT_OPEN_PARENT                                            = 1824
	ER_FK_INCORRECT_OPTION                                              = 1825
	ER_FK_DUP_NAME                                                      = 1826
	ER_PASSWORD_FORMAT                                                  = 1827
	ER_FK_COLUMN_CANNOT_DROP                                            = 1828
	ER_FK_COLUMN_CANNOT_DROP_CHILD                                      = 1829
	ER_FK_COLUMN_NOT_NULL                                               = 1830
	ER_DUP_INDEX                                                        = 1831
	ER_FK_COLUMN_CANNOT_CHANGE                                          = 1832
	ER_FK_COLUMN_CANNOT_CHANGE_CHILD                                    = 1833
	ER_FK_CANNOT_DELETE_PARENT                                          = 1834
	ER_MALFORMED_PACKET                                                 = 1835
	ER_READ_ONLY_MODE                                                   = 1836
	ER_GTID_NEXT_TYPE_UNDEFINED_GROUP                                   = 1837
	ER_VARIABLE_NOT_SETTABLE_IN_SP                                      = 1838
	ER_CANT_SET_GTID_PURGED_WHEN_GTID_MODE_IS_OFF                       = 1839
	ER_CANT_SET_GTID_PURGED_WHEN_GTID_EXECUTED_IS_NOT_EMPTY             = 1840
	ER_CANT_SET_GTID_PURGED_WHEN_OWNED_GTIDS_IS_NOT_EMPTY               = 1841
	ER_GTID_PURGED_WAS_CHANGED                                          = 1842
	ER_GTID_EXECUTED_WAS_CHANGED                                        = 1843
	ER_BINLOG_STMT_MODE_AND_NO_REPL_TABLES                              = 1844
	ER_ALTER_OPERATION_NOT_SUPPORTED                                    = 1845
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON                             = 1846
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_COPY                        = 1847
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_PARTITION                   = 1848
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_FK_RENAME                   = 1849
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_COLUMN_TYPE                 = 1850
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_FK_CHECK                    = 1851
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_IGNORE                      = 1852
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_NOPK                        = 1853
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_AUTOINC                     = 1854
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_HIDDEN_FTS                  = 1855
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_CHANGE_FTS                  = 1856
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_FTS                         = 1857
	ER_SQL_SLAVE_SKIP_COUNTER_NOT_SETTABLE_IN_GTID_MODE                 = 1858
	ER_DUP_UNKNOWN_IN_INDEX                                             = 1859
	ER_IDENT_CAUSES_TOO_LONG_PATH                                       = 1860
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_NOT_NULL                    = 1861
	ER_MUST_CHANGE_PASSWORD_LOGIN                                       = 1862
	ER_ROW_IN_WRONG_PARTITION                                           = 1863
	ER_ERROR_LAST                                                       = 1863
)

var MySQLErrName = map[uint16]string{
	ER_HASHCHK:                                       "hashchk",
	ER_NISAMCHK:                                      "isamchk",
	ER_NO:                                            "NO",
	ER_YES:                                           "YES",
	ER_CANT_CREATE_FILE:                              "Can't create file '%-.200s' (errno: %d - %s)",
	ER_CANT_CREATE_TABLE:                             "Can't create table '%-.200s' (errno: %d)",
	ER_CANT_CREATE_DB:                                "Can't create database '%-.192s' (errno: %d)",
	ER_DB_CREATE_EXISTS:                              "Can't create database '%-.192s'; database exists",
	ER_DB_DROP_EXISTS:                                "Can't drop database '%-.192s'; database doesn't exist",
	ER_DB_DROP_DELETE:                                "Error dropping database (can't delete '%-.192s', errno: %d)",
	ER_DB_DROP_RMDIR:                                 "Error dropping database (can't rmdir '%-.192s', errno: %d)",
	ER_CANT_DELETE_FILE:                              "Error on delete of '%-.192s' (errno: %d - %s)",
	ER_CANT_FIND_SYSTEM_REC:                          "Can't read record in system table",
	ER_CANT_GET_STAT:                                 "Can't get status of '%-.200s' (errno: %d - %s)",
	ER_CANT_GET_WD:                                   "Can't get working directory (errno: %d - %s)",
	ER_CANT_LOCK:                                     "Can't lock file (errno: %d - %s)",
	ER_CANT_OPEN_FILE:                                "Can't open file: '%-.200s' (errno: %d - %s)",
	ER_FILE_NOT_FOUND:                                "Can't find file: '%-.200s' (errno: %d - %s)",
	ER_CANT_READ_DIR:                                 "Can't read dir of '%-.192s' (errno: %d - %s)",
	ER_CANT_SET_WD:                                   "Can't change dir to '%-.192s' (errno: %d - %s)",
	ER_CHECKREAD:                                     "Record has changed since last read in table '%-.192s'",
	ER_DISK_FULL:                                     "Disk full (%s); waiting for someone to free some space... (errno: %d - %s)",
	ER_DUP_KEY:                                       "Can't write; duplicate key in table '%-.192s'",
	ER_ERROR_ON_CLOSE:                                "Error on close of '%-.192s' (errno: %d - %s)",
	ER_ERROR_ON_READ:                                 "Error reading file '%-.200s' (errno: %d - %s)",
	ER_ERROR_ON_RENAME:                               "Error on rename of '%-.210s' to '%-.210s' (errno: %d - %s)",
	ER_ERROR_ON_WRITE:                                "Error writing file '%-.200s' (errno: %d - %s)",
	ER_FILE_USED:                                     "'%-.192s' is locked against change",
	ER_FILSORT_ABORT:                                 "Sort aborted",
	ER_FORM_NOT_FOUND:                                "View '%-.192s' doesn't exist for '%-.192s'",
	ER_GET_ERRNO:                                     "Got error %d from storage engine",
	ER_ILLEGAL_HA:                                    "Table storage engine for '%-.192s' doesn't have this option",
	ER_KEY_NOT_FOUND:                                 "Can't find record in '%-.192s'",
	ER_NOT_FORM_FILE:                                 "Incorrect information in file: '%-.200s'",
	ER_NOT_KEYFILE:                                   "Incorrect key file for table '%-.200s'; try to repair it",
	ER_OLD_KEYFILE:                                   "Old key file for table '%-.192s'; repair it!",
	ER_OPEN_AS_READONLY:                              "Table '%-.192s' is read only",
	ER_OUTOFMEMORY:                                   "Out of memory; restart server and try again (needed %d bytes)",
	ER_OUT_OF_SORTMEMORY:                             "Out of sort memory, consider increasing server sort buffer size",
	ER_UNEXPECTED_EOF:                                "Unexpected EOF found when reading file '%-.192s' (errno: %d - %s)",
	ER_CON_COUNT_ERROR:                               "Too many connections",
	ER_OUT_OF_RESOURCES:                              "Out of memory; check if mysqld or some other process uses all available memory; if not, you may have to use 'ulimit' to allow mysqld to use more memory or you can add more swap space",
	ER_BAD_HOST_ERROR:                                "Can't get hostname for your address",
	ER_HANDSHAKE_ERROR:                               "Bad handshake",
	ER_DBACCESS_DENIED_ERROR:                         "Access denied for user '%-.48s'@'%-.64s' to database '%-.192s'",
	ER_ACCESS_DENIED_ERROR:                           "Access denied for user '%-.48s'@'%-.64s' (using password: %s)",
	ER_NO_DB_ERROR:                                   "No database selected",
	ER_UNKNOWN_COM_ERROR:                             "Unknown command",
	ER_BAD_NULL_ERROR:                                "Column '%-.192s' cannot be null",
	ER_BAD_DB_ERROR:                                  "Unknown database '%-.192s'",
	ER_TABLE_EXISTS_ERROR:                            "Table '%-.192s' already exists",
	ER_BAD_TABLE_ERROR:                               "Unknown table '%-.100s'",
	ER_NON_UNIQ_ERROR:                                "Column '%-.192s' in %-.192s is ambiguous",
	ER_SERVER_SHUTDOWN:                               "Server shutdown in progress",
	ER_BAD_FIELD_ERROR:                               "Unknown column '%-.192s' in '%-.192s'",
	ER_WRONG_FIELD_WITH_GROUP:                        "'%-.192s' isn't in GROUP BY",
	ER_WRONG_GROUP_FIELD:                             "Can't group on '%-.192s'",
	ER_WRONG_SUM_SELECT:                              "Statement has sum functions and columns in same statement",
	ER_WRONG_VALUE_COUNT:                             "Column count doesn't match value count",
	ER_TOO_LONG_IDENT:                                "Identifier name '%-.100s' is too long",
	ER_DUP_FIELDNAME:                                 "Duplicate column name '%-.192s'",
	ER_DUP_KEYNAME:                                   "Duplicate key name '%-.192s'",
	ER_DUP_ENTRY:                                     "Duplicate entry '%-.192s' for key %d",
	ER_WRONG_FIELD_SPEC:                              "Incorrect column specifier for column '%-.192s'",
	ER_PARSE_ERROR:                                   "%s near '%-.80s' at line %d",
	ER_EMPTY_QUERY:                                   "Query was empty",
	ER_NONUNIQ_TABLE:                                 "Not unique table/alias: '%-.192s'",
	ER_INVALID_DEFAULT:                               "Invalid default value for '%-.192s'",
	ER_MULTIPLE_PRI_KEY:                              "Multiple primary key defined",
	ER_TOO_MANY_KEYS:                                 "Too many keys specified; max %d keys allowed",
	ER_TOO_MANY_KEY_PARTS:                            "Too many key parts specified; max %d parts allowed",
	ER_TOO_LONG_KEY:                                  "Specified key was too long; max key length is %d bytes",
	ER_KEY_COLUMN_DOES_NOT_EXITS:                     "Key column '%-.192s' doesn't exist in table",
	ER_BLOB_USED_AS_KEY:                              "BLOB column '%-.192s' can't be used in key specification with the used table type",
	ER_TOO_BIG_FIELDLENGTH:                           "Column length too big for column '%-.192s' (max = %d); use BLOB or TEXT instead",
	ER_WRONG_AUTO_KEY:                                "Incorrect table definition; there can be only one auto column and it must be defined as a key",
	ER_READY:                                         "%s: ready for connections.\nVersion: '%s'  socket: '%s'  port: %d",
	ER_NORMAL_SHUTDOWN:                               "%s: Normal shutdown\n",
	ER_GOT_SIGNAL:                                    "%s: Got signal %d. Aborting!\n",
	ER_SHUTDOWN_COMPLETE:                             "%s: Shutdown complete\n",
	ER_FORCING_CLOSE:                                 "%s: Forcing close of thread %d  user: '%-.48s'\n",
	ER_IPSOCK_ERROR:                                  "Can't create IP socket",
	ER_NO_SUCH_INDEX:                                 "Table '%-.192s' has no index like the one used in CREATE INDEX; recreate the table",
	ER_WRONG_FIELD_TERMINATORS:                       "Field separator argument is not what is expected; check the manual",
	ER_BLOBS_AND_NO_TERMINATED:                       "You can't use fixed rowlength with BLOBs; please use 'fields terminated by'",
	ER_TEXTFILE_NOT_READABLE:                         "The file '%-.128s' must be in the database directory or be readable by all",
	ER_FILE_EXISTS_ERROR:                             "File '%-.200s' already exists",
	ER_LOAD_INFO:                                     "Records: %d  Deleted: %d  Skipped: %d  Warnings: %d",
	ER_ALTER_INFO:                                    "Records: %d  Duplicates: %d",
	ER_WRONG_SUB_KEY:                                 "Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys",
	ER_CANT_REMOVE_ALL_FIELDS:                        "You can't delete all columns with ALTER TABLE; use DROP TABLE instead",
	ER_CANT_DROP_FIELD_OR_KEY:                        "Can't DROP '%-.192s'; check that column/key exists",
	ER_INSERT_INFO:                                   "Records: %d  Duplicates: %d  Warnings: %d",
	ER_UPDATE_TABLE_USED:                             "You can't specify target table '%-.192s' for update in FROM clause",
	ER_NO_SUCH_THREAD:                                "Unknown thread id: %d",
	ER_KILL_DENIED_ERROR:                             "You are not owner of thread %d",
	ER_NO_TABLES_USED:                                "No tables used",
	ER_TOO_BIG_SET:                                   "Too many strings for column %-.192s and SET",
	ER_NO_UNIQUE_LOGFILE:                             "Can't generate a unique log-filename %-.200s.(1-999)\n",
	ER_TABLE_NOT_LOCKED_FOR_WRITE:                    "Table '%-.192s' was locked with a READ lock and can't be updated",
	ER_TABLE_NOT_LOCKED:                              "Table '%-.192s' was not locked with LOCK TABLES",
	ER_BLOB_CANT_HAVE_DEFAULT:                        "BLOB/TEXT column '%-.192s' can't have a default value",
	ER_WRONG_DB_NAME:                                 "Incorrect database name '%-.100s'",
	ER_WRONG_TABLE_NAME:                              "Incorrect table name '%-.100s'",
	ER_TOO_BIG_SELECT:                                "The SELECT would examine more than MAX_JOIN_SIZE rows; check your WHERE and use SET SQL_BIG_SELECTS=1 or SET MAX_JOIN_SIZE=# if the SELECT is okay",
	ER_UNKNOWN_ERROR:                                 "Unknown error",
	ER_UNKNOWN_PROCEDURE:                             "Unknown procedure '%-.192s'",
	ER_WRONG_PARAMCOUNT_TO_PROCEDURE:                 "Incorrect parameter count to procedure '%-.192s'",
	ER_WRONG_PARAMETERS_TO_PROCEDURE:                 "Incorrect parameters to procedure '%-.192s'",
	ER_UNKNOWN_TABLE:                                 "Unknown table '%-.192s' in %-.32s",
	ER_FIELD_SPECIFIED_TWICE:                         "Column '%-.192s' specified twice",
	ER_INVALID_GROUP_FUNC_USE:                        "Invalid use of group function",
	ER_UNSUPPORTED_EXTENSION:                         "Table '%-.192s' uses an extension that doesn't exist in this MySQL version",
	ER_TABLE_MUST_HAVE_COLUMNS:                       "A table must have at least 1 column",
	ER_RECORD_FILE_FULL:                              "The table '%-.192s' is full",
	ER_UNKNOWN_CHARACTER_SET:                         "Unknown character set: '%-.64s'",
	ER_TOO_MANY_TABLES:                               "Too many tables; MySQL can only use %d tables in a join",
	ER_TOO_MANY_FIELDS:                               "Too many columns",
	ER_TOO_BIG_ROWSIZE:                               "Row size too large. The maximum row size for the used table type, not counting BLOBs, is %d. This includes storage overhead, check the manual. You have to change some columns to TEXT or BLOBs",
	ER_STACK_OVERRUN:                                 "Thread stack overrun:  Used: %d of a %d stack.  Use 'mysqld --thread_stack=#' to specify a bigger stack if needed",
	ER_WRONG_OUTER_JOIN:                              "Cross dependency found in OUTER JOIN; examine your ON conditions",
	ER_NULL_COLUMN_IN_INDEX:                          "Table handler doesn't support NULL in given index. Please change column '%-.192s' to be NOT NULL or use another handler",
	ER_CANT_FIND_UDF:                                 "Can't load function '%-.192s'",
	ER_CANT_INITIALIZE_UDF:                           "Can't initialize function '%-.192s'; %-.80s",
	ER_UDF_NO_PATHS:                                  "No paths allowed for shared library",
	ER_UDF_EXISTS:                                    "Function '%-.192s' already exists",
	ER_CANT_OPEN_LIBRARY:                             "Can't open shared library '%-.192s' (errno: %d %-.128s)",
	ER_CANT_FIND_DL_ENTRY:                            "Can't find symbol '%-.128s' in library",
	ER_FUNCTION_NOT_DEFINED:                          "Function '%-.192s' is not defined",
	ER_HOST_IS_BLOCKED:                               "Host '%-.64s' is blocked because of many connection errors; unblock with 'mysqladmin flush-hosts'",
	ER_HOST_NOT_PRIVILEGED:                           "Host '%-.64s' is not allowed to connect to this MySQL server",
	ER_PASSWORD_ANONYMOUS_USER:                       "You are using MySQL as an anonymous user and anonymous users are not allowed to change passwords",
	ER_PASSWORD_NOT_ALLOWED:                          "You must have privileges to update tables in the mysql database to be able to change passwords for others",
	ER_PASSWORD_NO_MATCH:                             "Can't find any matching row in the user table",
	ER_UPDATE_INFO:                                   "Rows matched: %d  Changed: %d  Warnings: %d",
	ER_CANT_CREATE_THREAD:                            "Can't create a new thread (errno %d); if you are not out of available memory, you can consult the manual for a possible OS-dependent bug",
	ER_WRONG_VALUE_COUNT_ON_ROW:                      "Column count doesn't match value count at row %d",
	ER_CANT_REOPEN_TABLE:                             "Can't reopen table: '%-.192s'",
	ER_INVALID_USE_OF_NULL:                           "Invalid use of NULL value",
	ER_REGEXP_ERROR:                                  "Got error '%-.64s' from regexp",
	ER_MIX_OF_GROUP_FUNC_AND_FIELDS:                  "Mixing of GROUP columns (MIN(),MAX(),COUNT(),...) with no GROUP columns is illegal if there is no GROUP BY clause",
	ER_NONEXISTING_GRANT:                             "There is no such grant defined for user '%-.48s' on host '%-.64s'",
	ER_TABLEACCESS_DENIED_ERROR:                      "%-.128s command denied to user '%-.48s'@'%-.64s' for table '%-.64s'",
	ER_COLUMNACCESS_DENIED_ERROR:                     "%-.16s command denied to user '%-.48s'@'%-.64s' for column '%-.192s' in table '%-.192s'",
	ER_ILLEGAL_GRANT_FOR_TABLE:                       "Illegal GRANT/REVOKE command; please consult the manual to see which privileges can be used",
	ER_GRANT_WRONG_HOST_OR_USER:                      "The host or user argument to GRANT is too long",
	ER_NO_SUCH_TABLE:                                 "Table '%-.192s.%-.192s' doesn't exist",
	ER_NONEXISTING_TABLE_GRANT:                       "There is no such grant defined for user '%-.48s' on host '%-.64s' on table '%-.192s'",
	ER_NOT_ALLOWED_COMMAND:                           "The used command is not allowed with this MySQL version",
	ER_SYNTAX_ERROR:                                  "You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use",
	ER_DELAYED_CANT_CHANGE_LOCK:                      "Delayed insert thread couldn't get requested lock for table %-.192s",
	ER_TOO_MANY_DELAYED_THREADS:                      "Too many delayed threads in use",
	ER_ABORTING_CONNECTION:                           "Aborted connection %d to db: '%-.192s' user: '%-.48s' (%-.64s)",
	ER_NET_PACKET_TOO_LARGE:                          "Got a packet bigger than 'max_allowed_packet' bytes",
	ER_NET_READ_ERROR_FROM_PIPE:                      "Got a read error from the connection pipe",
	ER_NET_FCNTL_ERROR:                               "Got an error from fcntl()",
	ER_NET_PACKETS_OUT_OF_ORDER:                      "Got packets out of order",
	ER_NET_UNCOMPRESS_ERROR:                          "Couldn't uncompress communication packet",
	ER_NET_READ_ERROR:                                "Got an error reading communication packets",
	ER_NET_READ_INTERRUPTED:                          "Got timeout reading communication packets",
	ER_NET_ERROR_ON_WRITE:                            "Got an error writing communication packets",
	ER_NET_WRITE_INTERRUPTED:                         "Got timeout writing communication packets",
	ER_TOO_LONG_STRING:                               "Result string is longer than 'max_allowed_packet' bytes",
	ER_TABLE_CANT_HANDLE_BLOB:                        "The used table type doesn't support BLOB/TEXT columns",
	ER_TABLE_CANT_HANDLE_AUTO_INCREMENT:              "The used table type doesn't support AUTO_INCREMENT columns",
	ER_DELAYED_INSERT_TABLE_LOCKED:                   "INSERT DELAYED can't be used with table '%-.192s' because it is locked with LOCK TABLES",
	ER_WRONG_COLUMN_NAME:                             "Incorrect column name '%-.100s'",
	ER_WRONG_KEY_COLUMN:                              "The used storage engine can't index column '%-.192s'",
	ER_WRONG_MRG_TABLE:                               "Unable to open underlying table which is differently defined or of non-MyISAM type or doesn't exist",
	ER_DUP_UNIQUE:                                    "Can't write, because of unique constraint, to table '%-.192s'",
	ER_BLOB_KEY_WITHOUT_LENGTH:                       "BLOB/TEXT column '%-.192s' used in key specification without a key length",
	ER_PRIMARY_CANT_HAVE_NULL:                        "All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead",
	ER_TOO_MANY_ROWS:                                 "Result consisted of more than one row",
	ER_REQUIRES_PRIMARY_KEY:                          "This table type requires a primary key",
	ER_NO_RAID_COMPILED:                              "This version of MySQL is not compiled with RAID support",
	ER_UPDATE_WITHOUT_KEY_IN_SAFE_MODE:               "You are using safe update mode and you tried to update a table without a WHERE that uses a KEY column",
	ER_KEY_DOES_NOT_EXITS:                            "Key '%-.192s' doesn't exist in table '%-.192s'",
	ER_CHECK_NO_SUCH_TABLE:                           "Can't open table",
	ER_CHECK_NOT_IMPLEMENTED:                         "The storage engine for the table doesn't support %s",
	ER_CANT_DO_THIS_DURING_AN_TRANSACTION:            "You are not allowed to execute this command in a transaction",
	ER_ERROR_DURING_COMMIT:                           "Got error %d during COMMIT",
	ER_ERROR_DURING_ROLLBACK:                         "Got error %d during ROLLBACK",
	ER_ERROR_DURING_FLUSH_LOGS:                       "Got error %d during FLUSH_LOGS",
	ER_ERROR_DURING_CHECKPOINT:                       "Got error %d during CHECKPOINT",
	ER_NEW_ABORTING_CONNECTION:                       "Aborted connection %d to db: '%-.192s' user: '%-.48s' host: '%-.64s' (%-.64s)",
	ER_DUMP_NOT_IMPLEMENTED:                          "The storage engine for the table does not support binary table dump",
	ER_FLUSH_MASTER_BINLOG_CLOSED:                    "Binlog closed, cannot RESET MASTER",
	ER_INDEX_REBUILD:                                 "Failed rebuilding the index of  dumped table '%-.192s'",
	ER_MASTER:                                        "Error from master: '%-.64s'",
	ER_MASTER_NET_READ:                               "Net error reading from master",
	ER_MASTER_NET_WRITE:                              "Net error writing to master",
	ER_FT_MATCHING_KEY_NOT_FOUND:                     "Can't find FULLTEXT index matching the column list",
	ER_LOCK_OR_ACTIVE_TRANSACTION:                    "Can't execute the given command because you have active locked tables or an active transaction",
	ER_UNKNOWN_SYSTEM_VARIABLE:                       "Unknown system variable '%-.64s'",
	ER_CRASHED_ON_USAGE:                              "Table '%-.192s' is marked as crashed and should be repaired",
	ER_CRASHED_ON_REPAIR:                             "Table '%-.192s' is marked as crashed and last (automatic?) repair failed",
	ER_WARNING_NOT_COMPLETE_ROLLBACK:                 "Some non-transactional changed tables couldn't be rolled back",
	ER_TRANS_CACHE_FULL:                              "Multi-statement transaction required more than 'max_binlog_cache_size' bytes of storage; increase this mysqld variable and try again",
	ER_SLAVE_MUST_STOP:                               "This operation cannot be performed with a running slave; run STOP SLAVE first",
	ER_SLAVE_NOT_RUNNING:                             "This operation requires a running slave; configure slave and do START SLAVE",
	ER_BAD_SLAVE:                                     "The server is not configured as slave; fix in config file or with CHANGE MASTER TO",
	ER_MASTER_INFO:                                   "Could not initialize master info structure; more error messages can be found in the MySQL error log",
	ER_SLAVE_THREAD:                                  "Could not create slave thread; check system resources",
	ER_TOO_MANY_USER_CONNECTIONS:                     "User %-.64s already has more than 'max_user_connections' active connections",
	ER_SET_CONSTANTS_ONLY:                            "You may only use constant expressions with SET",
	ER_LOCK_WAIT_TIMEOUT:                             "Lock wait timeout exceeded; try restarting transaction",
	ER_LOCK_TABLE_FULL:                               "The total number of locks exceeds the lock table size",
	ER_READ_ONLY_TRANSACTION:                         "Update locks cannot be acquired during a READ UNCOMMITTED transaction",
	ER_DROP_DB_WITH_READ_LOCK:                        "DROP DATABASE not allowed while thread is holding global read lock",
	ER_CREATE_DB_WITH_READ_LOCK:                      "CREATE DATABASE not allowed while thread is holding global read lock",
	ER_WRONG_ARGUMENTS:                               "Incorrect arguments to %s",
	ER_NO_PERMISSION_TO_CREATE_USER:                  "'%-.48s'@'%-.64s' is not allowed to create new users",
	ER_UNION_TABLES_IN_DIFFERENT_DIR:                 "Incorrect table definition; all MERGE tables must be in the same database",
	ER_LOCK_DEADLOCK:                                 "Deadlock found when trying to get lock; try restarting transaction",
	ER_TABLE_CANT_HANDLE_FT:                          "The used table type doesn't support FULLTEXT indexes",
	ER_CANNOT_ADD_FOREIGN:                            "Cannot add foreign key constraint",
	ER_NO_REFERENCED_ROW:                             "Cannot add or update a child row: a foreign key constraint fails",
	ER_ROW_IS_REFERENCED:                             "Cannot delete or update a parent row: a foreign key constraint fails",
	ER_CONNECT_TO_MASTER:                             "Error connecting to master: %-.128s",
	ER_QUERY_ON_MASTER:                               "Error running query on master: %-.128s",
	ER_ERROR_WHEN_EXECUTING_COMMAND:                  "Error when executing command %s: %-.128s",
	ER_WRONG_USAGE:                                   "Incorrect usage of %s and %s",
	ER_WRONG_NUMBER_OF_COLUMNS_IN_SELECT:             "The used SELECT statements have a different number of columns",
	ER_CANT_UPDATE_WITH_READLOCK:                     "Can't execute the query because you have a conflicting read lock",
	ER_MIXING_NOT_ALLOWED:                            "Mixing of transactional and non-transactional tables is disabled",
	ER_DUP_ARGUMENT:                                  "Option '%s' used twice in statement",
	ER_USER_LIMIT_REACHED:                            "User '%-.64s' has exceeded the '%s' resource (current value: %d)",
	ER_SPECIFIC_ACCESS_DENIED_ERROR:                  "Access denied; you need (at least one of) the %-.128s privilege(s) for this operation",
	ER_LOCAL_VARIABLE:                                "Variable '%-.64s' is a SESSION variable and can't be used with SET GLOBAL",
	ER_GLOBAL_VARIABLE:                               "Variable '%-.64s' is a GLOBAL variable and should be set with SET GLOBAL",
	ER_NO_DEFAULT:                                    "Variable '%-.64s' doesn't have a default value",
	ER_WRONG_VALUE_FOR_VAR:                           "Variable '%-.64s' can't be set to the value of '%-.200s'",
	ER_WRONG_TYPE_FOR_VAR:                            "Incorrect argument type to variable '%-.64s'",
	ER_VAR_CANT_BE_READ:                              "Variable '%-.64s' can only be set, not read",
	ER_CANT_USE_OPTION_HERE:                          "Incorrect usage/placement of '%s'",
	ER_NOT_SUPPORTED_YET:                             "This version of MySQL doesn't yet support '%s'",
	ER_MASTER_FATAL_ERROR_READING_BINLOG:             "Got fatal error %d from master when reading data from binary log: '%-.320s'",
	ER_SLAVE_IGNORED_TABLE:                           "Slave SQL thread ignored the query because of replicate-*-table rules",
	ER_INCORRECT_GLOBAL_LOCAL_VAR:                    "Variable '%-.192s' is a %s variable",
	ER_WRONG_FK_DEF:                                  "Incorrect foreign key definition for '%-.192s': %s",
	ER_KEY_REF_DO_NOT_MATCH_TABLE_REF:                "Key reference and table reference don't match",
	ER_OPERAND_COLUMNS:                               "Operand should contain %d column(s)",
	ER_SUBQUERY_NO_1_ROW:                             "Subquery returns more than 1 row",
	ER_UNKNOWN_STMT_HANDLER:                          "Unknown prepared statement handler (%.*s) given to %s",
	ER_CORRUPT_HELP_DB:                               "Help database is corrupt or does not exist",
	ER_CYCLIC_REFERENCE:                              "Cyclic reference on subqueries",
	ER_AUTO_CONVERT:                                  "Converting column '%s' from %s to %s",
	ER_ILLEGAL_REFERENCE:                             "Reference '%-.64s' not supported (%s)",
	ER_DERIVED_MUST_HAVE_ALIAS:                       "Every derived table must have its own alias",
	ER_SELECT_REDUCED:                                "Select %d was reduced during optimization",
	ER_TABLENAME_NOT_ALLOWED_HERE:                    "Table '%-.192s' from one of the SELECTs cannot be used in %-.32s",
	ER_NOT_SUPPORTED_AUTH_MODE:                       "Client does not support authentication protocol requested by server; consider upgrading MySQL client",
	ER_SPATIAL_CANT_HAVE_NULL:                        "All parts of a SPATIAL index must be NOT NULL",
	ER_COLLATION_CHARSET_MISMATCH:                    "COLLATION '%s' is not valid for CHARACTER SET '%s'",
	ER_SLAVE_WAS_RUNNING:                             "Slave is already running",
	ER_SLAVE_WAS_NOT_RUNNING:                         "Slave already has been stopped",
	ER_TOO_BIG_FOR_UNCOMPRESS:                        "Uncompressed data size too large; the maximum size is %d (probably, length of uncompressed data was corrupted)",
	ER_ZLIB_Z_MEM_ERROR:                              "ZLIB: Not enough memory",
	ER_ZLIB_Z_BUF_ERROR:                              "ZLIB: Not enough room in the output buffer (probably, length of uncompressed data was corrupted)",
	ER_ZLIB_Z_DATA_ERROR:                             "ZLIB: Input data corrupted",
	ER_CUT_VALUE_GROUP_CONCAT:                        "Row %d was cut by GROUP_CONCAT()",
	ER_WARN_TOO_FEW_RECORDS:                          "Row %d doesn't contain data for all columns",
	ER_WARN_TOO_MANY_RECORDS:                         "Row %d was truncated; it contained more data than there were input columns",
	ER_WARN_NULL_TO_NOTNULL:                          "Column set to default value; NULL supplied to NOT NULL column '%s' at row %d",
	ER_WARN_DATA_OUT_OF_RANGE:                        "Out of range value for column '%s' at row %d",
	WARN_DATA_TRUNCATED:                              "Data truncated for column '%s' at row %d",
	ER_WARN_USING_OTHER_HANDLER:                      "Using storage engine %s for table '%s'",
	ER_CANT_AGGREGATE_2COLLATIONS:                    "Illegal mix of collations (%s,%s) and (%s,%s) for operation '%s'",
	ER_DROP_USER:                                     "Cannot drop one or more of the requested users",
	ER_REVOKE_GRANTS:                                 "Can't revoke all privileges for one or more of the requested users",
	ER_CANT_AGGREGATE_3COLLATIONS:                    "Illegal mix of collations (%s,%s), (%s,%s), (%s,%s) for operation '%s'",
	ER_CANT_AGGREGATE_NCOLLATIONS:                    "Illegal mix of collations for operation '%s'",
	ER_VARIABLE_IS_NOT_STRUCT:                        "Variable '%-.64s' is not a variable component (can't be used as XXXX.variable_name)",
	ER_UNKNOWN_COLLATION:                             "Unknown collation: '%-.64s'",
	ER_SLAVE_IGNORED_SSL_PARAMS:                      "SSL parameters in CHANGE MASTER are ignored because this MySQL slave was compiled without SSL support; they can be used later if MySQL slave with SSL is started",
	ER_SERVER_IS_IN_SECURE_AUTH_MODE:                 "Server is running in --secure-auth mode, but '%s'@'%s' has a password in the old format; please change the password to the new format",
	ER_WARN_FIELD_RESOLVED:                           "Field or reference '%-.192s%s%-.192s%s%-.192s' of SELECT #%d was resolved in SELECT #%d",
	ER_BAD_SLAVE_UNTIL_COND:                          "Incorrect parameter or combination of parameters for START SLAVE UNTIL",
	ER_MISSING_SKIP_SLAVE:                            "It is recommended to use --skip-slave-start when doing step-by-step replication with START SLAVE UNTIL; otherwise, you will get problems if you get an unexpected slave's mysqld restart",
	ER_UNTIL_COND_IGNORED:                            "SQL thread is not to be started so UNTIL options are ignored",
	ER_WRONG_NAME_FOR_INDEX:                          "Incorrect index name '%-.100s'",
	ER_WRONG_NAME_FOR_CATALOG:                        "Incorrect catalog name '%-.100s'",
	ER_WARN_QC_RESIZE:                                "Query cache failed to set size %d; new query cache size is %d",
	ER_BAD_FT_COLUMN:                                 "Column '%-.192s' cannot be part of FULLTEXT index",
	ER_UNKNOWN_KEY_CACHE:                             "Unknown key cache '%-.100s'",
	ER_WARN_HOSTNAME_WONT_WORK:                       "MySQL is started in --skip-name-resolve mode; you must restart it without this switch for this grant to work",
	ER_UNKNOWN_STORAGE_ENGINE:                        "Unknown storage engine '%s'",
	ER_WARN_DEPRECATED_SYNTAX:                        "'%s' is deprecated and will be removed in a future release. Please use %s instead",
	ER_NON_UPDATABLE_TABLE:                           "The target table %-.100s of the %s is not updatable",
	ER_FEATURE_DISABLED:                              "The '%s' feature is disabled; you need MySQL built with '%s' to have it working",
	ER_OPTION_PREVENTS_STATEMENT:                     "The MySQL server is running with the %s option so it cannot execute this statement",
	ER_DUPLICATED_VALUE_IN_TYPE:                      "Column '%-.100s' has duplicated value '%-.64s' in %s",
	ER_TRUNCATED_WRONG_VALUE:                         "Truncated incorrect %-.32s value: '%-.128s'",
	ER_TOO_MUCH_AUTO_TIMESTAMP_COLS:                  "Incorrect table definition; there can be only one TIMESTAMP column with CURRENT_TIMESTAMP in DEFAULT or ON UPDATE clause",
	ER_INVALID_ON_UPDATE:                             "Invalid ON UPDATE clause for '%-.192s' column",
	ER_UNSUPPORTED_PS:                                "This command is not supported in the prepared statement protocol yet",
	ER_GET_ERRMSG:                                    "Got error %d '%-.100s' from %s",
	ER_GET_TEMPORARY_ERRMSG:                          "Got temporary error %d '%-.100s' from %s",
	ER_UNKNOWN_TIME_ZONE:                             "Unknown or incorrect time zone: '%-.64s'",
	ER_WARN_INVALID_TIMESTAMP:                        "Invalid TIMESTAMP value in column '%s' at row %d",
	ER_INVALID_CHARACTER_STRING:                      "Invalid %s character string: '%.64s'",
	ER_WARN_ALLOWED_PACKET_OVERFLOWED:                "Result of %s() was larger than max_allowed_packet (%d) - truncated",
	ER_CONFLICTING_DECLARATIONS:                      "Conflicting declarations: '%s%s' and '%s%s'",
	ER_SP_NO_RECURSIVE_CREATE:                        "Can't create a %s from within another stored routine",
	ER_SP_ALREADY_EXISTS:                             "%s %s already exists",
	ER_SP_DOES_NOT_EXIST:                             "%s %s does not exist",
	ER_SP_DROP_FAILED:                                "Failed to DROP %s %s",
	ER_SP_STORE_FAILED:                               "Failed to CREATE %s %s",
	ER_SP_LILABEL_MISMATCH:                           "%s with no matching label: %s",
	ER_SP_LABEL_REDEFINE:                             "Redefining label %s",
	ER_SP_LABEL_MISMATCH:                             "End-label %s without match",
	ER_SP_UNINIT_VAR:                                 "Referring to uninitialized variable %s",
	ER_SP_BADSELECT:                                  "PROCEDURE %s can't return a result set in the given context",
	ER_SP_BADRETURN:                                  "RETURN is only allowed in a FUNCTION",
	ER_SP_BADSTATEMENT:                               "%s is not allowed in stored procedures",
	ER_UPDATE_LOG_DEPRECATED_IGNORED:                 "The update log is deprecated and replaced by the binary log; SET SQL_LOG_UPDATE has been ignored.",
	ER_UPDATE_LOG_DEPRECATED_TRANSLATED:              "The update log is deprecated and replaced by the binary log; SET SQL_LOG_UPDATE has been translated to SET SQL_LOG_BIN.",
	ER_QUERY_INTERRUPTED:                             "Query execution was interrupted",
	ER_SP_WRONG_NO_OF_ARGS:                           "Incorrect number of arguments for %s %s; expected %d, got %d",
	ER_SP_COND_MISMATCH:                              "Undefined CONDITION: %s",
	ER_SP_NORETURN:                                   "No RETURN found in FUNCTION %s",
	ER_SP_NORETURNEND:                                "FUNCTION %s ended without RETURN",
	ER_SP_BAD_CURSOR_QUERY:                           "Cursor statement must be a SELECT",
	ER_SP_BAD_CURSOR_SELECT:                          "Cursor SELECT must not have INTO",
	ER_SP_CURSOR_MISMATCH:                            "Undefined CURSOR: %s",
	ER_SP_CURSOR_ALREADY_OPEN:                        "Cursor is already open",
	ER_SP_CURSOR_NOT_OPEN:                            "Cursor is not open",
	ER_SP_UNDECLARED_VAR:                             "Undeclared variable: %s",
	ER_SP_WRONG_NO_OF_FETCH_ARGS:                     "Incorrect number of FETCH variables",
	ER_SP_FETCH_NO_DATA:                              "No data - zero rows fetched, selected, or processed",
	ER_SP_DUP_PARAM:                                  "Duplicate parameter: %s",
	ER_SP_DUP_VAR:                                    "Duplicate variable: %s",
	ER_SP_DUP_COND:                                   "Duplicate condition: %s",
	ER_SP_DUP_CURS:                                   "Duplicate cursor: %s",
	ER_SP_CANT_ALTER:                                 "Failed to ALTER %s %s",
	ER_SP_SUBSELECT_NYI:                              "Subquery value not supported",
	ER_STMT_NOT_ALLOWED_IN_SF_OR_TRG:                 "%s is not allowed in stored function or trigger",
	ER_SP_VARCOND_AFTER_CURSHNDLR:                    "Variable or condition declaration after cursor or handler declaration",
	ER_SP_CURSOR_AFTER_HANDLER:                       "Cursor declaration after handler declaration",
	ER_SP_CASE_NOT_FOUND:                             "Case not found for CASE statement",
	ER_FPARSER_TOO_BIG_FILE:                          "Configuration file '%-.192s' is too big",
	ER_FPARSER_BAD_HEADER:                            "Malformed file type header in file '%-.192s'",
	ER_FPARSER_EOF_IN_COMMENT:                        "Unexpected end of file while parsing comment '%-.200s'",
	ER_FPARSER_ERROR_IN_PARAMETER:                    "Error while parsing parameter '%-.192s' (line: '%-.192s')",
	ER_FPARSER_EOF_IN_UNKNOWN_PARAMETER:              "Unexpected end of file while skipping unknown parameter '%-.192s'",
	ER_VIEW_NO_EXPLAIN:                               "EXPLAIN/SHOW can not be issued; lacking privileges for underlying table",
	ER_FRM_UNKNOWN_TYPE:                              "File '%-.192s' has unknown type '%-.64s' in its header",
	ER_WRONG_OBJECT:                                  "'%-.192s.%-.192s' is not %s",
	ER_NONUPDATEABLE_COLUMN:                          "Column '%-.192s' is not updatable",
	ER_VIEW_SELECT_DERIVED:                           "View's SELECT contains a subquery in the FROM clause",
	ER_VIEW_SELECT_CLAUSE:                            "View's SELECT contains a '%s' clause",
	ER_VIEW_SELECT_VARIABLE:                          "View's SELECT contains a variable or parameter",
	ER_VIEW_SELECT_TMPTABLE:                          "View's SELECT refers to a temporary table '%-.192s'",
	ER_VIEW_WRONG_LIST:                               "View's SELECT and view's field list have different column counts",
	ER_WARN_VIEW_MERGE:                               "View merge algorithm can't be used here for now (assumed undefined algorithm)",
	ER_WARN_VIEW_WITHOUT_KEY:                         "View being updated does not have complete key of underlying table in it",
	ER_VIEW_INVALID:                                  "View '%-.192s.%-.192s' references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them",
	ER_SP_NO_DROP_SP:                                 "Can't drop or alter a %s from within another stored routine",
	ER_SP_GOTO_IN_HNDLR:                              "GOTO is not allowed in a stored procedure handler",
	ER_TRG_ALREADY_EXISTS:                            "Trigger already exists",
	ER_TRG_DOES_NOT_EXIST:                            "Trigger does not exist",
	ER_TRG_ON_VIEW_OR_TEMP_TABLE:                     "Trigger's '%-.192s' is view or temporary table",
	ER_TRG_CANT_CHANGE_ROW:                           "Updating of %s row is not allowed in %strigger",
	ER_TRG_NO_SUCH_ROW_IN_TRG:                        "There is no %s row in %s trigger",
	ER_NO_DEFAULT_FOR_FIELD:                          "Field '%-.192s' doesn't have a default value",
	ER_DIVISION_BY_ZERO:                              "Division by 0",
	ER_TRUNCATED_WRONG_VALUE_FOR_FIELD:               "Incorrect %-.32s value: '%-.128s' for column '%.192s' at row %d",
	ER_ILLEGAL_VALUE_FOR_TYPE:                        "Illegal %s '%-.192s' value found during parsing",
	ER_VIEW_NONUPD_CHECK:                             "CHECK OPTION on non-updatable view '%-.192s.%-.192s'",
	ER_VIEW_CHECK_FAILED:                             "CHECK OPTION failed '%-.192s.%-.192s'",
	ER_PROCACCESS_DENIED_ERROR:                       "%-.16s command denied to user '%-.48s'@'%-.64s' for routine '%-.192s'",
	ER_RELAY_LOG_FAIL:                                "Failed purging old relay logs: %s",
	ER_PASSWD_LENGTH:                                 "Password hash should be a %d-digit hexadecimal number",
	ER_UNKNOWN_TARGET_BINLOG:                         "Target log not found in binlog index",
	ER_IO_ERR_LOG_INDEX_READ:                         "I/O error reading log index file",
	ER_BINLOG_PURGE_PROHIBITED:                       "Server configuration does not permit binlog purge",
	ER_FSEEK_FAIL:                                    "Failed on fseek()",
	ER_BINLOG_PURGE_FATAL_ERR:                        "Fatal error during log purge",
	ER_LOG_IN_USE:                                    "A purgeable log is in use, will not purge",
	ER_LOG_PURGE_UNKNOWN_ERR:                         "Unknown error during log purge",
	ER_RELAY_LOG_INIT:                                "Failed initializing relay log position: %s",
	ER_NO_BINARY_LOGGING:                             "You are not using binary logging",
	ER_RESERVED_SYNTAX:                               "The '%-.64s' syntax is reserved for purposes internal to the MySQL server",
	ER_WSAS_FAILED:                                   "WSAStartup Failed",
	ER_DIFF_GROUPS_PROC:                              "Can't handle procedures with different groups yet",
	ER_NO_GROUP_FOR_PROC:                             "Select must have a group with this procedure",
	ER_ORDER_WITH_PROC:                               "Can't use ORDER clause with this procedure",
	ER_LOGGING_PROHIBIT_CHANGING_OF:                  "Binary logging and replication forbid changing the global server %s",
	ER_NO_FILE_MAPPING:                               "Can't map file: %-.200s, errno: %d",
	ER_WRONG_MAGIC:                                   "Wrong magic in %-.64s",
	ER_PS_MANY_PARAM:                                 "Prepared statement contains too many placeholders",
	ER_KEY_PART_0:                                    "Key part '%-.192s' length cannot be 0",
	ER_VIEW_CHECKSUM:                                 "View text checksum failed",
	ER_VIEW_MULTIUPDATE:                              "Can not modify more than one base table through a join view '%-.192s.%-.192s'",
	ER_VIEW_NO_INSERT_FIELD_LIST:                     "Can not insert into join view '%-.192s.%-.192s' without fields list",
	ER_VIEW_DELETE_MERGE_VIEW:                        "Can not delete from join view '%-.192s.%-.192s'",
	ER_CANNOT_USER:                                   "Operation %s failed for %.256s",
	ER_XAER_NOTA:                                     "XAER_NOTA: Unknown XID",
	ER_XAER_INVAL:                                    "XAER_INVAL: Invalid arguments (or unsupported command)",
	ER_XAER_RMFAIL:                                   "XAER_RMFAIL: The command cannot be executed when global transaction is in the  %.64s state",
	ER_XAER_OUTSIDE:                                  "XAER_OUTSIDE: Some work is done outside global transaction",
	ER_XAER_RMERR:                                    "XAER_RMERR: Fatal error occurred in the transaction branch - check your data for consistency",
	ER_XA_RBROLLBACK:                                 "XA_RBROLLBACK: Transaction branch was rolled back",
	ER_NONEXISTING_PROC_GRANT:                        "There is no such grant defined for user '%-.48s' on host '%-.64s' on routine '%-.192s'",
	ER_PROC_AUTO_GRANT_FAIL:                          "Failed to grant EXECUTE and ALTER ROUTINE privileges",
	ER_PROC_AUTO_REVOKE_FAIL:                         "Failed to revoke all privileges to dropped routine",
	ER_DATA_TOO_LONG:                                 "Data too long for column '%s' at row %d",
	ER_SP_BAD_SQLSTATE:                               "Bad SQLSTATE: '%s'",
	ER_STARTUP:                                       "%s: ready for connections.\nVersion: '%s'  socket: '%s'  port: %d  %s",
	ER_LOAD_FROM_FIXED_SIZE_ROWS_TO_VAR:              "Can't load value from file with fixed size rows to variable",
	ER_CANT_CREATE_USER_WITH_GRANT:                   "You are not allowed to create a user with GRANT",
	ER_WRONG_VALUE_FOR_TYPE:                          "Incorrect %-.32s value: '%-.128s' for function %-.32s",
	ER_TABLE_DEF_CHANGED:                             "Table definition has changed, please retry transaction",
	ER_SP_DUP_HANDLER:                                "Duplicate handler declared in the same block",
	ER_SP_NOT_VAR_ARG:                                "OUT or INOUT argument %d for routine %s is not a variable or NEW pseudo-variable in BEFORE trigger",
	ER_SP_NO_RETSET:                                  "Not allowed to return a result set from a %s",
	ER_CANT_CREATE_GEOMETRY_OBJECT:                   "Cannot get geometry object from data you send to the GEOMETRY field",
	ER_FAILED_ROUTINE_BREAK_BINLOG:                   "A routine failed and has neither NO SQL nor READS SQL DATA in its declaration and binary logging is enabled; if non-transactional tables were updated, the binary log will miss their changes",
	ER_BINLOG_UNSAFE_ROUTINE:                         "This function has none of DETERMINISTIC, NO SQL, or READS SQL DATA in its declaration and binary logging is enabled (you *might* want to use the less safe log_bin_trust_function_creators variable)",
	ER_BINLOG_CREATE_ROUTINE_NEED_SUPER:              "You do not have the SUPER privilege and binary logging is enabled (you *might* want to use the less safe log_bin_trust_function_creators variable)",
	ER_EXEC_STMT_WITH_OPEN_CURSOR:                    "You can't execute a prepared statement which has an open cursor associated with it. Reset the statement to re-execute it.",
	ER_STMT_HAS_NO_OPEN_CURSOR:                       "The statement (%d) has no open cursor.",
	ER_COMMIT_NOT_ALLOWED_IN_SF_OR_TRG:               "Explicit or implicit commit is not allowed in stored function or trigger.",
	ER_NO_DEFAULT_FOR_VIEW_FIELD:                     "Field of view '%-.192s.%-.192s' underlying table doesn't have a default value",
	ER_SP_NO_RECURSION:                               "Recursive stored functions and triggers are not allowed.",
	ER_TOO_BIG_SCALE:                                 "Too big scale %d specified for column '%-.192s'. Maximum is %d.",
	ER_TOO_BIG_PRECISION:                             "Too big precision %d specified for column '%-.192s'. Maximum is %d.",
	ER_M_BIGGER_THAN_D:                               "For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '%-.192s').",
	ER_WRONG_LOCK_OF_SYSTEM_TABLE:                    "You can't combine write-locking of system tables with other tables or lock types",
	ER_CONNECT_TO_FOREIGN_DATA_SOURCE:                "Unable to connect to foreign data source: %.64s",
	ER_QUERY_ON_FOREIGN_DATA_SOURCE:                  "There was a problem processing the query on the foreign data source. Data source error: %-.64s",
	ER_FOREIGN_DATA_SOURCE_DOESNT_EXIST:              "The foreign data source you are trying to reference does not exist. Data source error:  %-.64s",
	ER_FOREIGN_DATA_STRING_INVALID_CANT_CREATE:       "Can't create federated table. The data source connection string '%-.64s' is not in the correct format",
	ER_FOREIGN_DATA_STRING_INVALID:                   "The data source connection string '%-.64s' is not in the correct format",
	ER_CANT_CREATE_FEDERATED_TABLE:                   "Can't create federated table. Foreign data src error:  %-.64s",
	ER_TRG_IN_WRONG_SCHEMA:                           "Trigger in wrong schema",
	ER_STACK_OVERRUN_NEED_MORE:                       "Thread stack overrun:  %d bytes used of a %d byte stack, and %d bytes needed.  Use 'mysqld --thread_stack=#' to specify a bigger stack.",
	ER_TOO_LONG_BODY:                                 "Routine body for '%-.100s' is too long",
	ER_WARN_CANT_DROP_DEFAULT_KEYCACHE:               "Cannot drop default keycache",
	ER_TOO_BIG_DISPLAYWIDTH:                          "Display width out of range for column '%-.192s' (max = %d)",
	ER_XAER_DUPID:                                    "XAER_DUPID: The XID already exists",
	ER_DATETIME_FUNCTION_OVERFLOW:                    "Datetime function: %-.32s field overflow",
	ER_CANT_UPDATE_USED_TABLE_IN_SF_OR_TRG:           "Can't update table '%-.192s' in stored function/trigger because it is already used by statement which invoked this stored function/trigger.",
	ER_VIEW_PREVENT_UPDATE:                           "The definition of table '%-.192s' prevents operation %.192s on table '%-.192s'.",
	ER_PS_NO_RECURSION:                               "The prepared statement contains a stored routine call that refers to that same statement. It's not allowed to execute a prepared statement in such a recursive manner",
	ER_SP_CANT_SET_AUTOCOMMIT:                        "Not allowed to set autocommit from a stored function or trigger",
	ER_MALFORMED_DEFINER:                             "Definer is not fully qualified",
	ER_VIEW_FRM_NO_USER:                              "View '%-.192s'.'%-.192s' has no definer information (old table format). Current user is used as definer. Please recreate the view!",
	ER_VIEW_OTHER_USER:                               "You need the SUPER privilege for creation view with '%-.192s'@'%-.192s' definer",
	ER_NO_SUCH_USER:                                  "The user specified as a definer ('%-.64s'@'%-.64s') does not exist",
	ER_FORBID_SCHEMA_CHANGE:                          "Changing schema from '%-.192s' to '%-.192s' is not allowed.",
	ER_ROW_IS_REFERENCED_2:                           "Cannot delete or update a parent row: a foreign key constraint fails (%.192s)",
	ER_NO_REFERENCED_ROW_2:                           "Cannot add or update a child row: a foreign key constraint fails (%.192s)",
	ER_SP_BAD_VAR_SHADOW:                             "Variable '%-.64s' must be quoted with `...`, or renamed",
	ER_TRG_NO_DEFINER:                                "No definer attribute for trigger '%-.192s'.'%-.192s'. The trigger will be activated under the authorization of the invoker, which may have insufficient privileges. Please recreate the trigger.",
	ER_OLD_FILE_FORMAT:                               "'%-.192s' has an old format, you should re-create the '%s' object(s)",
	ER_SP_RECURSION_LIMIT:                            "Recursive limit %d (as set by the max_sp_recursion_depth variable) was exceeded for routine %.192s",
	ER_SP_PROC_TABLE_CORRUPT:                         "Failed to load routine %-.192s. The table mysql.proc is missing, corrupt, or contains bad data (internal code %d)",
	ER_SP_WRONG_NAME:                                 "Incorrect routine name '%-.192s'",
	ER_TABLE_NEEDS_UPGRADE:                           "Table upgrade required. Please do \"REPAIR TABLE `%-.32s`\" or dump/reload to fix it!",
	ER_SP_NO_AGGREGATE:                               "AGGREGATE is not supported for stored functions",
	ER_MAX_PREPARED_STMT_COUNT_REACHED:               "Can't create more than max_prepared_stmt_count statements (current value: %d)",
	ER_VIEW_RECURSIVE:                                "`%-.192s`.`%-.192s` contains view recursion",
	ER_NON_GROUPING_FIELD_USED:                       "Non-grouping field '%-.192s' is used in %-.64s clause",
	ER_TABLE_CANT_HANDLE_SPKEYS:                      "The used table type doesn't support SPATIAL indexes",
	ER_NO_TRIGGERS_ON_SYSTEM_SCHEMA:                  "Triggers can not be created on system tables",
	ER_REMOVED_SPACES:                                "Leading spaces are removed from name '%s'",
	ER_AUTOINC_READ_FAILED:                           "Failed to read auto-increment value from storage engine",
	ER_USERNAME:                                      "user name",
	ER_HOSTNAME:                                      "host name",
	ER_WRONG_STRING_LENGTH:                           "String '%-.70s' is too long for %s (should be no longer than %d)",
	ER_NON_INSERTABLE_TABLE:                          "The target table %-.100s of the %s is not insertable-into",
	ER_ADMIN_WRONG_MRG_TABLE:                         "Table '%-.64s' is differently defined or of non-MyISAM type or doesn't exist",
	ER_TOO_HIGH_LEVEL_OF_NESTING_FOR_SELECT:          "Too high level of nesting for select",
	ER_NAME_BECOMES_EMPTY:                            "Name '%-.64s' has become ''",
	ER_AMBIGUOUS_FIELD_TERM:                          "First character of the FIELDS TERMINATED string is ambiguous; please use non-optional and non-empty FIELDS ENCLOSED BY",
	ER_FOREIGN_SERVER_EXISTS:                         "The foreign server, %s, you are trying to create already exists.",
	ER_FOREIGN_SERVER_DOESNT_EXIST:                   "The foreign server name you are trying to reference does not exist. Data source error:  %-.64s",
	ER_ILLEGAL_HA_CREATE_OPTION:                      "Table storage engine '%-.64s' does not support the create option '%.64s'",
	ER_PARTITION_REQUIRES_VALUES_ERROR:               "Syntax error: %-.64s PARTITIONING requires definition of VALUES %-.64s for each partition",
	ER_PARTITION_WRONG_VALUES_ERROR:                  "Only %-.64s PARTITIONING can use VALUES %-.64s in partition definition",
	ER_PARTITION_MAXVALUE_ERROR:                      "MAXVALUE can only be used in last partition definition",
	ER_PARTITION_SUBPARTITION_ERROR:                  "Subpartitions can only be hash partitions and by key",
	ER_PARTITION_SUBPART_MIX_ERROR:                   "Must define subpartitions on all partitions if on one partition",
	ER_PARTITION_WRONG_NO_PART_ERROR:                 "Wrong number of partitions defined, mismatch with previous setting",
	ER_PARTITION_WRONG_NO_SUBPART_ERROR:              "Wrong number of subpartitions defined, mismatch with previous setting",
	ER_WRONG_EXPR_IN_PARTITION_FUNC_ERROR:            "Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed",
	ER_NO_CONST_EXPR_IN_RANGE_OR_LIST_ERROR:          "Expression in RANGE/LIST VALUES must be constant",
	ER_FIELD_NOT_FOUND_PART_ERROR:                    "Field in list of fields for partition function not found in table",
	ER_LIST_OF_FIELDS_ONLY_IN_HASH_ERROR:             "List of fields is only allowed in KEY partitions",
	ER_INCONSISTENT_PARTITION_INFO_ERROR:             "The partition info in the frm file is not consistent with what can be written into the frm file",
	ER_PARTITION_FUNC_NOT_ALLOWED_ERROR:              "The %-.192s function returns the wrong type",
	ER_PARTITIONS_MUST_BE_DEFINED_ERROR:              "For %-.64s partitions each partition must be defined",
	ER_RANGE_NOT_INCREASING_ERROR:                    "VALUES LESS THAN value must be strictly increasing for each partition",
	ER_INCONSISTENT_TYPE_OF_FUNCTIONS_ERROR:          "VALUES value must be of same type as partition function",
	ER_MULTIPLE_DEF_CONST_IN_LIST_PART_ERROR:         "Multiple definition of same constant in list partitioning",
	ER_PARTITION_ENTRY_ERROR:                         "Partitioning can not be used stand-alone in query",
	ER_MIX_HANDLER_ERROR:                             "The mix of handlers in the partitions is not allowed in this version of MySQL",
	ER_PARTITION_NOT_DEFINED_ERROR:                   "For the partitioned engine it is necessary to define all %-.64s",
	ER_TOO_MANY_PARTITIONS_ERROR:                     "Too many partitions (including subpartitions) were defined",
	ER_SUBPARTITION_ERROR:                            "It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning",
	ER_CANT_CREATE_HANDLER_FILE:                      "Failed to create specific handler file",
	ER_BLOB_FIELD_IN_PART_FUNC_ERROR:                 "A BLOB field is not allowed in partition function",
	ER_UNIQUE_KEY_NEED_ALL_FIELDS_IN_PF:              "A %-.192s must include all columns in the table's partitioning function",
	ER_NO_PARTS_ERROR:                                "Number of %-.64s = 0 is not an allowed value",
	ER_PARTITION_MGMT_ON_NONPARTITIONED:              "Partition management on a not partitioned table is not possible",
	ER_FOREIGN_KEY_ON_PARTITIONED:                    "Foreign key clause is not yet supported in conjunction with partitioning",
	ER_DROP_PARTITION_NON_EXISTENT:                   "Error in list of partitions to %-.64s",
	ER_DROP_LAST_PARTITION:                           "Cannot remove all partitions, use DROP TABLE instead",
	ER_COALESCE_ONLY_ON_HASH_PARTITION:               "COALESCE PARTITION can only be used on HASH/KEY partitions",
	ER_REORG_HASH_ONLY_ON_SAME_NO:                    "REORGANIZE PARTITION can only be used to reorganize partitions not to change their numbers",
	ER_REORG_NO_PARAM_ERROR:                          "REORGANIZE PARTITION without parameters can only be used on auto-partitioned tables using HASH PARTITIONs",
	ER_ONLY_ON_RANGE_LIST_PARTITION:                  "%-.64s PARTITION can only be used on RANGE/LIST partitions",
	ER_ADD_PARTITION_SUBPART_ERROR:                   "Trying to Add partition(s) with wrong number of subpartitions",
	ER_ADD_PARTITION_NO_NEW_PARTITION:                "At least one partition must be added",
	ER_COALESCE_PARTITION_NO_PARTITION:               "At least one partition must be coalesced",
	ER_REORG_PARTITION_NOT_EXIST:                     "More partitions to reorganize than there are partitions",
	ER_SAME_NAME_PARTITION:                           "Duplicate partition name %-.192s",
	ER_NO_BINLOG_ERROR:                               "It is not allowed to shut off binlog on this command",
	ER_CONSECUTIVE_REORG_PARTITIONS:                  "When reorganizing a set of partitions they must be in consecutive order",
	ER_REORG_OUTSIDE_RANGE:                           "Reorganize of range partitions cannot change total ranges except for last partition where it can extend the range",
	ER_PARTITION_FUNCTION_FAILURE:                    "Partition function not supported in this version for this handler",
	ER_PART_STATE_ERROR:                              "Partition state cannot be defined from CREATE/ALTER TABLE",
	ER_LIMITED_PART_RANGE:                            "The %-.64s handler only supports 32 bit integers in VALUES",
	ER_PLUGIN_IS_NOT_LOADED:                          "Plugin '%-.192s' is not loaded",
	ER_WRONG_VALUE:                                   "Incorrect %-.32s value: '%-.128s'",
	ER_NO_PARTITION_FOR_GIVEN_VALUE:                  "Table has no partition for value %-.64s",
	ER_FILEGROUP_OPTION_ONLY_ONCE:                    "It is not allowed to specify %s more than once",
	ER_CREATE_FILEGROUP_FAILED:                       "Failed to create %s",
	ER_DROP_FILEGROUP_FAILED:                         "Failed to drop %s",
	ER_TABLESPACE_AUTO_EXTEND_ERROR:                  "The handler doesn't support autoextend of tablespaces",
	ER_WRONG_SIZE_NUMBER:                             "A size parameter was incorrectly specified, either number or on the form 10M",
	ER_SIZE_OVERFLOW_ERROR:                           "The size number was correct but we don't allow the digit part to be more than 2 billion",
	ER_ALTER_FILEGROUP_FAILED:                        "Failed to alter: %s",
	ER_BINLOG_ROW_LOGGING_FAILED:                     "Writing one row to the row-based binary log failed",
	ER_BINLOG_ROW_WRONG_TABLE_DEF:                    "Table definition on master and slave does not match: %s",
	ER_BINLOG_ROW_RBR_TO_SBR:                         "Slave running with --log-slave-updates must use row-based binary logging to be able to replicate row-based binary log events",
	ER_EVENT_ALREADY_EXISTS:                          "Event '%-.192s' already exists",
	ER_EVENT_STORE_FAILED:                            "Failed to store event %s. Error code %d from storage engine.",
	ER_EVENT_DOES_NOT_EXIST:                          "Unknown event '%-.192s'",
	ER_EVENT_CANT_ALTER:                              "Failed to alter event '%-.192s'",
	ER_EVENT_DROP_FAILED:                             "Failed to drop %s",
	ER_EVENT_INTERVAL_NOT_POSITIVE_OR_TOO_BIG:        "INTERVAL is either not positive or too big",
	ER_EVENT_ENDS_BEFORE_STARTS:                      "ENDS is either invalid or before STARTS",
	ER_EVENT_EXEC_TIME_IN_THE_PAST:                   "Event execution time is in the past. Event has been disabled",
	ER_EVENT_OPEN_TABLE_FAILED:                       "Failed to open mysql.event",
	ER_EVENT_NEITHER_M_EXPR_NOR_M_AT:                 "No datetime expression provided",
	ER_OBSOLETE_COL_COUNT_DOESNT_MATCH_CORRUPTED:     "Column count of mysql.%s is wrong. Expected %d, found %d. The table is probably corrupted",
	ER_OBSOLETE_CANNOT_LOAD_FROM_TABLE:               "Cannot load from mysql.%s. The table is probably corrupted",
	ER_EVENT_CANNOT_DELETE:                           "Failed to delete the event from mysql.event",
	ER_EVENT_COMPILE_ERROR:                           "Error during compilation of event's body",
	ER_EVENT_SAME_NAME:                               "Same old and new event name",
	ER_EVENT_DATA_TOO_LONG:                           "Data for column '%s' too long",
	ER_DROP_INDEX_FK:                                 "Cannot drop index '%-.192s': needed in a foreign key constraint",
	ER_WARN_DEPRECATED_SYNTAX_WITH_VER:               "The syntax '%s' is deprecated and will be removed in MySQL %s. Please use %s instead",
	ER_CANT_WRITE_LOCK_LOG_TABLE:                     "You can't write-lock a log table. Only read access is possible",
	ER_CANT_LOCK_LOG_TABLE:                           "You can't use locks with log tables.",
	ER_FOREIGN_DUPLICATE_KEY_OLD_UNUSED:              "Upholding foreign key constraints for table '%.192s', entry '%-.192s', key %d would lead to a duplicate entry",
	ER_COL_COUNT_DOESNT_MATCH_PLEASE_UPDATE:          "Column count of mysql.%s is wrong. Expected %d, found %d. Created with MySQL %d, now running %d. Please use mysql_upgrade to fix this error.",
	ER_TEMP_TABLE_PREVENTS_SWITCH_OUT_OF_RBR:         "Cannot switch out of the row-based binary log format when the session has open temporary tables",
	ER_STORED_FUNCTION_PREVENTS_SWITCH_BINLOG_FORMAT: "Cannot change the binary logging format inside a stored function or trigger",
	ER_NDB_CANT_SWITCH_BINLOG_FORMAT:                 "The NDB cluster engine does not support changing the binlog format on the fly yet",
	ER_PARTITION_NO_TEMPORARY:                        "Cannot create temporary table with partitions",
	ER_PARTITION_CONST_DOMAIN_ERROR:                  "Partition constant is out of partition function domain",
	ER_PARTITION_FUNCTION_IS_NOT_ALLOWED:             "This partition function is not allowed",
	ER_DDL_LOG_ERROR:                                 "Error in DDL log",
	ER_NULL_IN_VALUES_LESS_THAN:                      "Not allowed to use NULL value in VALUES LESS THAN",
	ER_WRONG_PARTITION_NAME:                          "Incorrect partition name",
	ER_CANT_CHANGE_TX_CHARACTERISTICS:                "Transaction characteristics can't be changed while a transaction is in progress",
	ER_DUP_ENTRY_AUTOINCREMENT_CASE:                  "ALTER TABLE causes auto_increment resequencing, resulting in duplicate entry '%-.192s' for key '%-.192s'",
	ER_EVENT_MODIFY_QUEUE_ERROR:                      "Internal scheduler error %d",
	ER_EVENT_SET_VAR_ERROR:                           "Error during starting/stopping of the scheduler. Error code %d",
	ER_PARTITION_MERGE_ERROR:                         "Engine cannot be used in partitioned tables",
	ER_CANT_ACTIVATE_LOG:                             "Cannot activate '%-.64s' log",
	ER_RBR_NOT_AVAILABLE:                             "The server was not built with row-based replication",
	ER_BASE64_DECODE_ERROR:                           "Decoding of base64 string failed",
	ER_EVENT_RECURSION_FORBIDDEN:                     "Recursion of EVENT DDL statements is forbidden when body is present",
	ER_EVENTS_DB_ERROR:                               "Cannot proceed because system tables used by Event Scheduler were found damaged at server start",
	ER_ONLY_INTEGERS_ALLOWED:                         "Only integers allowed as number here",
	ER_UNSUPORTED_LOG_ENGINE:                         "This storage engine cannot be used for log tables\"",
	ER_BAD_LOG_STATEMENT:                             "You cannot '%s' a log table if logging is enabled",
	ER_CANT_RENAME_LOG_TABLE:                         "Cannot rename '%s'. When logging enabled, rename to/from log table must rename two tables: the log table to an archive table and another table back to '%s'",
	ER_WRONG_PARAMCOUNT_TO_NATIVE_FCT:                "Incorrect parameter count in the call to native function '%-.192s'",
	ER_WRONG_PARAMETERS_TO_NATIVE_FCT:                "Incorrect parameters in the call to native function '%-.192s'",
	ER_WRONG_PARAMETERS_TO_STORED_FCT:                "Incorrect parameters in the call to stored function '%-.192s'",
	ER_NATIVE_FCT_NAME_COLLISION:                     "This function '%-.192s' has the same name as a native function",
	ER_DUP_ENTRY_WITH_KEY_NAME:                       "Duplicate entry '%-.64s' for key '%-.192s'",
	ER_BINLOG_PURGE_EMFILE:                           "Too many files opened, please execute the command again",
	ER_EVENT_CANNOT_CREATE_IN_THE_PAST:               "Event execution time is in the past and ON COMPLETION NOT PRESERVE is set. The event was dropped immediately after creation.",
	ER_EVENT_CANNOT_ALTER_IN_THE_PAST:                "Event execution time is in the past and ON COMPLETION NOT PRESERVE is set. The event was not changed. Specify a time in the future.",
	ER_SLAVE_INCIDENT:                                "The incident %s occurred on the master. Message: %-.64s",
	ER_NO_PARTITION_FOR_GIVEN_VALUE_SILENT:           "Table has no partition for some existing values",
	ER_BINLOG_UNSAFE_STATEMENT:                       "Unsafe statement written to the binary log using statement format since BINLOG_FORMAT = STATEMENT. %s",
	ER_SLAVE_FATAL_ERROR:                             "Fatal error: %s",
	ER_SLAVE_RELAY_LOG_READ_FAILURE:                  "Relay log read failure: %s",
	ER_SLAVE_RELAY_LOG_WRITE_FAILURE:                 "Relay log write failure: %s",
	ER_SLAVE_CREATE_EVENT_FAILURE:                    "Failed to create %s",
	ER_SLAVE_MASTER_COM_FAILURE:                      "Master command %s failed: %s",
	ER_BINLOG_LOGGING_IMPOSSIBLE:                     "Binary logging not possible. Message: %s",
	ER_VIEW_NO_CREATION_CTX:                          "View `%-.64s`.`%-.64s` has no creation context",
	ER_VIEW_INVALID_CREATION_CTX:                     "Creation context of view `%-.64s`.`%-.64s' is invalid",
	ER_SR_INVALID_CREATION_CTX:                       "Creation context of stored routine `%-.64s`.`%-.64s` is invalid",
	ER_TRG_CORRUPTED_FILE:                            "Corrupted TRG file for table `%-.64s`.`%-.64s`",
	ER_TRG_NO_CREATION_CTX:                           "Triggers for table `%-.64s`.`%-.64s` have no creation context",
	ER_TRG_INVALID_CREATION_CTX:                      "Trigger creation context of table `%-.64s`.`%-.64s` is invalid",
	ER_EVENT_INVALID_CREATION_CTX:                    "Creation context of event `%-.64s`.`%-.64s` is invalid",
	ER_TRG_CANT_OPEN_TABLE:                           "Cannot open table for trigger `%-.64s`.`%-.64s`",
	ER_CANT_CREATE_SROUTINE:                          "Cannot create stored routine `%-.64s`. Check warnings",
	ER_NEVER_USED:                                    "Ambiguous slave modes combination. %s",
	ER_NO_FORMAT_DESCRIPTION_EVENT_BEFORE_BINLOG_STATEMENT:              "The BINLOG statement of type `%s` was not preceded by a format description BINLOG statement.",
	ER_SLAVE_CORRUPT_EVENT:                                              "Corrupted replication event was detected",
	ER_LOAD_DATA_INVALID_COLUMN:                                         "Invalid column reference (%-.64s) in LOAD DATA",
	ER_LOG_PURGE_NO_FILE:                                                "Being purged log %s was not found",
	ER_XA_RBTIMEOUT:                                                     "XA_RBTIMEOUT: Transaction branch was rolled back: took too long",
	ER_XA_RBDEADLOCK:                                                    "XA_RBDEADLOCK: Transaction branch was rolled back: deadlock was detected",
	ER_NEED_REPREPARE:                                                   "Prepared statement needs to be re-prepared",
	ER_DELAYED_NOT_SUPPORTED:                                            "DELAYED option not supported for table '%-.192s'",
	WARN_NO_MASTER_INFO:                                                 "The master info structure does not exist",
	WARN_OPTION_IGNORED:                                                 "<%-.64s> option ignored",
	WARN_PLUGIN_DELETE_BUILTIN:                                          "Built-in plugins cannot be deleted",
	WARN_PLUGIN_BUSY:                                                    "Plugin is busy and will be uninstalled on shutdown",
	ER_VARIABLE_IS_READONLY:                                             "%s variable '%s' is read-only. Use SET %s to assign the value",
	ER_WARN_ENGINE_TRANSACTION_ROLLBACK:                                 "Storage engine %s does not support rollback for this statement. Transaction rolled back and must be restarted",
	ER_SLAVE_HEARTBEAT_FAILURE:                                          "Unexpected master's heartbeat data: %s",
	ER_SLAVE_HEARTBEAT_VALUE_OUT_OF_RANGE:                               "The requested value for the heartbeat period is either negative or exceeds the maximum allowed (%s seconds).",
	ER_NDB_REPLICATION_SCHEMA_ERROR:                                     "Bad schema for mysql.ndb_replication table. Message: %-.64s",
	ER_CONFLICT_FN_PARSE_ERROR:                                          "Error in parsing conflict function. Message: %-.64s",
	ER_EXCEPTIONS_WRITE_ERROR:                                           "Write to exceptions table failed. Message: %-.128s\"",
	ER_TOO_LONG_TABLE_COMMENT:                                           "Comment for table '%-.64s' is too long (max = %d)",
	ER_TOO_LONG_FIELD_COMMENT:                                           "Comment for field '%-.64s' is too long (max = %d)",
	ER_FUNC_INEXISTENT_NAME_COLLISION:                                   "FUNCTION %s does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual",
	ER_DATABASE_NAME:                                                    "Database",
	ER_TABLE_NAME:                                                       "Table",
	ER_PARTITION_NAME:                                                   "Partition",
	ER_SUBPARTITION_NAME:                                                "Subpartition",
	ER_TEMPORARY_NAME:                                                   "Temporary",
	ER_RENAMED_NAME:                                                     "Renamed",
	ER_TOO_MANY_CONCURRENT_TRXS:                                         "Too many active concurrent transactions",
	WARN_NON_ASCII_SEPARATOR_NOT_IMPLEMENTED:                            "Non-ASCII separator arguments are not fully supported",
	ER_DEBUG_SYNC_TIMEOUT:                                               "debug sync point wait timed out",
	ER_DEBUG_SYNC_HIT_LIMIT:                                             "debug sync point hit limit reached",
	ER_DUP_SIGNAL_SET:                                                   "Duplicate condition information item '%s'",
	ER_SIGNAL_WARN:                                                      "Unhandled user-defined warning condition",
	ER_SIGNAL_NOT_FOUND:                                                 "Unhandled user-defined not found condition",
	ER_SIGNAL_EXCEPTION:                                                 "Unhandled user-defined exception condition",
	ER_RESIGNAL_WITHOUT_ACTIVE_HANDLER:                                  "RESIGNAL when handler not active",
	ER_SIGNAL_BAD_CONDITION_TYPE:                                        "SIGNAL/RESIGNAL can only use a CONDITION defined with SQLSTATE",
	WARN_COND_ITEM_TRUNCATED:                                            "Data truncated for condition item '%s'",
	ER_COND_ITEM_TOO_LONG:                                               "Data too long for condition item '%s'",
	ER_UNKNOWN_LOCALE:                                                   "Unknown locale: '%-.64s'",
	ER_SLAVE_IGNORE_SERVER_IDS:                                          "The requested server id %d clashes with the slave startup option --replicate-same-server-id",
	ER_QUERY_CACHE_DISABLED:                                             "Query cache is disabled; restart the server with query_cache_type=1 to enable it",
	ER_SAME_NAME_PARTITION_FIELD:                                        "Duplicate partition field name '%-.192s'",
	ER_PARTITION_COLUMN_LIST_ERROR:                                      "Inconsistency in usage of column lists for partitioning",
	ER_WRONG_TYPE_COLUMN_VALUE_ERROR:                                    "Partition column values of incorrect type",
	ER_TOO_MANY_PARTITION_FUNC_FIELDS_ERROR:                             "Too many fields in '%-.192s'",
	ER_MAXVALUE_IN_VALUES_IN:                                            "Cannot use MAXVALUE as value in VALUES IN",
	ER_TOO_MANY_VALUES_ERROR:                                            "Cannot have more than one value for this type of %-.64s partitioning",
	ER_ROW_SINGLE_PARTITION_FIELD_ERROR:                                 "Row expressions in VALUES IN only allowed for multi-field column partitioning",
	ER_FIELD_TYPE_NOT_ALLOWED_AS_PARTITION_FIELD:                        "Field '%-.192s' is of a not allowed type for this type of partitioning",
	ER_PARTITION_FIELDS_TOO_LONG:                                        "The total length of the partitioning fields is too large",
	ER_BINLOG_ROW_ENGINE_AND_STMT_ENGINE:                                "Cannot execute statement: impossible to write to binary log since both row-incapable engines and statement-incapable engines are involved.",
	ER_BINLOG_ROW_MODE_AND_STMT_ENGINE:                                  "Cannot execute statement: impossible to write to binary log since BINLOG_FORMAT = ROW and at least one table uses a storage engine limited to statement-based logging.",
	ER_BINLOG_UNSAFE_AND_STMT_ENGINE:                                    "Cannot execute statement: impossible to write to binary log since statement is unsafe, storage engine is limited to statement-based logging, and BINLOG_FORMAT = MIXED. %s",
	ER_BINLOG_ROW_INJECTION_AND_STMT_ENGINE:                             "Cannot execute statement: impossible to write to binary log since statement is in row format and at least one table uses a storage engine limited to statement-based logging.",
	ER_BINLOG_STMT_MODE_AND_ROW_ENGINE:                                  "Cannot execute statement: impossible to write to binary log since BINLOG_FORMAT = STATEMENT and at least one table uses a storage engine limited to row-based logging.%s",
	ER_BINLOG_ROW_INJECTION_AND_STMT_MODE:                               "Cannot execute statement: impossible to write to binary log since statement is in row format and BINLOG_FORMAT = STATEMENT.",
	ER_BINLOG_MULTIPLE_ENGINES_AND_SELF_LOGGING_ENGINE:                  "Cannot execute statement: impossible to write to binary log since more than one engine is involved and at least one engine is self-logging.",
	ER_BINLOG_UNSAFE_LIMIT:                                              "The statement is unsafe because it uses a LIMIT clause. This is unsafe because the set of rows included cannot be predicted.",
	ER_BINLOG_UNSAFE_INSERT_DELAYED:                                     "The statement is unsafe because it uses INSERT DELAYED. This is unsafe because the times when rows are inserted cannot be predicted.",
	ER_BINLOG_UNSAFE_SYSTEM_TABLE:                                       "The statement is unsafe because it uses the general log, slow query log, or performance_schema table(s). This is unsafe because system tables may differ on slaves.",
	ER_BINLOG_UNSAFE_AUTOINC_COLUMNS:                                    "Statement is unsafe because it invokes a trigger or a stored function that inserts into an AUTO_INCREMENT column. Inserted values cannot be logged correctly.",
	ER_BINLOG_UNSAFE_UDF:                                                "Statement is unsafe because it uses a UDF which may not return the same value on the slave.",
	ER_BINLOG_UNSAFE_SYSTEM_VARIABLE:                                    "Statement is unsafe because it uses a system variable that may have a different value on the slave.",
	ER_BINLOG_UNSAFE_SYSTEM_FUNCTION:                                    "Statement is unsafe because it uses a system function that may return a different value on the slave.",
	ER_BINLOG_UNSAFE_NONTRANS_AFTER_TRANS:                               "Statement is unsafe because it accesses a non-transactional table after accessing a transactional table within the same transaction.",
	ER_MESSAGE_AND_STATEMENT:                                            "%s Statement: %s",
	ER_SLAVE_CONVERSION_FAILED:                                          "Column %d of table '%-.192s.%-.192s' cannot be converted from type '%-.32s' to type '%-.32s'",
	ER_SLAVE_CANT_CREATE_CONVERSION:                                     "Can't create conversion table for table '%-.192s.%-.192s'",
	ER_INSIDE_TRANSACTION_PREVENTS_SWITCH_BINLOG_FORMAT:                 "Cannot modify @@session.binlog_format inside a transaction",
	ER_PATH_LENGTH:                                                      "The path specified for %.64s is too long.",
	ER_WARN_DEPRECATED_SYNTAX_NO_REPLACEMENT:                            "'%s' is deprecated and will be removed in a future release.",
	ER_WRONG_NATIVE_TABLE_STRUCTURE:                                     "Native table '%-.64s'.'%-.64s' has the wrong structure",
	ER_WRONG_PERFSCHEMA_USAGE:                                           "Invalid performance_schema usage.",
	ER_WARN_I_S_SKIPPED_TABLE:                                           "Table '%s'.'%s' was skipped since its definition is being modified by concurrent DDL statement",
	ER_INSIDE_TRANSACTION_PREVENTS_SWITCH_BINLOG_DIRECT:                 "Cannot modify @@session.binlog_direct_non_transactional_updates inside a transaction",
	ER_STORED_FUNCTION_PREVENTS_SWITCH_BINLOG_DIRECT:                    "Cannot change the binlog direct flag inside a stored function or trigger",
	ER_SPATIAL_MUST_HAVE_GEOM_COL:                                       "A SPATIAL index may only contain a geometrical type column",
	ER_TOO_LONG_INDEX_COMMENT:                                           "Comment for index '%-.64s' is too long (max = %d)",
	ER_LOCK_ABORTED:                                                     "Wait on a lock was aborted due to a pending exclusive lock",
	ER_DATA_OUT_OF_RANGE:                                                "%s value is out of range in '%s'",
	ER_WRONG_SPVAR_TYPE_IN_LIMIT:                                        "A variable of a non-integer based type in LIMIT clause",
	ER_BINLOG_UNSAFE_MULTIPLE_ENGINES_AND_SELF_LOGGING_ENGINE:           "Mixing self-logging and non-self-logging engines in a statement is unsafe.",
	ER_BINLOG_UNSAFE_MIXED_STATEMENT:                                    "Statement accesses nontransactional table as well as transactional or temporary table, and writes to any of them.",
	ER_INSIDE_TRANSACTION_PREVENTS_SWITCH_SQL_LOG_BIN:                   "Cannot modify @@session.sql_log_bin inside a transaction",
	ER_STORED_FUNCTION_PREVENTS_SWITCH_SQL_LOG_BIN:                      "Cannot change the sql_log_bin inside a stored function or trigger",
	ER_FAILED_READ_FROM_PAR_FILE:                                        "Failed to read from the .par file",
	ER_VALUES_IS_NOT_INT_TYPE_ERROR:                                     "VALUES value for partition '%-.64s' must have type INT",
	ER_ACCESS_DENIED_NO_PASSWORD_ERROR:                                  "Access denied for user '%-.48s'@'%-.64s'",
	ER_SET_PASSWORD_AUTH_PLUGIN:                                         "SET PASSWORD has no significance for users authenticating via plugins",
	ER_GRANT_PLUGIN_USER_EXISTS:                                         "GRANT with IDENTIFIED WITH is illegal because the user %-.*s already exists",
	ER_TRUNCATE_ILLEGAL_FK:                                              "Cannot truncate a table referenced in a foreign key constraint (%.192s)",
	ER_PLUGIN_IS_PERMANENT:                                              "Plugin '%s' is force_plus_permanent and can not be unloaded",
	ER_SLAVE_HEARTBEAT_VALUE_OUT_OF_RANGE_MIN:                           "The requested value for the heartbeat period is less than 1 millisecond. The value is reset to 0, meaning that heartbeating will effectively be disabled.",
	ER_SLAVE_HEARTBEAT_VALUE_OUT_OF_RANGE_MAX:                           "The requested value for the heartbeat period exceeds the value of `slave_net_timeout' seconds. A sensible value for the period should be less than the timeout.",
	ER_STMT_CACHE_FULL:                                                  "Multi-row statements required more than 'max_binlog_stmt_cache_size' bytes of storage; increase this mysqld variable and try again",
	ER_MULTI_UPDATE_KEY_CONFLICT:                                        "Primary key/partition key update is not allowed since the table is updated both as '%-.192s' and '%-.192s'.",
	ER_TABLE_NEEDS_REBUILD:                                              "Table rebuild required. Please do \"ALTER TABLE `%-.32s` FORCE\" or dump/reload to fix it!",
	WARN_OPTION_BELOW_LIMIT:                                             "The value of '%s' should be no less than the value of '%s'",
	ER_INDEX_COLUMN_TOO_LONG:                                            "Index column size too large. The maximum column size is %d bytes.",
	ER_ERROR_IN_TRIGGER_BODY:                                            "Trigger '%-.64s' has an error in its body: '%-.256s'",
	ER_ERROR_IN_UNKNOWN_TRIGGER_BODY:                                    "Unknown trigger has an error in its body: '%-.256s'",
	ER_INDEX_CORRUPT:                                                    "Index %s is corrupted",
	ER_UNDO_RECORD_TOO_BIG:                                              "Undo log record is too big.",
	ER_BINLOG_UNSAFE_INSERT_IGNORE_SELECT:                               "INSERT IGNORE... SELECT is unsafe because the order in which rows are retrieved by the SELECT determines which (if any) rows are ignored. This order cannot be predicted and may differ on master and the slave.",
	ER_BINLOG_UNSAFE_INSERT_SELECT_UPDATE:                               "INSERT... SELECT... ON DUPLICATE KEY UPDATE is unsafe because the order in which rows are retrieved by the SELECT determines which (if any) rows are updated. This order cannot be predicted and may differ on master and the slave.",
	ER_BINLOG_UNSAFE_REPLACE_SELECT:                                     "REPLACE... SELECT is unsafe because the order in which rows are retrieved by the SELECT determines which (if any) rows are replaced. This order cannot be predicted and may differ on master and the slave.",
	ER_BINLOG_UNSAFE_CREATE_IGNORE_SELECT:                               "CREATE... IGNORE SELECT is unsafe because the order in which rows are retrieved by the SELECT determines which (if any) rows are ignored. This order cannot be predicted and may differ on master and the slave.",
	ER_BINLOG_UNSAFE_CREATE_REPLACE_SELECT:                              "CREATE... REPLACE SELECT is unsafe because the order in which rows are retrieved by the SELECT determines which (if any) rows are replaced. This order cannot be predicted and may differ on master and the slave.",
	ER_BINLOG_UNSAFE_UPDATE_IGNORE:                                      "UPDATE IGNORE is unsafe because the order in which rows are updated determines which (if any) rows are ignored. This order cannot be predicted and may differ on master and the slave.",
	ER_PLUGIN_NO_UNINSTALL:                                              "Plugin '%s' is marked as not dynamically uninstallable. You have to stop the server to uninstall it.",
	ER_PLUGIN_NO_INSTALL:                                                "Plugin '%s' is marked as not dynamically installable. You have to stop the server to install it.",
	ER_BINLOG_UNSAFE_WRITE_AUTOINC_SELECT:                               "Statements writing to a table with an auto-increment column after selecting from another table are unsafe because the order in which rows are retrieved determines what (if any) rows will be written. This order cannot be predicted and may differ on master and the slave.",
	ER_BINLOG_UNSAFE_CREATE_SELECT_AUTOINC:                              "CREATE TABLE... SELECT...  on a table with an auto-increment column is unsafe because the order in which rows are retrieved by the SELECT determines which (if any) rows are inserted. This order cannot be predicted and may differ on master and the slave.",
	ER_BINLOG_UNSAFE_INSERT_TWO_KEYS:                                    "INSERT... ON DUPLICATE KEY UPDATE  on a table with more than one UNIQUE KEY is unsafe",
	ER_TABLE_IN_FK_CHECK:                                                "Table is being used in foreign key check.",
	ER_UNSUPPORTED_ENGINE:                                               "Storage engine '%s' does not support system tables. [%s.%s]",
	ER_BINLOG_UNSAFE_AUTOINC_NOT_FIRST:                                  "INSERT into autoincrement field which is not the first part in the composed primary key is unsafe.",
	ER_CANNOT_LOAD_FROM_TABLE_V2:                                        "Cannot load from %s.%s. The table is probably corrupted",
	ER_MASTER_DELAY_VALUE_OUT_OF_RANGE:                                  "The requested value %d for the master delay exceeds the maximum %d",
	ER_ONLY_FD_AND_RBR_EVENTS_ALLOWED_IN_BINLOG_STATEMENT:               "Only Format_description_log_event and row events are allowed in BINLOG statements (but %s was provided)",
	ER_PARTITION_EXCHANGE_DIFFERENT_OPTION:                              "Non matching attribute '%-.64s' between partition and table",
	ER_PARTITION_EXCHANGE_PART_TABLE:                                    "Table to exchange with partition is partitioned: '%-.64s'",
	ER_PARTITION_EXCHANGE_TEMP_TABLE:                                    "Table to exchange with partition is temporary: '%-.64s'",
	ER_PARTITION_INSTEAD_OF_SUBPARTITION:                                "Subpartitioned table, use subpartition instead of partition",
	ER_UNKNOWN_PARTITION:                                                "Unknown partition '%-.64s' in table '%-.64s'",
	ER_TABLES_DIFFERENT_METADATA:                                        "Tables have different definitions",
	ER_ROW_DOES_NOT_MATCH_PARTITION:                                     "Found a row that does not match the partition",
	ER_BINLOG_CACHE_SIZE_GREATER_THAN_MAX:                               "Option binlog_cache_size (%d) is greater than max_binlog_cache_size (%d); setting binlog_cache_size equal to max_binlog_cache_size.",
	ER_WARN_INDEX_NOT_APPLICABLE:                                        "Cannot use %-.64s access on index '%-.64s' due to type or collation conversion on field '%-.64s'",
	ER_PARTITION_EXCHANGE_FOREIGN_KEY:                                   "Table to exchange with partition has foreign key references: '%-.64s'",
	ER_NO_SUCH_KEY_VALUE:                                                "Key value '%-.192s' was not found in table '%-.192s.%-.192s'",
	ER_RPL_INFO_DATA_TOO_LONG:                                           "Data for column '%s' too long",
	ER_NETWORK_READ_EVENT_CHECKSUM_FAILURE:                              "Replication event checksum verification failed while reading from network.",
	ER_BINLOG_READ_EVENT_CHECKSUM_FAILURE:                               "Replication event checksum verification failed while reading from a log file.",
	ER_BINLOG_STMT_CACHE_SIZE_GREATER_THAN_MAX:                          "Option binlog_stmt_cache_size (%d) is greater than max_binlog_stmt_cache_size (%d); setting binlog_stmt_cache_size equal to max_binlog_stmt_cache_size.",
	ER_CANT_UPDATE_TABLE_IN_CREATE_TABLE_SELECT:                         "Can't update table '%-.192s' while '%-.192s' is being created.",
	ER_PARTITION_CLAUSE_ON_NONPARTITIONED:                               "PARTITION () clause on non partitioned table",
	ER_ROW_DOES_NOT_MATCH_GIVEN_PARTITION_SET:                           "Found a row not matching the given partition set",
	ER_NO_SUCH_PARTITION__UNUSED:                                        "partition '%-.64s' doesn't exist",
	ER_CHANGE_RPL_INFO_REPOSITORY_FAILURE:                               "Failure while changing the type of replication repository: %s.",
	ER_WARNING_NOT_COMPLETE_ROLLBACK_WITH_CREATED_TEMP_TABLE:            "The creation of some temporary tables could not be rolled back.",
	ER_WARNING_NOT_COMPLETE_ROLLBACK_WITH_DROPPED_TEMP_TABLE:            "Some temporary tables were dropped, but these operations could not be rolled back.",
	ER_MTS_FEATURE_IS_NOT_SUPPORTED:                                     "%s is not supported in multi-threaded slave mode. %s",
	ER_MTS_UPDATED_DBS_GREATER_MAX:                                      "The number of modified databases exceeds the maximum %d; the database names will not be included in the replication event metadata.",
	ER_MTS_CANT_PARALLEL:                                                "Cannot execute the current event group in the parallel mode. Encountered event %s, relay-log name %s, position %s which prevents execution of this event group in parallel mode. Reason: %s.",
	ER_MTS_INCONSISTENT_DATA:                                            "%s",
	ER_FULLTEXT_NOT_SUPPORTED_WITH_PARTITIONING:                         "FULLTEXT index is not supported for partitioned tables.",
	ER_DA_INVALID_CONDITION_NUMBER:                                      "Invalid condition number",
	ER_INSECURE_PLAIN_TEXT:                                              "Sending passwords in plain text without SSL/TLS is extremely insecure.",
	ER_INSECURE_CHANGE_MASTER:                                           "Storing MySQL user name or password information in the master.info repository is not secure and is therefore not recommended. Please see the MySQL Manual for more about this issue and possible alternatives.",
	ER_FOREIGN_DUPLICATE_KEY_WITH_CHILD_INFO:                            "Foreign key constraint for table '%.192s', record '%-.192s' would lead to a duplicate entry in table '%.192s', key '%.192s'",
	ER_FOREIGN_DUPLICATE_KEY_WITHOUT_CHILD_INFO:                         "Foreign key constraint for table '%.192s', record '%-.192s' would lead to a duplicate entry in a child table",
	ER_SQLTHREAD_WITH_SECURE_SLAVE:                                      "Setting authentication options is not possible when only the Slave SQL Thread is being started.",
	ER_TABLE_HAS_NO_FT:                                                  "The table does not have FULLTEXT index to support this query",
	ER_VARIABLE_NOT_SETTABLE_IN_SF_OR_TRIGGER:                           "The system variable %.200s cannot be set in stored functions or triggers.",
	ER_VARIABLE_NOT_SETTABLE_IN_TRANSACTION:                             "The system variable %.200s cannot be set when there is an ongoing transaction.",
	ER_GTID_NEXT_IS_NOT_IN_GTID_NEXT_LIST:                               "The system variable @@SESSION.GTID_NEXT has the value %.200s, which is not listed in @@SESSION.GTID_NEXT_LIST.",
	ER_CANT_CHANGE_GTID_NEXT_IN_TRANSACTION_WHEN_GTID_NEXT_LIST_IS_NULL: "When @@SESSION.GTID_NEXT_LIST == NULL, the system variable @@SESSION.GTID_NEXT cannot change inside a transaction.",
	ER_SET_STATEMENT_CANNOT_INVOKE_FUNCTION:                             "The statement 'SET %.200s' cannot invoke a stored function.",
	ER_GTID_NEXT_CANT_BE_AUTOMATIC_IF_GTID_NEXT_LIST_IS_NON_NULL:        "The system variable @@SESSION.GTID_NEXT cannot be 'AUTOMATIC' when @@SESSION.GTID_NEXT_LIST is non-NULL.",
	ER_SKIPPING_LOGGED_TRANSACTION:                                      "Skipping transaction %.200s because it has already been executed and logged.",
	ER_MALFORMED_GTID_SET_SPECIFICATION:                                 "Malformed GTID set specification '%.200s'.",
	ER_MALFORMED_GTID_SET_ENCODING:                                      "Malformed GTID set encoding.",
	ER_MALFORMED_GTID_SPECIFICATION:                                     "Malformed GTID specification '%.200s'.",
	ER_GNO_EXHAUSTED:                                                    "Impossible to generate Global Transaction Identifier: the integer component reached the maximal value. Restart the server with a new server_uuid.",
	ER_BAD_SLAVE_AUTO_POSITION:                                          "Parameters MASTER_LOG_FILE, MASTER_LOG_POS, RELAY_LOG_FILE and RELAY_LOG_POS cannot be set when MASTER_AUTO_POSITION is active.",
	ER_AUTO_POSITION_REQUIRES_GTID_MODE_ON:                              "CHANGE MASTER TO MASTER_AUTO_POSITION = 1 can only be executed when @@GLOBAL.GTID_MODE = ON.",
	ER_CANT_DO_IMPLICIT_COMMIT_IN_TRX_WHEN_GTID_NEXT_IS_SET:             "Cannot execute statements with implicit commit inside a transaction when @@SESSION.GTID_NEXT != AUTOMATIC or @@SESSION.GTID_NEXT_LIST != NULL.",
	ER_GTID_MODE_2_OR_3_REQUIRES_ENFORCE_GTID_CONSISTENCY_ON:            "@@GLOBAL.GTID_MODE = ON or UPGRADE_STEP_2 requires @@GLOBAL.ENFORCE_GTID_CONSISTENCY = 1.",
	ER_GTID_MODE_REQUIRES_BINLOG:                                        "@@GLOBAL.GTID_MODE = ON or UPGRADE_STEP_1 or UPGRADE_STEP_2 requires --log-bin and --log-slave-updates.",
	ER_CANT_SET_GTID_NEXT_TO_GTID_WHEN_GTID_MODE_IS_OFF:                 "@@SESSION.GTID_NEXT cannot be set to UUID:NUMBER when @@GLOBAL.GTID_MODE = OFF.",
	ER_CANT_SET_GTID_NEXT_TO_ANONYMOUS_WHEN_GTID_MODE_IS_ON:             "@@SESSION.GTID_NEXT cannot be set to ANONYMOUS when @@GLOBAL.GTID_MODE = ON.",
	ER_CANT_SET_GTID_NEXT_LIST_TO_NON_NULL_WHEN_GTID_MODE_IS_OFF:        "@@SESSION.GTID_NEXT_LIST cannot be set to a non-NULL value when @@GLOBAL.GTID_MODE = OFF.",
	ER_FOUND_GTID_EVENT_WHEN_GTID_MODE_IS_OFF:                           "Found a Gtid_log_event or Previous_gtids_log_event when @@GLOBAL.GTID_MODE = OFF.",
	ER_GTID_UNSAFE_NON_TRANSACTIONAL_TABLE:                              "When @@GLOBAL.ENFORCE_GTID_CONSISTENCY = 1, updates to non-transactional tables can only be done in either autocommitted statements or single-statement transactions, and never in the same statement as updates to transactional tables.",
	ER_GTID_UNSAFE_CREATE_SELECT:                                        "CREATE TABLE ... SELECT is forbidden when @@GLOBAL.ENFORCE_GTID_CONSISTENCY = 1.",
	ER_GTID_UNSAFE_CREATE_DROP_TEMPORARY_TABLE_IN_TRANSACTION:           "When @@GLOBAL.ENFORCE_GTID_CONSISTENCY = 1, the statements CREATE TEMPORARY TABLE and DROP TEMPORARY TABLE can be executed in a non-transactional context only, and require that AUTOCOMMIT = 1.",
	ER_GTID_MODE_CAN_ONLY_CHANGE_ONE_STEP_AT_A_TIME:                     "The value of @@GLOBAL.GTID_MODE can only change one step at a time: OFF <-> UPGRADE_STEP_1 <-> UPGRADE_STEP_2 <-> ON. Also note that this value must be stepped up or down simultaneously on all servers; see the Manual for instructions.",
	ER_MASTER_HAS_PURGED_REQUIRED_GTIDS:                                 "The slave is connecting using CHANGE MASTER TO MASTER_AUTO_POSITION = 1, but the master has purged binary logs containing GTIDs that the slave requires.",
	ER_CANT_SET_GTID_NEXT_WHEN_OWNING_GTID:                              "@@SESSION.GTID_NEXT cannot be changed by a client that owns a GTID. The client owns %s. Ownership is released on COMMIT or ROLLBACK.",
	ER_UNKNOWN_EXPLAIN_FORMAT:                                           "Unknown EXPLAIN format name: '%s'",
	ER_CANT_EXECUTE_IN_READ_ONLY_TRANSACTION:                            "Cannot execute statement in a READ ONLY transaction.",
	ER_TOO_LONG_TABLE_PARTITION_COMMENT:                                 "Comment for table partition '%-.64s' is too long (max = %d)",
	ER_SLAVE_CONFIGURATION:                                              "Slave is not configured or failed to initialize properly. You must at least set --server-id to enable either a master or a slave. Additional error messages can be found in the MySQL error log.",
	ER_INNODB_FT_LIMIT:                                                  "InnoDB presently supports one FULLTEXT index creation at a time",
	ER_INNODB_NO_FT_TEMP_TABLE:                                          "Cannot create FULLTEXT index on temporary InnoDB table",
	ER_INNODB_FT_WRONG_DOCID_COLUMN:                                     "Column '%-.192s' is of wrong type for an InnoDB FULLTEXT index",
	ER_INNODB_FT_WRONG_DOCID_INDEX:                                      "Index '%-.192s' is of wrong type for an InnoDB FULLTEXT index",
	ER_INNODB_ONLINE_LOG_TOO_BIG:                                        "Creating index '%-.192s' required more than 'innodb_online_alter_log_max_size' bytes of modification log. Please try again.",
	ER_UNKNOWN_ALTER_ALGORITHM:                                          "Unknown ALGORITHM '%s'",
	ER_UNKNOWN_ALTER_LOCK:                                               "Unknown LOCK type '%s'",
	ER_MTS_CHANGE_MASTER_CANT_RUN_WITH_GAPS:                             "CHANGE MASTER cannot be executed when the slave was stopped with an error or killed in MTS mode. Consider using RESET SLAVE or START SLAVE UNTIL.",
	ER_MTS_RECOVERY_FAILURE:                                             "Cannot recover after SLAVE errored out in parallel execution mode. Additional error messages can be found in the MySQL error log.",
	ER_MTS_RESET_WORKERS:                                                "Cannot clean up worker info tables. Additional error messages can be found in the MySQL error log.",
	ER_COL_COUNT_DOESNT_MATCH_CORRUPTED_V2:                              "Column count of %s.%s is wrong. Expected %d, found %d. The table is probably corrupted",
	ER_SLAVE_SILENT_RETRY_TRANSACTION:                                   "Slave must silently retry current transaction",
	ER_DISCARD_FK_CHECKS_RUNNING:                                        "There is a foreign key check running on table '%-.192s'. Cannot discard the table.",
	ER_TABLE_SCHEMA_MISMATCH:                                            "Schema mismatch (%s)",
	ER_TABLE_IN_SYSTEM_TABLESPACE:                                       "Table '%-.192s' in system tablespace",
	ER_IO_READ_ERROR:                                                    "IO Read error: (%d, %s) %s",
	ER_IO_WRITE_ERROR:                                                   "IO Write error: (%d, %s) %s",
	ER_TABLESPACE_MISSING:                                               "Tablespace is missing for table '%-.192s'",
	ER_TABLESPACE_EXISTS:                                                "Tablespace for table '%-.192s' exists. Please DISCARD the tablespace before IMPORT.",
	ER_TABLESPACE_DISCARDED:                                             "Tablespace has been discarded for table '%-.192s'",
	ER_INTERNAL_ERROR:                                                   "Internal error: %s",
	ER_INNODB_IMPORT_ERROR:                                              "ALTER TABLE '%-.192s' IMPORT TABLESPACE failed with error %d : '%s'",
	ER_INNODB_INDEX_CORRUPT:                                             "Index corrupt: %s",
	ER_INVALID_YEAR_COLUMN_LENGTH:                                       "YEAR(%d) column type is deprecated. Creating YEAR(4) column instead.",
	ER_NOT_VALID_PASSWORD:                                               "Your password does not satisfy the current policy requirements",
	ER_MUST_CHANGE_PASSWORD:                                             "You must SET PASSWORD before executing this statement",
	ER_FK_NO_INDEX_CHILD:                                                "Failed to add the foreign key constaint. Missing index for constraint '%s' in the foreign table '%s'",
	ER_FK_NO_INDEX_PARENT:                                               "Failed to add the foreign key constaint. Missing index for constraint '%s' in the referenced table '%s'",
	ER_FK_FAIL_ADD_SYSTEM:                                               "Failed to add the foreign key constraint '%s' to system tables",
	ER_FK_CANNOT_OPEN_PARENT:                                            "Failed to open the referenced table '%s'",
	ER_FK_INCORRECT_OPTION:                                              "Failed to add the foreign key constraint on table '%s'. Incorrect options in FOREIGN KEY constraint '%s'",
	ER_FK_DUP_NAME:                                                      "Duplicate foreign key constraint name '%s'",
	ER_PASSWORD_FORMAT:                                                  "The password hash doesn't have the expected format. Check if the correct password algorithm is being used with the PASSWORD() function.",
	ER_FK_COLUMN_CANNOT_DROP:                                            "Cannot drop column '%-.192s': needed in a foreign key constraint '%-.192s'",
	ER_FK_COLUMN_CANNOT_DROP_CHILD:                                      "Cannot drop column '%-.192s': needed in a foreign key constraint '%-.192s' of table '%-.192s'",
	ER_FK_COLUMN_NOT_NULL:                                               "Column '%-.192s' cannot be NOT NULL: needed in a foreign key constraint '%-.192s' SET NULL",
	ER_DUP_INDEX:                                                        "Duplicate index '%-.64s' defined on the table '%-.64s.%-.64s'. This is deprecated and will be disallowed in a future release.",
	ER_FK_COLUMN_CANNOT_CHANGE:                                          "Cannot change column '%-.192s': used in a foreign key constraint '%-.192s'",
	ER_FK_COLUMN_CANNOT_CHANGE_CHILD:                                    "Cannot change column '%-.192s': used in a foreign key constraint '%-.192s' of table '%-.192s'",
	ER_FK_CANNOT_DELETE_PARENT:                                          "Cannot delete rows from table which is parent in a foreign key constraint '%-.192s' of table '%-.192s'",
	ER_MALFORMED_PACKET:                                                 "Malformed communication packet.",
	ER_READ_ONLY_MODE:                                                   "Running in read-only mode",
	ER_GTID_NEXT_TYPE_UNDEFINED_GROUP:                                   "When @@SESSION.GTID_NEXT is set to a GTID, you must explicitly set it again after a COMMIT or ROLLBACK. If you see this error message in the slave SQL thread, it means that a table in the current transaction is transactional on the master and non-transactional on the slave. In a client connection, it means that you executed SET @@SESSION.GTID_NEXT before a transaction and forgot to set @@SESSION.GTID_NEXT to a different identifier or to 'AUTOMATIC' after COMMIT or ROLLBACK. Current @@SESSION.GTID_NEXT is '%s'.",
	ER_VARIABLE_NOT_SETTABLE_IN_SP:                                      "The system variable %.200s cannot be set in stored procedures.",
	ER_CANT_SET_GTID_PURGED_WHEN_GTID_MODE_IS_OFF:                       "@@GLOBAL.GTID_PURGED can only be set when @@GLOBAL.GTID_MODE = ON.",
	ER_CANT_SET_GTID_PURGED_WHEN_GTID_EXECUTED_IS_NOT_EMPTY:             "@@GLOBAL.GTID_PURGED can only be set when @@GLOBAL.GTID_EXECUTED is empty.",
	ER_CANT_SET_GTID_PURGED_WHEN_OWNED_GTIDS_IS_NOT_EMPTY:               "@@GLOBAL.GTID_PURGED can only be set when there are no ongoing transactions (not even in other clients).",
	ER_GTID_PURGED_WAS_CHANGED:                                          "@@GLOBAL.GTID_PURGED was changed from '%s' to '%s'.",
	ER_GTID_EXECUTED_WAS_CHANGED:                                        "@@GLOBAL.GTID_EXECUTED was changed from '%s' to '%s'.",
	ER_BINLOG_STMT_MODE_AND_NO_REPL_TABLES:                              "Cannot execute statement: impossible to write to binary log since BINLOG_FORMAT = STATEMENT, and both replicated and non replicated tables are written to.",
	ER_ALTER_OPERATION_NOT_SUPPORTED:                                    "%s is not supported for this operation. Try %s.",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON:                             "%s is not supported. Reason: %s. Try %s.",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_COPY:                        "COPY algorithm requires a lock",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_PARTITION:                   "Partition specific operations do not yet support LOCK/ALGORITHM",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_FK_RENAME:                   "Columns participating in a foreign key are renamed",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_COLUMN_TYPE:                 "Cannot change column type INPLACE",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_FK_CHECK:                    "Adding foreign keys needs foreign_key_checks=OFF",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_IGNORE:                      "Creating unique indexes with IGNORE requires COPY algorithm to remove duplicate rows",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_NOPK:                        "Dropping a primary key is not allowed without also adding a new primary key",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_AUTOINC:                     "Adding an auto-increment column requires a lock",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_HIDDEN_FTS:                  "Cannot replace hidden FTS_DOC_ID with a user-visible one",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_CHANGE_FTS:                  "Cannot drop or rename FTS_DOC_ID",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_FTS:                         "Fulltext index creation requires a lock",
	ER_SQL_SLAVE_SKIP_COUNTER_NOT_SETTABLE_IN_GTID_MODE:                 "sql_slave_skip_counter can not be set when the server is running with @@GLOBAL.GTID_MODE = ON. Instead, for each transaction that you want to skip, generate an empty transaction with the same GTID as the transaction",
	ER_DUP_UNKNOWN_IN_INDEX:                                             "Duplicate entry for key '%-.192s'",
	ER_IDENT_CAUSES_TOO_LONG_PATH:                                       "Long database name and identifier for object resulted in path length exceeding %d characters. Path: '%s'.",
	ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_NOT_NULL:                    "cannot silently convert NULL values, as required in this SQL_MODE",
	ER_MUST_CHANGE_PASSWORD_LOGIN:                                       "Your password has expired. To log in you must change it using a client that supports expired passwords.",
	ER_ROW_IN_WRONG_PARTITION:                                           "Found a row in wrong partition %s",
}

const (
	CLIENT_LONG_PASSWORD uint32 = 1 << iota
	CLIENT_FOUND_ROWS
	CLIENT_LONG_FLAG
	CLIENT_CONNECT_WITH_DB
	CLIENT_NO_SCHEMA
	CLIENT_COMPRESS
	CLIENT_ODBC
	CLIENT_LOCAL_FILES
	CLIENT_IGNORE_SPACE
	CLIENT_PROTOCOL_41
	CLIENT_INTERACTIVE
	CLIENT_SSL
	CLIENT_IGNORE_SIGPIPE
	CLIENT_TRANSACTIONS
	CLIENT_RESERVED
	CLIENT_SECURE_CONNECTION
	CLIENT_MULTI_STATEMENTS
	CLIENT_MULTI_RESULTS
	CLIENT_PS_MULTI_RESULTS
	CLIENT_PLUGIN_AUTH
	CLIENT_CONNECT_ATTRS
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA
)

const (
	DEFAULT_CHARSET               = "utf8"
	DEFAULT_COLLATION_ID   uint8  = 33
	DEFAULT_COLLATION_NAME string = "utf8_general_ci"
)

const (
	MinProtocolVersion byte   = 10
	MaxPayloadLen      int    = 1<<24 - 1
	TimeFormat         string = "2006-01-02 15:04:05"
)

const (
	AUTH_MYSQL_OLD_PASSWORD    = "mysql_old_password"
	AUTH_NATIVE_PASSWORD       = "mysql_native_password"
	AUTH_CLEAR_PASSWORD        = "mysql_clear_password"
	AUTH_CACHING_SHA2_PASSWORD = "caching_sha2_password"
	AUTH_SHA256_PASSWORD       = "sha256_password"
)

const (
	NOT_NULL_FLAG       = 1
	PRI_KEY_FLAG        = 2
	UNIQUE_KEY_FLAG     = 4
	BLOB_FLAG           = 16
	UNSIGNED_FLAG       = 32
	ZEROFILL_FLAG       = 64
	BINARY_FLAG         = 128
	ENUM_FLAG           = 256
	AUTO_INCREMENT_FLAG = 512
	TIMESTAMP_FLAG      = 1024
	SET_FLAG            = 2048
	NUM_FLAG            = 32768
	PART_KEY_FLAG       = 16384
	GROUP_FLAG          = 32768
	UNIQUE_FLAG         = 65536
)

const (
	MYSQL_TYPE_DECIMAL byte = iota
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT

	//mysql 5.6
	MYSQL_TYPE_TIMESTAMP2
	MYSQL_TYPE_DATETIME2
	MYSQL_TYPE_TIME2
)

const (
	MYSQL_TYPE_JSON byte = iota + 0xf5
	MYSQL_TYPE_NEWDECIMAL
	MYSQL_TYPE_ENUM
	MYSQL_TYPE_SET
	MYSQL_TYPE_TINY_BLOB
	MYSQL_TYPE_MEDIUM_BLOB
	MYSQL_TYPE_LONG_BLOB
	MYSQL_TYPE_BLOB
	MYSQL_TYPE_VAR_STRING
	MYSQL_TYPE_STRING
	MYSQL_TYPE_GEOMETRY
)

type EventType byte

const (
	UNKNOWN_EVENT EventType = iota
	START_EVENT_V3
	QUERY_EVENT
	STOP_EVENT
	ROTATE_EVENT
	INTVAR_EVENT
	LOAD_EVENT
	SLAVE_EVENT
	CREATE_FILE_EVENT
	APPEND_BLOCK_EVENT
	EXEC_LOAD_EVENT
	DELETE_FILE_EVENT
	NEW_LOAD_EVENT
	RAND_EVENT
	USER_VAR_EVENT
	FORMAT_DESCRIPTION_EVENT
	XID_EVENT
	BEGIN_LOAD_QUERY_EVENT
	EXECUTE_LOAD_QUERY_EVENT
	TABLE_MAP_EVENT
	WRITE_ROWS_EVENTv0
	UPDATE_ROWS_EVENTv0
	DELETE_ROWS_EVENTv0
	WRITE_ROWS_EVENTv1
	UPDATE_ROWS_EVENTv1
	DELETE_ROWS_EVENTv1
	INCIDENT_EVENT
	HEARTBEAT_EVENT
	IGNORABLE_EVENT
	ROWS_QUERY_EVENT
	WRITE_ROWS_EVENTv2
	UPDATE_ROWS_EVENTv2
	DELETE_ROWS_EVENTv2
	GTID_EVENT
	ANONYMOUS_GTID_EVENT
	PREVIOUS_GTIDS_EVENT
	TRANSACTION_CONTEXT_EVENT
	VIEW_CHANGE_EVENT
	XA_PREPARE_LOG_EVENT
)

func (e EventType) String() string {
	switch e {
	case UNKNOWN_EVENT:
		return "UnknownEvent"
	case START_EVENT_V3:
		return "StartEventV3"
	case QUERY_EVENT:
		return "QueryEvent"
	case STOP_EVENT:
		return "StopEvent"
	case ROTATE_EVENT:
		return "RotateEvent"
	case INTVAR_EVENT:
		return "IntVarEvent"
	case LOAD_EVENT:
		return "LoadEvent"
	case SLAVE_EVENT:
		return "SlaveEvent"
	case CREATE_FILE_EVENT:
		return "CreateFileEvent"
	case APPEND_BLOCK_EVENT:
		return "AppendBlockEvent"
	case EXEC_LOAD_EVENT:
		return "ExecLoadEvent"
	case DELETE_FILE_EVENT:
		return "DeleteFileEvent"
	case NEW_LOAD_EVENT:
		return "NewLoadEvent"
	case RAND_EVENT:
		return "RandEvent"
	case USER_VAR_EVENT:
		return "UserVarEvent"
	case FORMAT_DESCRIPTION_EVENT:
		return "FormatDescriptionEvent"
	case XID_EVENT:
		return "XIDEvent"
	case BEGIN_LOAD_QUERY_EVENT:
		return "BeginLoadQueryEvent"
	case EXECUTE_LOAD_QUERY_EVENT:
		return "ExectueLoadQueryEvent"
	case TABLE_MAP_EVENT:
		return "TableMapEvent"
	case WRITE_ROWS_EVENTv0:
		return "WriteRowsEventV0"
	case UPDATE_ROWS_EVENTv0:
		return "UpdateRowsEventV0"
	case DELETE_ROWS_EVENTv0:
		return "DeleteRowsEventV0"
	case WRITE_ROWS_EVENTv1:
		return "WriteRowsEventV1"
	case UPDATE_ROWS_EVENTv1:
		return "UpdateRowsEventV1"
	case DELETE_ROWS_EVENTv1:
		return "DeleteRowsEventV1"
	case INCIDENT_EVENT:
		return "IncidentEvent"
	case HEARTBEAT_EVENT:
		return "HeartbeatEvent"
	case IGNORABLE_EVENT:
		return "IgnorableEvent"
	case ROWS_QUERY_EVENT:
		return "RowsQueryEvent"
	case WRITE_ROWS_EVENTv2:
		return "WriteRowsEventV2"
	case UPDATE_ROWS_EVENTv2:
		return "UpdateRowsEventV2"
	case DELETE_ROWS_EVENTv2:
		return "DeleteRowsEventV2"
	case GTID_EVENT:
		return "GTIDEvent"
	case ANONYMOUS_GTID_EVENT:
		return "AnonymousGTIDEvent"
	case PREVIOUS_GTIDS_EVENT:
		return "PreviousGTIDsEvent"
	case MARIADB_ANNOTATE_ROWS_EVENT:
		return "MariadbAnnotateRowsEvent"
	case MARIADB_BINLOG_CHECKPOINT_EVENT:
		return "MariadbBinLogCheckPointEvent"
	case MARIADB_GTID_EVENT:
		return "MariadbGTIDEvent"
	case MARIADB_GTID_LIST_EVENT:
		return "MariadbGTIDListEvent"
	case TRANSACTION_CONTEXT_EVENT:
		return "TransactionContextEvent"
	case VIEW_CHANGE_EVENT:
		return "ViewChangeEvent"
	case XA_PREPARE_LOG_EVENT:
		return "XAPrepareLogEvent"

	default:
		return "UnknownEvent"
	}
}

const (
	// MariaDB event starts from 160
	MARIADB_ANNOTATE_ROWS_EVENT EventType = 160 + iota
	MARIADB_BINLOG_CHECKPOINT_EVENT
	MARIADB_GTID_EVENT
	MARIADB_GTID_LIST_EVENT
)

const (
	EventHeaderSize            = 19
	SidLength                  = 16
	LogicalTimestampTypeCode   = 2
	PartLogicalTimestampLength = 8
	BinlogChecksumLength       = 4
	UndefinedServerVer         = 999999 // UNDEFINED_SERVER_VERSION
)
