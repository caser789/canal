package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pingcap/errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrSkip = errors.New("Handler error, but skipped")
)

var binlogExp *regexp.Regexp
var useExp *regexp.Regexp
var valuesExp *regexp.Regexp
var gtidExp *regexp.Regexp

func init() {
	binlogExp = regexp.MustCompile(`^CHANGE MASTER TO MASTER_LOG_FILE='(.+)', MASTER_LOG_POS=(\d+);`)
	useExp = regexp.MustCompile("^USE `(.+)`;")
	valuesExp = regexp.MustCompile("^INSERT INTO `(.+?)` VALUES \\((.+)\\);$")
	// The pattern will only match MySQL GTID, as you know SET GLOBAL gtid_slave_pos='0-1-4' is used for MariaDB.
	// SET @@GLOBAL.GTID_PURGED='1638041a-0457-11e9-bb9f-00505690b730:1-429405150';
	// https://dev.mysql.com/doc/refman/5.7/en/replication-gtids-concepts.html
	gtidExp = regexp.MustCompile(`(\w{8}(-\w{4}){3}-\w{12}(:\d+(-\d+)?)+)`)
}

type Dumper struct {
	// mysqldump execution path, like mysqldump or /usr/bin/mysqldump, etc...
	ExecutionPath string

	Addr     string
	User     string
	Password string
	Protocol string

	// Will override Databases
	Tables  []string
	TableDB string

	Databases []string

	Where   string
	Charset string

	IgnoreTables map[string][]string

	ExtraOptions []string

	ErrOut io.Writer

	masterDataSkipped bool
	maxAllowedPacket  int
	hexBlob           bool

	// see detectColumnStatisticsParamSupported
	isColumnStatisticsParamSupported bool
}

type ParseHandler interface {
	// Parse CHANGE MASTER TO MASTER_LOG_FILE=name, MASTER_LOG_POS=pos;
	BinLog(name string, pos uint64) error
	GtidSet(gtidsets string) error
	Data(schema string, table string, values []string) error
}

// Parse the dump data with Dumper generate.
// It can not parse all the data formats with mysqldump outputs
func Parse(r io.Reader, h ParseHandler, parseBinlogPos bool) error {
	rb := bufio.NewReaderSize(r, 1024*16)

	var db string
	var binlogParsed bool

	for {
		line, err := rb.ReadString('\n')
		if err != nil && err != io.EOF {
			return errors.Trace(err)
		}

		if ErrorEqual(err, io.EOF) {
			break
		}

		// Ignore '\n' on Linux or '\r\n' on Windows
		line = strings.TrimRightFunc(line, func(c rune) bool {
			return c == '\r' || c == '\n'
		})

		if parseBinlogPos && !binlogParsed {
			// [1] parse gtid set from mysqldump
			// gtid comes before binlog file-position
			if m := gtidExp.FindAllStringSubmatch(line, -1); len(m) == 1 {
				gtidStr := m[0][1]
				if gtidStr != "" {
					if err := h.GtidSet(gtidStr); err != nil {
						return errors.Trace(err)
					}
				}
			}

			// [2] parse binlog
			//     master log file name
			//     master log file position
			if m := binlogExp.FindAllStringSubmatch(line, -1); len(m) == 1 {
				name := m[0][1]
				pos, err := strconv.ParseUint(m[0][2], 10, 64)
				if err != nil {
					return errors.Errorf("parse binlog %v err, invalid number", line)
				}

				if err = h.BinLog(name, pos); err != nil && err != ErrSkip {
					return errors.Trace(err)
				}

				binlogParsed = true
			}
		}

		// [3] parse USE
		//     database
		if m := useExp.FindAllStringSubmatch(line, -1); len(m) == 1 {
			db = m[0][1]
		}

		// [4] parse values
		//     table
		//     values (insert into values)
		if m := valuesExp.FindAllStringSubmatch(line, -1); len(m) == 1 {
			table := m[0][1]

			values, err := parseValues(m[0][2])
			if err != nil {
				return errors.Errorf("parse values %v err", line)
			}

			if err = h.Data(db, table, values); err != nil && err != ErrSkip {
				return errors.Trace(err)
			}
		}
	}

	return nil
}

