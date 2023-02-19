package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/pingcap/errors"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

var (
	errSyncRunning         = errors.New("Sync is running, must Close first")
	SemiSyncIndicator byte = 0xef
)

// Like vitess, use flavor for different MySQL versions,
const (
	MySQLFlavor   = "mysql"
	MariaDBFlavor = "mariadb"
)

func NewBinlogSyncer(cfg BinlogSyncerConfig) *BinlogSyncer {
	if cfg.ServerID == 0 {
		fmt.Println("can't use 0 as the server ID")
		return nil
	}

	if cfg.Dialer == nil {
		dialer := &net.Dialer{}
		cfg.Dialer = dialer.DialContext
	}

	b := new(BinlogSyncer)
	b.cfg = cfg
	b.parser = NewBinlogParser()
	b.parser.SetFlavor(cfg.Flavor)
	b.parser.SetRawMode(b.cfg.RawModeEnabled)
	b.parser.SetParseTime(b.cfg.ParseTime)
	b.parser.SetTimestampStringLocation(b.cfg.TimestampStringLocation)
	b.parser.SetUseDecimal(b.cfg.UseDecimal)
	b.parser.SetVerifyChecksum(b.cfg.VerifyChecksum)
	b.parser.SetRowsEventDecodeFunc(b.cfg.RowsEventDecodeFunc)
	b.running = false
	b.ctx, b.cancel = context.WithCancel(context.Background())
	return b
}

// BinlogSyncer syncs binlog event from server.
type BinlogSyncer struct {
	parser *BinlogParser
	c      *Conn
	cfg    BinlogSyncerConfig

	nextPos            Position
	prevGset, currGset GTIDSet
	// instead of GTIDSet.Clone, use this to speed up calculate prevGset
	prevMySQLGTIDEvent *GTIDEvent

	m      sync.RWMutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	running          bool
	retryCount       int
	lastConnectionID uint32
}

// StartBackup  Like mysqlbinlog remote raw backup
// Backup remote binlog from position (filename, offset) and write in backupDir
func (b *BinlogSyncer) StartBackup(backupDir string, p Position, timeout time.Duration) error {
	if timeout == 0 {
		// a very long timeout here
		timeout = 30 * 3600 * 24 * time.Second
	}

	// Force use raw mode
	b.parser.SetRawMode(true)

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return errors.Trace(err)
	}

	// [1] register slave
	// [2] enable semi-sync
	// [3] send dump cmd
	// [4] listen on events and send to streamer
	s, err := b.StartSync(p)
	if err != nil {
		return errors.Trace(err)
	}

	var f *os.File
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	// [5] get event from streamer
	// [6] handle event
	var filename string
	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		e, err := s.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			return nil
		}

		if err != nil {
			return errors.Trace(err)
		}

		offset := e.Header.LogPos
		switch e.Header.EventType {
		case ROTATE_EVENT:
			rotateEvent := e.Event.(*RotateEvent)
			filename = string(rotateEvent.NextLogName)
			if e.Header.Timestamp == 0 || offset == 0 {
				// fake rotate event
				continue
			}
		case FORMAT_DESCRIPTION_EVENT:
			// FormateDescriptionEvent is the first event in binlog, we will close old one and create a new
			if f != nil {
				f.Close()
			}

			if len(filename) == 0 {
				return errors.Errorf("empty binlog filename for FormateDescriptionEvent")
			}

			f, err = os.OpenFile(path.Join(backupDir, filename), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return errors.Trace(err)
			}

			// write binlog header fe'bin'
			if _, err = f.Write(BinLogFileHeader); err != nil {
				return errors.Trace(err)
			}
		}

		n, err := f.Write(e.RawData)
		if err != nil {
			return errors.Trace(err)
		}

		if n != len(e.RawData) {
			return errors.Trace(io.ErrShortWrite)
		}
	}
}