func parseValues(str string) ([]string, error) {
	// values are separated by comma, but we can not split using comma directly
	// string is enclosed by single quote

	// a simple implementation, may be more robust later.

	values := make([]string, 0, 8)

	i := 0
	for i < len(str) {
		if str[i] != '\'' {
			// no string, read until comma
			j := i + 1
			for ; j < len(str) && str[j] != ','; j++ {
			}
			values = append(values, str[i:j])
			// skip ,
			i = j + 1
		} else {
			// read string until another single quote
			j := i + 1

			escaped := false
			for j < len(str) {
				if str[j] == '\\' {
					// skip escaped character
					j += 2
					escaped = true
					continue
				} else if str[j] == '\'' {
					break
				} else {
					j++
				}
			}

			if j >= len(str) {
				return nil, fmt.Errorf("parse quote values error")
			}

			value := str[i : j+1]
			if escaped {
				value = unescapeString(value)
			}
			values = append(values, value)
			// skip ' and ,
			i = j + 2
		}

		// need skip blank???
	}

	return values, nil
}

// unescapeString un-escapes the string.
// mysqldump will escape the string when dumps,
// Refer http://dev.mysql.com/doc/refman/5.7/en/string-literals.html
func unescapeString(s string) string {
	i := 0

	value := make([]byte, 0, len(s))
	for i < len(s) {
		if s[i] == '\\' {
			j := i + 1
			if j == len(s) {
				// The last char is \, remove
				break
			}

			value = append(value, unescapeChar(s[j]))
			i += 2
		} else {
			value = append(value, s[i])
			i++
		}
	}

	return string(value)
}

func unescapeChar(ch byte) byte {
	// \" \' \\ \n \0 \b \Z \r \t ==> escape to one char
	switch ch {
	case 'n':
		ch = '\n'
	case '0':
		ch = 0
	case 'b':
		ch = 8
	case 'Z':
		ch = 26
	case 'r':
		ch = '\r'
	case 't':
		ch = '\t'
	}
	return ch
}

func (d *Dumper) Dump(w io.Writer) error {
	args := make([]string, 0, 16)

	// Common args
	if strings.Contains(d.Addr, "/") {
		args = append(args, fmt.Sprintf("--socket=%s", d.Addr))
	} else {
		seps := strings.SplitN(d.Addr, ":", 2)
		args = append(args, fmt.Sprintf("--host=%s", seps[0]))
		if len(seps) > 1 {
			args = append(args, fmt.Sprintf("--port=%s", seps[1]))
		}
	}

	args = append(args, fmt.Sprintf("--user=%s", d.User))
	passwordArg := fmt.Sprintf("--password=%s", d.Password)
	args = append(args, passwordArg)
	passwordArgIndex := len(args) - 1

	if !d.masterDataSkipped {
		args = append(args, "--master-data")
	}

	if d.maxAllowedPacket > 0 {
		// mysqldump param should be --max-allowed-packet=%dM not be --max_allowed_packet=%dM
		args = append(args, fmt.Sprintf("--max-allowed-packet=%dM", d.maxAllowedPacket))
	}

	if d.Protocol != "" {
		args = append(args, fmt.Sprintf("--protocol=%s", d.Protocol))
	}

	args = append(args, "--single-transaction")
	args = append(args, "--skip-lock-tables")

	// Disable uncessary data
	args = append(args, "--compact")
	args = append(args, "--skip-opt")
	args = append(args, "--quick")

	// We only care about data
	args = append(args, "--no-create-info")

	// Multi row is easy for us to parse the data
	args = append(args, "--skip-extended-insert")
	args = append(args, "--skip-tz-utc")
	if d.hexBlob {
		// Use hex for the binary type
		args = append(args, "--hex-blob")
	}

	for db, tables := range d.IgnoreTables {
		for _, table := range tables {
			args = append(args, fmt.Sprintf("--ignore-table=%s.%s", db, table))
		}
	}

	if len(d.Charset) != 0 {
		args = append(args, fmt.Sprintf("--default-character-set=%s", d.Charset))
	}

	if len(d.Where) != 0 {
		args = append(args, fmt.Sprintf("--where=%s", d.Where))
	}

	if len(d.ExtraOptions) != 0 {
		args = append(args, d.ExtraOptions...)
	}

	if d.isColumnStatisticsParamSupported {
		args = append(args, `--column-statistics=0`)
	}

	if len(d.Tables) == 0 && len(d.Databases) == 0 {
		args = append(args, "--all-databases")
	} else if len(d.Tables) == 0 {
		args = append(args, "--databases")
		args = append(args, d.Databases...)
	} else {
		args = append(args, d.TableDB)
		args = append(args, d.Tables...)

		// If we only dump some tables, the dump data will not have database name
		// which makes us hard to parse, so here we add it manually.

		_, err := w.Write([]byte(fmt.Sprintf("USE `%s`;\n", d.TableDB)))
		if err != nil {
			return fmt.Errorf(`could not write USE command: %w`, err)
		}
	}

	args[passwordArgIndex] = "--password=******"
	fmt.Printf("exec mysqldump with %v\n", args)
	args[passwordArgIndex] = passwordArg
	cmd := exec.Command(d.ExecutionPath, args...)

	cmd.Stderr = d.ErrOut
	cmd.Stdout = w

	return cmd.Run()
}

func NewDumper(executionPath string, addr string, user string, password string) (*Dumper, error) {
	if len(executionPath) == 0 {
		return nil, nil
	}

	path, err := exec.LookPath(executionPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d := new(Dumper)
	d.ExecutionPath = path
	d.Addr = addr
	d.User = user
	d.Password = password
	d.Tables = make([]string, 0, 16)
	d.Databases = make([]string, 0, 16)
	d.Charset = DEFAULT_CHARSET
	d.IgnoreTables = make(map[string][]string)
	d.ExtraOptions = make([]string, 0, 5)
	d.masterDataSkipped = false
	d.isColumnStatisticsParamSupported = d.detectColumnStatisticsParamSupported()

	d.ErrOut = os.Stderr

	return d, nil
}

// New mysqldump versions try to send queries to information_schema.COLUMN_STATISTICS table which does not exist in old MySQL (<5.x).
// And we got error: "Unknown table 'COLUMN_STATISTICS' in information_schema (1109)".
//
// mysqldump may not send this query if it is started with parameter --column-statistics.
// But this parameter exists only for versions >=8.0.2 (https://dev.mysql.com/doc/relnotes/mysql/8.0/en/news-8-0-2.html).
//
// For environments where the version of mysql-server and mysqldump differs, we try to check this parameter and use it if available.
func (d *Dumper) detectColumnStatisticsParamSupported() bool {
	out, err := exec.Command(d.ExecutionPath, `--help`).CombinedOutput()
	if err != nil {
		return false
	}
	return bytes.Contains(out, []byte(`--column-statistics`))
}

func (d *Dumper) SetCharset(charset string) {
	d.Charset = charset
}

func (d *Dumper) SetProtocol(protocol string) {
	d.Protocol = protocol
}

func (d *Dumper) SetWhere(where string) {
	d.Where = where
}

func (d *Dumper) SetExtraOptions(options []string) {
	d.ExtraOptions = options
}

func (d *Dumper) SetErrOut(o io.Writer) {
	d.ErrOut = o
}

// SkipMasterData: In some cloud MySQL, we have no privilege to use `--master-data`.
func (d *Dumper) SkipMasterData(v bool) {
	d.masterDataSkipped = v
}

func (d *Dumper) SetMaxAllowedPacket(i int) {
	d.maxAllowedPacket = i
}

func (d *Dumper) SetHexBlob(v bool) {
	d.hexBlob = v
}

func (d *Dumper) AddDatabases(dbs ...string) {
	d.Databases = append(d.Databases, dbs...)
}

func (d *Dumper) AddTables(db string, tables ...string) {
	if d.TableDB != db {
		d.TableDB = db
		d.Tables = d.Tables[0:0]
	}

	d.Tables = append(d.Tables, tables...)
}

func (d *Dumper) AddIgnoreTables(db string, tables ...string) {
	t := d.IgnoreTables[db]
	t = append(t, tables...)
	d.IgnoreTables[db] = t
}

func (d *Dumper) Reset() {
	d.Tables = d.Tables[0:0]
	d.TableDB = ""
	d.IgnoreTables = make(map[string][]string)
	d.Databases = d.Databases[0:0]
	d.Where = ""
}

// DumpAndParse: Dump MySQL and parse immediately
func (d *Dumper) DumpAndParse(h ParseHandler) error {
	r, w := io.Pipe()

	done := make(chan error, 1)
	go func() {
		err := Parse(r, h, !d.masterDataSkipped)
		_ = r.CloseWithError(err)
		done <- err
	}()

	err := d.Dump(w)
	_ = w.CloseWithError(err)

	err = <-done

	return errors.Trace(err)
}