// StartSync starts syncing from the `pos` position.
func (b *BinlogSyncer) StartSync(pos Position) (*BinlogStreamer, error) {
	b.m.Lock()
	defer b.m.Unlock()

	if b.running {
		return nil, errors.Trace(errSyncRunning)
	}

	// [1] register slave
	// [2] semi-sync
	// [3] send dump cmd
	if err := b.prepareSyncPos(pos); err != nil {
		return nil, errors.Trace(err)
	}

	// [4] turn on running
	// [5] create new streamer
	// [6] listen on event
	return b.startDumpStream(), nil
}

func (b *BinlogSyncer) prepareSyncPos(pos Position) error {
	// always start from position 4
	if pos.Pos < 4 {
		pos.Pos = 4
	}

	// [1] register slave
	// [2] enable semi-sync
	if err := b.prepare(); err != nil {
		return errors.Trace(err)
	}

	// [3] send dump cmd
	if err := b.writeBinlogDumpCommand(pos); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncer) prepare() error {
	if b.isClosed() {
		return errors.Trace(ErrSyncClosed)
	}

	if err := b.registerSlave(); err != nil {
		return errors.Trace(err)
	}

	if err := b.enableSemiSync(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncer) isClosed() bool {
	select {
	case <-b.ctx.Done():
		return true
	default:
		return false
	}
}

// [1] create new connection and apply options
// [2] request server to set conn charset
// [3] set conn read timeout
// [4] set conn read buffer
// [5] request server to kill last conn
// [6] handle conn checksum
// [7] request server to set heartbeat
// [8] request server to register slave
func (b *BinlogSyncer) registerSlave() error {
	if b.c != nil {
		b.c.Close()
	}

	var err error
	b.c, err = b.newConnection(b.ctx)
	if err != nil {
		return errors.Trace(err)
	}

	if b.cfg.Option != nil {
		if err = b.cfg.Option(b.c); err != nil {
			return errors.Trace(err)
		}
	}

	if len(b.cfg.Charset) != 0 {
		if err = b.c.SetCharset(b.cfg.Charset); err != nil {
			return errors.Trace(err)
		}
	}

	//set read timeout
	if b.cfg.ReadTimeout > 0 {
		_ = b.c.SetReadDeadline(time.Now().Add(b.cfg.ReadTimeout))
	}

	if b.cfg.RecvBufferSize > 0 {
		if tcp, ok := b.c.Packet.Conn.(*net.TCPConn); ok {
			_ = tcp.SetReadBuffer(b.cfg.RecvBufferSize)
		}
	}

	// kill last connection id
	if b.lastConnectionID > 0 {
		b.killConnection(b.c, b.lastConnectionID)
	}

	// save last last connection id for kill
	b.lastConnectionID = b.c.GetConnectionID()

	//for mysql 5.6+, binlog has a crc32 checksum
	//before mysql 5.6, this will not work, don't matter.:-)
	r, err := b.c.Execute("SHOW GLOBAL VARIABLES LIKE 'BINLOG_CHECKSUM'")
	if err != nil {
		return errors.Trace(err)
	}

	s, _ := r.GetString(0, 1)
	if s != "" {
		// maybe CRC32 or NONE

		// mysqlbinlog.cc use NONE, see its below comments:
		// Make a notice to the server that this client
		// is checksum-aware. It does not need the first fake Rotate
		// necessary checksummed.
		// That preference is specified below.

		if _, err = b.c.Execute(`SET @master_binlog_checksum='NONE'`); err != nil {
			return errors.Trace(err)
		}

		// if _, err = b.c.Execute(`SET @master_binlog_checksum=@@global.binlog_checksum`); err != nil {
		// 	return errors.Trace(err)
		// }
	}

	if b.cfg.Flavor == MariaDBFlavor {
		// Refer https://github.com/alibaba/canal/wiki/BinlogChange(MariaDB5&10)
		// Tell the server that we understand GTIDs by setting our slave capability
		// to MARIA_SLAVE_CAPABILITY_GTID = 4 (MariaDB >= 10.0.1).
		if _, err := b.c.Execute("SET @mariadb_slave_capability=4"); err != nil {
			return errors.Errorf("failed to set @mariadb_slave_capability=4: %v", err)
		}
	}

	if b.cfg.HeartbeatPeriod > 0 {
		_, err = b.c.Execute(fmt.Sprintf("SET @master_heartbeat_period=%d;", b.cfg.HeartbeatPeriod))
		if err != nil {
			fmt.Printf("failed to set @master_heartbeat_period=%d, err: %v\n", b.cfg.HeartbeatPeriod, err)
			return errors.Trace(err)
		}
	}

	if err = b.writeRegisterSlaveCommand(); err != nil {
		return errors.Trace(err)
	}

	if _, err = b.c.ReadOKPacket(); err != nil {
		return errors.Trace(err)
	}

	serverUUID, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("failed to get new uud %v\n", err)
		return errors.Trace(err)
	}

	if _, err = b.c.Execute(fmt.Sprintf("SET @slave_uuid = '%s', @replica_uuid = '%s'", serverUUID, serverUUID)); err != nil {
		fmt.Printf("failed to set @slave_uuid = '%s', err: %v\n", serverUUID, err)
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncer) enableSemiSync() error {
	if !b.cfg.SemiSyncEnabled {
		return nil
	}

	r, err := b.c.Execute("SHOW VARIABLES LIKE 'rpl_semi_sync_master_enabled';")
	if err != nil {
		return errors.Trace(err)
	}

	s, _ := r.GetString(0, 1)
	if s != "ON" {
		fmt.Println("master does not support semi synchronous replication, use no semi-sync")
		b.cfg.SemiSyncEnabled = false
		return nil
	}

	_, err = b.c.Execute(`SET @rpl_semi_sync_slave = 1;`)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncer) writeBinlogDumpCommand(p Position) error {
	b.c.ResetSequence()

	data := make([]byte, 4+1+4+2+4+len(p.Name))

	pos := 4
	data[pos] = COM_BINLOG_DUMP
	pos++

	binary.LittleEndian.PutUint32(data[pos:], p.Pos)
	pos += 4

	binary.LittleEndian.PutUint16(data[pos:], b.cfg.DumpCommandFlag)
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	copy(data[pos:], p.Name)

	return b.c.WritePacket(data)
}

func (b *BinlogSyncer) startDumpStream() *BinlogStreamer {
	b.running = true

	s := NewBinlogStreamer()

	b.wg.Add(1)
	go b.onStream(s)
	return s
}

func (b *BinlogSyncer) onStream(s *BinlogStreamer) {
	defer func() {
		if e := recover(); e != nil {
			s.closeWithError(fmt.Errorf("Err: %v\n Stack: %s", e, Pstack()))
		}
		b.wg.Done()
	}()

	// for each packet
	//   parse event
	//   send event to streamer channel
	// handle retry
	for {
		data, err := b.c.ReadPacket()
		select {
		case <-b.ctx.Done():
			s.close()
			return
		default:
		}

		if err != nil {
			fmt.Println(err)
			// we meet connection error, should re-connect again with
			// last nextPos or nextGTID we got.
			if len(b.nextPos.Name) == 0 && b.prevGset == nil {
				// we can't get the correct position, close.
				s.closeWithError(err)
				return
			}

			if b.cfg.DisableRetrySync {
				fmt.Println("retry sync is disabled")
				s.closeWithError(err)
				return
			}

			for {
				select {
				case <-b.ctx.Done():
					s.close()
					return
				case <-time.After(time.Second):
					b.retryCount++
					if err = b.retrySync(); err != nil {
						if b.cfg.MaxReconnectAttempts > 0 && b.retryCount >= b.cfg.MaxReconnectAttempts {
							fmt.Printf("retry sync err: %v, exceeded max retries (%d)\n", err, b.cfg.MaxReconnectAttempts)
							s.closeWithError(err)
							return
						}

						fmt.Printf("retry sync err: %v, wait 1s and retry again\n", err)
						continue
					}
				}

				break
			}

			// we connect the server and begin to re-sync again.
			continue
		}

		//set read timeout
		if b.cfg.ReadTimeout > 0 {
			_ = b.c.SetReadDeadline(time.Now().Add(b.cfg.ReadTimeout))
		}

		// Reset retry count on successful packet receieve
		b.retryCount = 0

		switch data[0] {
		case OK_HEADER:
			if err = b.parseEvent(s, data); err != nil {
				s.closeWithError(err)
				return
			}
		case ERR_HEADER:
			err = b.c.HandleErrorPacket(data)
			s.closeWithError(err)
			return
		case EOF_HEADER:
			// refer to https://dev.mysql.com/doc/internals/en/com-binlog-dump.html#binlog-dump-non-block
			// when COM_BINLOG_DUMP command use BINLOG_DUMP_NON_BLOCK flag,
			// if there is no more event to send an EOF_Packet instead of blocking the connection
			fmt.Println("receive EOF packet, no more binlog event now.")
			continue
		default:
			fmt.Printf("invalid stream header %c\n", data[0])
			continue
		}
	}
}

// [1] data 是从 conn 读到的数据
// [2] parse 出的 event 发送到 streamer 的 channel
//
// [0] 1 byte ok header
// [1] 1 byte semi sync flag
// [2:] event
func (b *BinlogSyncer) parseEvent(s *BinlogStreamer, data []byte) error {
	//skip OK byte, 0x00
	data = data[1:]

	needACK := false
	if b.cfg.SemiSyncEnabled && (data[0] == SemiSyncIndicator) {
		needACK = data[1] == 0x01
		//skip semi sync header
		data = data[2:]
	}

	e, err := b.parser.Parse(data)
	if err != nil {
		return errors.Trace(err)
	}

	if e.Header.LogPos > 0 {
		// Some events like FormatDescriptionEvent return 0, ignore.
		b.nextPos.Pos = e.Header.LogPos
	}

	getCurrentGtidSet := func() GTIDSet {
		if b.currGset == nil {
			return nil
		}
		return b.currGset.Clone()
	}

	switch event := e.Event.(type) {
	case *RotateEvent:
		b.nextPos.Name = string(event.NextLogName)
		b.nextPos.Pos = uint32(event.Position)
		fmt.Printf("rotate to %s\n", b.nextPos)
	case *GTIDEvent:
		if b.prevGset == nil {
			break
		}
		if b.currGset == nil {
			b.currGset = b.prevGset.Clone()
		}
		u, _ := uuid.FromBytes(event.SID)
		b.currGset.(*MySQLGTIDSet).AddGTID(u, event.GNO)
		if b.prevMySQLGTIDEvent != nil {
			u, _ = uuid.FromBytes(b.prevMySQLGTIDEvent.SID)
			b.prevGset.(*MySQLGTIDSet).AddGTID(u, b.prevMySQLGTIDEvent.GNO)
		}
		b.prevMySQLGTIDEvent = event
	// case *MariadbGTIDEvent:
	// 	if b.prevGset == nil {
	// 		break
	// 	}
	// 	if b.currGset == nil {
	// 		b.currGset = b.prevGset.Clone()
	// 	}
	// 	prev := b.currGset.Clone()
	// 	err = b.currGset.(*MariadbGTIDSet).AddSet(&event.GTID)
	// 	if err != nil {
	// 		return errors.Trace(err)
	// 	}
	// 	// right after reconnect we will see same gtid as we saw before, thus currGset will not get changed
	// 	if !b.currGset.Equal(prev) {
	// 		b.prevGset = prev
	// 	}
	case *XIDEvent:
		if !b.cfg.DiscardGTIDSet {
			event.GSet = getCurrentGtidSet()
		}
	case *QueryEvent:
		if !b.cfg.DiscardGTIDSet {
			event.GSet = getCurrentGtidSet()
		}
	}

	needStop := false
	select {
	case s.ch <- e:
	case <-b.ctx.Done():
		needStop = true
	}

	if needACK {
		err := b.replySemiSyncACK(b.nextPos)
		if err != nil {
			return errors.Trace(err)
		}
	}

	if needStop {
		return errors.New("sync is been closing...")
	}

	return nil
}

func (b *BinlogSyncer) replySemiSyncACK(p Position) error {
	b.c.ResetSequence()

	data := make([]byte, 4+1+8+len(p.Name))
	pos := 4
	// semi sync indicator
	data[pos] = SemiSyncIndicator
	pos++

	binary.LittleEndian.PutUint64(data[pos:], uint64(p.Pos))
	pos += 8

	copy(data[pos:], p.Name)

	err := b.c.WritePacket(data)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncer) writeRegisterSlaveCommand() error {
	b.c.ResetSequence()

	hostname := b.localHostname()

	// This should be the name of slave host not the host we are connecting to.
	data := make([]byte, 4+1+4+1+len(hostname)+1+len(b.cfg.User)+1+len(b.cfg.Password)+2+4+4)
	pos := 4

	data[pos] = COM_REGISTER_SLAVE
	pos++

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	// This should be the name of slave hostname not the host we are connecting to.
	data[pos] = uint8(len(hostname))
	pos++
	n := copy(data[pos:], hostname)
	pos += n

	data[pos] = uint8(len(b.cfg.User))
	pos++
	n = copy(data[pos:], b.cfg.User)
	pos += n

	data[pos] = uint8(len(b.cfg.Password))
	pos++
	n = copy(data[pos:], b.cfg.Password)
	pos += n

	binary.LittleEndian.PutUint16(data[pos:], b.cfg.Port)
	pos += 2

	//replication rank, not used
	binary.LittleEndian.PutUint32(data[pos:], 0)
	pos += 4

	// master ID, 0 is OK
	binary.LittleEndian.PutUint32(data[pos:], 0)

	return b.c.WritePacket(data)
}

// localHostname returns the hostname that register slave would register as.
func (b *BinlogSyncer) localHostname() string {
	if len(b.cfg.Localhost) == 0 {
		h, _ := os.Hostname()
		return h
	}

	return b.cfg.Localhost
}

func (b *BinlogSyncer) newConnection(ctx context.Context) (*Conn, error) {
	var addr string
	if b.cfg.Port != 0 {
		addr = net.JoinHostPort(b.cfg.Host, strconv.Itoa(int(b.cfg.Port)))
	} else {
		addr = b.cfg.Host
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return ConnectWithDialer(
		timeoutCtx, "", addr, b.cfg.User, b.cfg.Password, "",
		b.cfg.Dialer,
		func(c *Conn) { c.SetTLSConfig(b.cfg.TLSConfig) },
	)
}

func (b *BinlogSyncer) killConnection(conn *Conn, id uint32) {
	cmd := fmt.Sprintf("KILL %d", id)
	if _, err := conn.Execute(cmd); err != nil {
		fmt.Printf("kill connection %d error %v\n", id, err)
		// Unknown thread id
		if code := ErrorCode(err.Error()); code != ER_NO_SUCH_THREAD {
			fmt.Println(errors.Trace(err))
		}
	}
	fmt.Printf("kill last connection id %d\n", id)
}

func (b *BinlogSyncer) retrySync() error {
	b.m.Lock()
	defer b.m.Unlock()

	b.parser.Reset()
	b.prevMySQLGTIDEvent = nil

	if b.prevGset != nil {
		msg := fmt.Sprintf("begin to re-sync from %s", b.prevGset.String())
		if b.currGset != nil {
			msg = fmt.Sprintf("%v (last read GTID=%v)", msg, b.currGset)
		}
		fmt.Println(msg)

		if err := b.prepareSyncGTID(b.prevGset); err != nil {
			return errors.Trace(err)
		}
	} else {
		fmt.Printf("begin to re-sync from %s\n", b.nextPos)
		if err := b.prepareSyncPos(b.nextPos); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (b *BinlogSyncer) prepareSyncGTID(gset GTIDSet) error {
	var err error

	// re establishing network connection here and will start getting binlog events from "gset + 1", thus until first
	// MariadbGTIDEvent/GTIDEvent event is received - we effectively do not have a "current GTID"
	b.currGset = nil

	// register slave
	// enable semi-sync
	if err = b.prepare(); err != nil {
		return errors.Trace(err)
	}

	switch b.cfg.Flavor {
	case MariaDBFlavor:
		err = b.writeBinlogDumpMariadbGTIDCommand(gset)
	default:
		// default use MySQL
		err = b.writeBinlogDumpMysqlGTIDCommand(gset)
	}

	if err != nil {
		return err
	}
	return nil
}

func (b *BinlogSyncer) writeBinlogDumpMariadbGTIDCommand(gset GTIDSet) error {
	// Copy from vitess

	startPos := gset.String()

	// Set the slave_connect_state variable before issuing COM_BINLOG_DUMP to
	// provide the start position in GTID form.
	query := fmt.Sprintf("SET @slave_connect_state='%s'", startPos)

	if _, err := b.c.Execute(query); err != nil {
		return errors.Errorf("failed to set @slave_connect_state='%s': %v", startPos, err)
	}

	// Real slaves set this upon connecting if their gtid_strict_mode option was
	// enabled. We always use gtid_strict_mode because we need it to make our
	// internal GTID comparisons safe.
	if _, err := b.c.Execute("SET @slave_gtid_strict_mode=1"); err != nil {
		return errors.Errorf("failed to set @slave_gtid_strict_mode=1: %v", err)
	}

	// Since we use @slave_connect_state, the file and position here are ignored.
	return b.writeBinlogDumpCommand(Position{Name: "", Pos: 0})
}

func (b *BinlogSyncer) writeBinlogDumpMysqlGTIDCommand(gset GTIDSet) error {
	p := Position{Name: "", Pos: 4}
	gtidData := gset.Encode()

	b.c.ResetSequence()

	data := make([]byte, 4+1+2+4+4+len(p.Name)+8+4+len(gtidData))
	pos := 4
	data[pos] = COM_BINLOG_DUMP_GTID
	pos++

	binary.LittleEndian.PutUint16(data[pos:], 0)
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	binary.LittleEndian.PutUint32(data[pos:], uint32(len(p.Name)))
	pos += 4

	n := copy(data[pos:], p.Name)
	pos += n

	binary.LittleEndian.PutUint64(data[pos:], uint64(p.Pos))
	pos += 8

	binary.LittleEndian.PutUint32(data[pos:], uint32(len(gtidData)))
	pos += 4
	n = copy(data[pos:], gtidData)
	pos += n

	data = data[0:pos]

	return b.c.WritePacket(data)
}

// StartSyncGTID starts syncing from the `gset` GTIDSet.
func (b *BinlogSyncer) StartSyncGTID(gset GTIDSet) (*BinlogStreamer, error) {
	fmt.Printf("begin to sync binlog from GTID set %s\n", gset)

	b.prevMySQLGTIDEvent = nil
	b.prevGset = gset

	b.m.Lock()
	defer b.m.Unlock()

	if b.running {
		return nil, errors.Trace(errSyncRunning)
	}

	// establishing network connection here and will start getting binlog events from "gset + 1", thus until first
	// MariadbGTIDEvent/GTIDEvent event is received - we effectively do not have a "current GTID"
	b.currGset = nil

	if err := b.prepare(); err != nil {
		return nil, errors.Trace(err)
	}

	var err error
	switch b.cfg.Flavor {
	case MariaDBFlavor:
		err = b.writeBinlogDumpMariadbGTIDCommand(gset)
	default:
		// default use MySQL
		err = b.writeBinlogDumpMysqlGTIDCommand(gset)
	}

	if err != nil {
		return nil, err
	}

	return b.startDumpStream(), nil
}
